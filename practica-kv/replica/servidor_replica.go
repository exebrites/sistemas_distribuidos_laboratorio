package main 

import ( 
	"encoding/binary" 
    "context"
    "log"
    "sync"
    "time"
    "google.golang.org/grpc"
	
    "net"
    "os"
    "os/signal"
    "strconv"
    "syscall"

    "google.golang.org/grpc/credentials/insecure"

    pb "practica-kv/proto" // Ajusta esta ruta
)
 
//Parte del paso8
const (
    colorReset  = "\033[0m"
    colorRed    = "\033[31m"
    colorGreen  = "\033[32m"
    colorYellow = "\033[33m"
    colorBlue   = "\033[34m"
    colorPurple = "\033[35m"
    colorCyan   = "\033[36m"
)


// VectorReloj representa un reloj vectorial de longitud 3 (tres réplicas).
 type VectorReloj [3]uint64 


// ValorConVersion guarda el valor y su reloj vectorial asociado
type ValorConVersion struct {
	Valor []byte
	RelojVector VectorReloj
} 


// ServidorReplica implementa pb.ReplicaServer
type ServidorReplica struct {
	pb.UnimplementedReplicaServer

	mu sync.Mutex
	almacen map[string]ValorConVersion // map[clave]ValorConVersión
	relojVector VectorReloj
	idReplica int // 0, 1 o 2
	clientesPeer []pb.ReplicaClient // stubs gRPC a las otras réplicas
}

// NewServidorReplica crea una instancia de ServidorReplica
// idReplica: 0, 1 o 2
// peerAddrs: direcciones gRPC de los otros dos peers (ej.: []string{":50052", ":50053"})
// NewServidorReplica crea una instancia de ServidorReplica
func NewServidorReplica(idReplica int, peerAddrs []string) *ServidorReplica {
    sr := &ServidorReplica{
        almacen:      make(map[string]ValorConVersion),
        idReplica:    idReplica,
        clientesPeer: make([]pb.ReplicaClient, len(peerAddrs)),
    }

    // Crear conexiones gRPC a los peers con el nuevo método
    for i, addr := range peerAddrs {
        conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err != nil {
            log.Fatalf("no se pudo conectar a peer %s: %v", addr, err)
        }
        sr.clientesPeer[i] = pb.NewReplicaClient(conn)
    }

    return sr
}



// GuardarLocal recibe la petición del Coordinador para almacenar clave/valor
func (r *ServidorReplica) GuardarLocal(ctx context.Context, req *pb.SolicitudGuardar) (*pb.RespuestaGuardar, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    // 1. Fusionar primero con reloj del cliente si existe
    if len(req.RelojVector) > 0 {
        relojCliente := decodeVector(req.RelojVector)
        r.relojVector.Fusionar(relojCliente)
    }

    // 2. Incrementar nuestro componente
    r.relojVector.Incrementar(r.idReplica)

    // 3. Guardar en almacén local
    r.almacen[req.Clave] = ValorConVersion{
        Valor:       req.Valor,
        RelojVector: r.relojVector,
    }


    // Parte del paso8
    log.Printf("%sRéplica %d - GUARDAR clave: %s, valor: %s, reloj: %v%s", 
    colorGreen, r.idReplica, req.Clave, req.Valor, r.relojVector, colorReset)


    // 4. Preparar mutación para replicación
    mutacion := &pb.Mutacion{
        Tipo:        pb.Mutacion_GUARDAR,
        Clave:       req.Clave,
        Valor:       req.Valor,
        RelojVector: encodeVector(r.relojVector),
    }

    // 5. Replicación síncrona con timeout
    errCh := make(chan error, len(r.clientesPeer))
    var wg sync.WaitGroup

    for _, peer := range r.clientesPeer {
        if peer == nil {
            continue
        }
        wg.Add(1)
        go func(p pb.ReplicaClient) {
            defer wg.Done()
            ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
            defer cancel()
            _, err := p.ReplicarMutacion(ctx, mutacion)
            if err != nil {
                errCh <- err
            }
        }(peer)
    }
    wg.Wait()
    close(errCh)

    // Verificar errores de replicación
    var errores []error
    for err := range errCh {
        errores = append(errores, err)
    }

    if len(errores) > 0 {
        log.Printf("Errores en replicación: %v", errores)
        // Considerar revertir la operación si es crítica
    }

    return &pb.RespuestaGuardar{
        Exito:            true,
        NuevoRelojVector: encodeVector(r.relojVector),
    }, nil
}



