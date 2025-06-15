package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	
	pb "practica-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type VectorReloj [3]uint64

func decodeVectorLocal(b []byte) VectorReloj {
	var vr VectorReloj
	for i := 0; i < 3; i++ {
		vr[i] = binary.BigEndian.Uint64(b[i*8 : (i+1)*8])
	}
	return vr
}

func (vr VectorReloj) String() string {
	return fmt.Sprintf("[%d %d %d]", vr[0], vr[1], vr[2])
}

func runConflictClient(clientID string) {
	conn, err := grpc.NewClient("localhost:6000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("no se pudo conectar: %v", err)
	}
	defer conn.Close()

	cliente := pb.NewCoordinadorClient(conn)
	ctx := context.Background()

	valorUnico := "valor_" + clientID
	clave := "conflictoX"

	resp, err := cliente.Guardar(ctx, &pb.SolicitudGuardar{
		Clave: clave,
		Valor: []byte(valorUnico),
	})
	if err != nil {
		log.Fatalf("Error al guardar: %v", err)
	}

	log.Printf("Cliente %s guardó: %s. Reloj resultante: %v", 
		clientID, valorUnico, decodeVectorLocal(resp.NuevoRelojVector))
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Uso: go run cliente/conflicto/conflicto.go <client_id>")
	}
	runConflictClient(os.Args[1])
}