package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	pb "practica-kv/proto" // Asegúrate que esta ruta coincida con tu módulo

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServidorCoordinador implementa pb.CoordinadorServer.
type ServidorCoordinador struct {
	pb.UnimplementedCoordinadorServer
	listaReplicas []string // ej: []string{":50051", ":50052", ":50053"}
	indiceRR      uint64   // contador atómico para round-robin
}

// NewServidorCoordinador crea un Coordinador con direcciones de réplica.
func NewServidorCoordinador(replicas []string) *ServidorCoordinador {
	return &ServidorCoordinador{
		listaReplicas: replicas,
		indiceRR:      0,
	}
}

// elegirReplicaParaEscritura: round-robin simple (ignora la clave).
func (c *ServidorCoordinador) elegirReplicaParaEscritura(clave string) string {
	idx := atomic.AddUint64(&c.indiceRR, 1)
	return c.listaReplicas[int(idx)%len(c.listaReplicas)]
}

// elegirReplicaParaLectura: también round-robin.
func (c *ServidorCoordinador) elegirReplicaParaLectura() string {
	idx := atomic.AddUint64(&c.indiceRR, 1)
	return c.listaReplicas[int(idx)%len(c.listaReplicas)]
}

// Obtener: redirige petición de lectura a una réplica.
func (c *ServidorCoordinador) Obtener(ctx context.Context, req *pb.SolicitudObtener) (*pb.RespuestaObtener, error) {
    replicaAddr := c.elegirReplicaParaLectura()
    conn, err := grpc.NewClient(replicaAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, fmt.Errorf("no se pudo conectar a réplica: %v", err)
    }
    defer conn.Close()
    
    cliente := pb.NewReplicaClient(conn)
    return cliente.ObtenerLocal(ctx, req)
}



// Guardar: redirige petición de escritura a una réplica elegida.
func (c *ServidorCoordinador) Guardar(ctx context.Context, req *pb.SolicitudGuardar) (*pb.RespuestaGuardar, error) {
    replicaAddr := c.elegirReplicaParaEscritura(req.Clave)
    
    // Conexión actualizada usando NewClient y WithTransportCredentials
    conn, err := grpc.NewClient(
        replicaAddr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        return nil, fmt.Errorf("error al conectar con la réplica %s: %v", replicaAddr, err)
    }
    defer conn.Close()

    cliente := pb.NewReplicaClient(conn)
    return cliente.GuardarLocal(ctx, req)
}



// Eliminar: redirige petición de eliminación a una réplica elegida.
func (c *ServidorCoordinador) Eliminar(ctx context.Context, req *pb.SolicitudEliminar) (*pb.RespuestaEliminar, error) {
    replicaAddr := c.elegirReplicaParaEscritura(req.Clave)
    
    // Conexión actualizada usando NewClient y WithTransportCredentials
    conn, err := grpc.NewClient(
        replicaAddr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        return nil, fmt.Errorf("error al conectar con la réplica %s: %v", replicaAddr, err)
    }
    defer conn.Close()

    cliente := pb.NewReplicaClient(conn)
    return cliente.EliminarLocal(ctx, req)
}



func main() {
	// Definir bandera para la dirección de escucha del Coordinador.
	listen := flag.String("listen", ":6000", "dirección para que escuche el Coordinador (p.ej., :6000)")
	flag.Parse()
	replicas := flag.Args()
	
	if len(replicas) < 3 {
		log.Fatalf("Debe proveer al menos 3 direcciones de réplicas, p.ej.: go run servidor_coordinador.go -listen :6000 :50051 :50052 :50053")
	}

	// Configurar servidor gRPC
	lis, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	s := grpc.NewServer()
	coordinador := NewServidorCoordinador(replicas)
	pb.RegisterCoordinadorServer(s, coordinador)

	log.Printf("Coordinador iniciando en %s (réplicas: %v)", *listen, replicas)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}