// EliminarLocal recibe la petición del Coordinador para borrar una clave
func (r *ServidorReplica) EliminarLocal(ctx context.Context, req *pb.SolicitudEliminar) (*pb.RespuestaEliminar, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    // 1. Incrementar nuestro componente del reloj vectorial
    r.relojVector.Incrementar(r.idReplica)

    // 2. Fusionar con reloj del cliente si existe
    if len(req.RelojVector) > 0 {
        relojCliente := decodeVector(req.RelojVector)
        r.relojVector.Fusionar(relojCliente)
    }

    // 3. Borrar del mapa local (si existe)
    _, existe := r.almacen[req.Clave]
    delete(r.almacen, req.Clave)

    // 4. Construir mutación de eliminación
    mutacion := &pb.Mutacion{
        Tipo:        pb.Mutacion_ELIMINAR,
        Clave:       req.Clave,
        RelojVector: encodeVector(r.relojVector),
    }

    // 5. Replicar a peers asíncronamente
    go func() {
        for i, peer := range r.clientesPeer {
            if peer == nil {
                continue // Saltar peers no conectados
            }

            ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
            defer cancel()

            _, err := peer.ReplicarMutacion(ctx, mutacion)
            if err != nil {
                log.Printf("Error replicando eliminación a peer %d: %v", i, err)
            }
        }
    }()

    // 6. Responder al Coordinador
    return &pb.RespuestaEliminar{
        Exito:            existe, // true si la clave existía
        NuevoRelojVector: encodeVector(r.relojVector),
    }, nil
}



// ObtenerLocal retorna el valor y reloj vectorial de una clave en esta réplica
func (r *ServidorReplica) ObtenerLocal(ctx context.Context, req *pb.SolicitudObtener) (*pb.RespuestaObtener, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    // Buscar la clave en el almacén
    valor, existe := r.almacen[req.Clave]
    if !existe {
        return &pb.RespuestaObtener{
            Existe: false,
        }, nil
    }

    // Retornar valor y metadatos
    return &pb.RespuestaObtener{
        Valor:        valor.Valor,
        RelojVector: encodeVector(valor.RelojVector),
        Existe:      true,
    }, nil
}




func (r *ServidorReplica) ReplicarMutacion(ctx context.Context, m *pb.Mutacion) (*pb.Reconocimiento, error) {
    // Log de recepción inicial
    log.Printf("%sRéplica %d - Recibida mutación: Tipo: %v, Clave: %s, Reloj remoto: %v%s", 
        colorYellow, r.idReplica, m.Tipo, m.Clave, decodeVector(m.RelojVector), colorReset)

    r.mu.Lock()
    defer r.mu.Unlock()

    relojRemoto := decodeVector(m.RelojVector)
    valorLocal, existe := r.almacen[m.Clave]

    // Log del estado actual antes de procesar
    if existe {
        log.Printf("%sRéplica %d - Estado actual - Clave: %s, Valor: %s, Reloj local: %v%s",
            colorCyan, r.idReplica, m.Clave, string(valorLocal.Valor), valorLocal.RelojVector, colorReset)
    } else {
        log.Printf("%sRéplica %d - Clave %s no existe localmente%s",
            colorCyan, r.idReplica, m.Clave, colorReset)
    }

    // Determinar si debemos aplicar la mutación
    aplicar := false

    if !existe {
        aplicar = true
        log.Printf("%sRéplica %d - Aplicando mutación (clave nueva)%s",
            colorGreen, r.idReplica, colorReset)
    } else {
        if valorLocal.RelojVector.AntesDe(relojRemoto) {
            aplicar = true
            log.Printf("%sRéplica %d - Aplicando mutación (reloj remoto más nuevo)%s",
                colorGreen, r.idReplica, colorReset)
        } else if relojRemoto.AntesDe(valorLocal.RelojVector) {
            log.Printf("%sRéplica %d - Ignorando mutación (ya tenemos versión más reciente)%s",
                colorPurple, r.idReplica, colorReset)
        } else {
            // Conflicto de versiones concurrentes
            log.Printf("%sRéplica %d - CONFLICTO DETECTADO! Reloj local: %v, Reloj remoto: %v%s",
                colorRed, r.idReplica, valorLocal.RelojVector, relojRemoto, colorReset)
            
            // Política de resolución alternativa: comparar valores de reloj
            // Gana la mutación con mayor componente en la posición de la réplica actual
            if relojRemoto[r.idReplica] > valorLocal.RelojVector[r.idReplica] {
                aplicar = true
                log.Printf("%sRéplica %d - Resolución: Aceptando versión remota (mayor componente local)%s",
                    colorRed, r.idReplica, colorReset)
            } else {
                log.Printf("%sRéplica %d - Resolución: Conservando versión local%s",
                    colorRed, r.idReplica, colorReset)
            }
        }
    }

    if aplicar {
        // Fusionar relojes primero
        r.relojVector.Fusionar(relojRemoto)
        
        // Aplicar mutación
        if m.Tipo == pb.Mutacion_GUARDAR {
            r.almacen[m.Clave] = ValorConVersion{
                Valor:       m.Valor,
                RelojVector: relojRemoto,
            }
            log.Printf("%sRéplica %d - GUARDADO - Clave: %s, Valor: %s, Reloj: %v%s",
                colorBlue, r.idReplica, m.Clave, m.Valor, relojRemoto, colorReset)
        } else {
            delete(r.almacen, m.Clave)
            log.Printf("%sRéplica %d - ELIMINADO - Clave: %s%s",
                colorBlue, r.idReplica, m.Clave, colorReset)
        }

        // Incrementar nuestro reloj después de aplicar cambios
        r.relojVector.Incrementar(r.idReplica)
        
        log.Printf("%sRéplica %d - Nuevo reloj vectorial: %v%s",
            colorBlue, r.idReplica, r.relojVector, colorReset)
    }

    return &pb.Reconocimiento{
        Ok:            true,
        RelojVectorAck: encodeVector(r.relojVector),
    }, nil
}


