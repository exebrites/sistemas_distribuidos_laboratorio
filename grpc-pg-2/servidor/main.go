package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"grpc-pg-2/proto" // Usamos el código generado por el .proto

	"google.golang.org/grpc"
)

// Definimos el servidor, que implementa el servicio Monitor
type servidor struct {
	proto.UnimplementedMonitorServer
	mu          sync.Mutex
	ultimaVista map[string]time.Time
}

// Implementa el método EnviarHeartbeat del servicio Monitor
func (s *servidor) EnviarHeartbeat(stream proto.Monitor_EnviarHeartbeatServer) error {
	var nodoId string

	for {
		hb, err := stream.Recv() // Espera un nuevo heartbeat del cliente
		if err == io.EOF {
			// Si el stream termina normalmente, responde con un Ack
			return stream.SendAndClose(&proto.Ack{Mensaje: "Stream cerrado"})
		}
		if err != nil {
			//log.Printf("Error en stream: %v", err)
			return err
		}

		nodoId = hb.NodoId
		s.mu.Lock()
		s.ultimaVista[nodoId] = time.Unix(hb.MarcaTiempo, 0) // Guarda la marca de tiempo
		s.mu.Unlock()
		//log.Printf("[HEARTBEAT] %v %v", nodoId, hb.MarcaTiempo)
	}
}

// Proceso que revisa periódicamente si algún nodo ha fallado
func (s *servidor) detectorFallas(intervalo time.Duration) {
	for {
		time.Sleep(intervalo)
		s.mu.Lock()
		ahora := time.Now()
		fmt.Println("---- Estado de nodos ----")
		for nodo, ultimo := range s.ultimaVista {
			// Si hace más de 3 veces el intervalo sin ver al nodo → falla
			if ahora.Sub(ultimo) > 3*intervalo {
				//log.Printf("Fallo en Nodo %v inactivo desde hace %.0fs", nodo, ahora.Sub(ultimo).Seconds())
				fmt.Printf("Nodo %v: Inactivo (último hace %.0fs)\n", nodo, ahora.Sub(ultimo).Seconds())
			} else {
				fmt.Printf("Nodo %v: Activo (último hace %.0fs)\n", nodo, ahora.Sub(ultimo).Seconds())
			}
		}
		fmt.Println("--------------------------")

		s.mu.Unlock()
	}
}

func main() {

	// Escucha conexiones en el puerto 8000
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	s := grpc.NewServer()

	// Crea el servidor y la estructura para registrar heartbeats
	servidor := &servidor{ultimaVista: make(map[string]time.Time)}

	// Registra el servicio Monitor para que sea accesible desde gRPC
	proto.RegisterMonitorServer(s, servidor)

	// Inicia la detección de fallas en segundo plano
	go servidor.detectorFallas(5 * time.Second)

	fmt.Println("Servidor escuchando en :8000")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
