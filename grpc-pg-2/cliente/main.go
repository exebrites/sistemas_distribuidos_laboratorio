package main

import (
	"context"
	"log"
	"os"
	"time"

	"grpc-pg-2/proto" // Importamos el código generado desde el .proto

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Validamos: se necesita un argumento que identifique al nodo
	if len(os.Args) != 2 {
		log.Fatal("Debe especificar un id de nodo como argumento")
	}
	nodo := os.Args[1] // Lee el ID del nodo (ej: "nodo1")

	// Conexión al servidor gRPC en localhost:8000
	conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	// Crea un cliente para el servicio Monitor (definido en proto)
	cliente := proto.NewMonitorClient(conn)

	// Abre el stream para enviar heartbeats
	stream, err := cliente.EnviarHeartbeat(context.Background())
	if err != nil {
		log.Fatalf("No se pudo abrir el stream: %v", err)
	}

	// Loop que envía heartbeats cada 5 segundos
	for {
		hb := &proto.Heartbeat{
			NodoId:      nodo,              // ID de este cliente
			MarcaTiempo: time.Now().Unix(), // Marca de tiempo actual
		}

		// Envía el heartbeat al servidor
		if err := stream.Send(hb); err != nil {
			log.Fatalf("Error enviando heartbeat: %v", err)
		}
		log.Printf("[%v] Enviado heartbeat", nodo) // Log del envío
		time.Sleep(5 * time.Second)                // Espera 5 segundos
	}
}