func (r *ServidorReplica) resolverConflicto(m *pb.Mutacion, local *ValorConVersion) bool {
    relojRemoto := decodeVector(m.RelojVector)
    
    // Política 1: Preferir réplica con ID mayor
    // Buscamos el componente más significativo diferente
    for i := 2; i >= 0; i-- {
        if relojRemoto[i] > local.RelojVector[i] {
            return true
        } else if relojRemoto[i] < local.RelojVector[i] {
            return false
        }
    }
    
    // Política 2: Si todos los componentes son iguales, preferir la mutación más reciente
    // (asumiendo que la mutación tiene un timestamp)
    // return m.Timestamp > local.Timestamp
    
    // Por defecto (si no hay timestamp): aceptar la mutación remota
    return true
}






// Incrementar aumenta en 1 el componente correspondiente a la réplica que llama.
 func (vr *VectorReloj) Incrementar(idReplica int) { 
    vr[idReplica]++
} 

// Fusionar toma el máximo elemento a elemento entre dos vectores.
 func (vr *VectorReloj) Fusionar(otro VectorReloj) { 
    for i := 0; i < 3; i++ {
        if otro[i] > vr[i] {
            vr[i] = otro[i]
        }
    }
} 

// AntesDe devuelve true si vr < otro en el sentido estricto (strictly less).
 func (vr VectorReloj) AntesDe(otro VectorReloj) bool { 
  menor := false
    for i := 0; i < 3; i++ {
        if vr[i] > otro[i] {
            return false
        }
        if vr[i] < otro[i] {
            menor = true
        }
    }
    return menor
}


// encodeVector serializa el VectorReloj a []byte para enviarlo por gRPC.
func encodeVector(vr VectorReloj) []byte { 
  	buf := make([]byte, 8*3) 
	for i := 0; i < 3; i++ { 
    	binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], vr[i]) 
  	} 
	return buf 
} 

// decodeVector convierte []byte a VectorReloj.
 func decodeVector(b []byte) VectorReloj { 
	var vr VectorReloj; 
	for i := 0; i < 3; i++ { 
    vr[i] = binary.BigEndian.Uint64(b[i*8 : (i+1)*8]) 
  } 
	return vr;
}



func main() {
    // Verificar argumentos de línea de comandos
    if len(os.Args) != 5 {
        log.Fatalf("Uso: %s <idReplica> <direccionEscucha> <peer1> <peer2>", os.Args[0])
    }

    // 1. Parsear argumentos
    idReplica, err := strconv.Atoi(os.Args[1])
    if err != nil || idReplica < 0 || idReplica > 2 {
        log.Fatal("idReplica debe ser 0, 1 o 2")
    }

    direccionEscucha := os.Args[2]
    peerAddrs := []string{os.Args[3], os.Args[4]}

    // 2. Configurar y registrar el servidor gRPC
    lis, err := net.Listen("tcp", direccionEscucha)
    if err != nil {
        log.Fatalf("Error al escuchar en %s: %v", direccionEscucha, err)
    }

    // Configurar opciones del servidor
    serverOpts := []grpc.ServerOption{
        grpc.MaxRecvMsgSize(10 * 1024 * 1024), // 10MB
        grpc.MaxSendMsgSize(10 * 1024 * 1024), // 10MB
    }

    s := grpc.NewServer(serverOpts...)

    // 3. Crear instancia del servidor réplica
    replica := NewServidorReplica(idReplica, peerAddrs)
    pb.RegisterReplicaServer(s, replica)

    // 4. Configurar manejo de señales para apagado elegante
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-stop
        log.Println("Recibida señal de apagado, deteniendo servidor...")
        s.GracefulStop()
    }()

    // 5. Iniciar el servidor
    log.Printf("Réplica %d iniciando en %s (peers: %v)", idReplica, direccionEscucha, peerAddrs)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Error al servir: %v", err)
    }
}