package main

import (
	"context"
	"log"
	"os"
	"time"

	"grpc-pg-2/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Debe especificar un id de nodo como argumento")
	}
	nodo := os.Args[1]

	conn, err := grpc.NewClient("localhost:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	cliente := proto.NewMonitorClient(conn)
	stream, err := cliente.EnviarHeartbeat(context.Background())
	if err != nil {
		log.Fatalf("No se pudo abrir el stream: %v", err)
	}

	for {
		hb := &proto.Heartbeat{
			NodoId:      nodo,
			MarcaTiempo: time.Now().Unix(),
		}
		if err := stream.Send(hb); err != nil {
			log.Fatalf("Error enviando heartbeat: %v", err)
		}
		log.Printf("[%v] Enviado heartbeat", nodo)
		time.Sleep(5 * time.Second)
	}
}
