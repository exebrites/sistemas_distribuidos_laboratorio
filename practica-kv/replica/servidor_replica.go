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

    pb "practica-kv/proto" // Ajusta esta ruta
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
		relojVector:  VectorReloj{0, 0, 0},
		clientesPeer: make([]pb.ReplicaClient, len(peerAddrs)),
	}

	// Crear conexiones gRPC a los peers
	for i, addr := range peerAddrs {
		conn, err := grpc.Dial(addr, 
			grpc.WithInsecure(),
			grpc.WithTimeout(3*time.Second))
		if err != nil {
			log.Printf("No se pudo conectar a peer %s: %v", addr, err)
			continue
		}
		sr.clientesPeer[i] = pb.NewReplicaClient(conn)
	}

	return sr
}



// GuardarLocal recibe la petición del Coordinador para almacenar clave/valor
func (r *ServidorReplica) GuardarLocal(ctx context.Context, req *pb.SolicitudGuardar) (*pb.RespuestaGuardar, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 1. Incrementar nuestro componente del reloj vectorial
	r.relojVector.Incrementar(r.idReplica)

	// 2. Fusionar con reloj del cliente si existe
	if len(req.RelojVector) > 0 {
		relojCliente := decodeVector(req.RelojVector)
		r.relojVector.Fusionar(relojCliente)
	}

	// 3. Guardar en el mapa local
	r.almacen[req.Clave] = ValorConVersion{
		Valor:       req.Valor,
		RelojVector: r.relojVector,
	}

	// 4. Construir mutación para replicar a peers
	mutacion := &pb.Mutacion{
		Tipo:        pb.Mutacion_GUARDAR,
		Clave:       req.Clave,
		Valor:       req.Valor,
		RelojVector: encodeVector(r.relojVector),
	}

	// 5. Replicar asíncronamente a cada peer
	go func() {
		for i, peer := range r.clientesPeer {
			if peer == nil {
				continue // Saltar peers no conectados
			}
			
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			
			_, err := peer.ReplicarMutacion(ctx, mutacion)
			if err != nil {
				log.Printf("Error replicando a peer %d: %v", i, err)
			}
		}
	}()

	// 6. Responder al Coordinador con el nuevo reloj vectorial
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
    r.mu.Lock()
    defer r.mu.Unlock()

    // 1. Decodificar el reloj vectorial de la mutación
    relojRemoto := decodeVector(m.RelojVector)
    valorActual, existe := r.almacen[m.Clave]

    // 2. Determinar si aplicar la mutación
    aplicar := false
    if !existe {
        aplicar = true
    } else {
        relojLocal := valorActual.RelojVector
        
        if relojLocal.AntesDe(relojRemoto) {
            aplicar = true
        } else if !relojRemoto.AntesDe(relojLocal) {
            // Conflicto concurrente - resolver
            aplicar = r.resolverConflicto(m, &valorActual)
        }
    }

    // 3. Aplicar mutación si corresponde
    if aplicar {
        if m.Tipo == pb.Mutacion_GUARDAR {
            r.almacen[m.Clave] = ValorConVersion{
                Valor:       m.Valor,
                RelojVector: relojRemoto,
            }
        } else { // ELIMINAR
            delete(r.almacen, m.Clave)
        }
        r.relojVector.Fusionar(relojRemoto)
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
} 

// Fusionar toma el máximo elemento a elemento entre dos vectores.
 func (vr *VectorReloj) Fusionar(otro VectorReloj) { 
} 

// AntesDe devuelve true si vr < otro en el sentido estricto (strictly less).
 func (vr VectorReloj) AntesDe(otro VectorReloj) bool { 
  menor := false 

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