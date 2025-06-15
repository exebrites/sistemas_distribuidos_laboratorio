package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	pb "practica-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Definición local del tipo VectorReloj (debe coincidir con el de las réplicas)
type VectorReloj [3]uint64

// Implementación local de decodeVector para el cliente
func decodeVectorLocal(b []byte) VectorReloj {
	var vr VectorReloj
	for i := 0; i < 3; i++ {
		vr[i] = binary.BigEndian.Uint64(b[i*8 : (i+1)*8])
	}
	return vr
}

// Método para imprimir el vector de forma legible
func (vr VectorReloj) String() string {
	return fmt.Sprintf("[%d %d %d]", vr[0], vr[1], vr[2])
}

func main() {
	// 1. Conectar al coordinador
	conn, err := grpc.NewClient("localhost:6000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("no se pudo conectar al coordinador: %v", err)
	}
	defer conn.Close()
	
	cliente := pb.NewCoordinadorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. Guardar clave-valor
	log.Println("\n=== Operación Guardar ===")
	respGuardar, err := cliente.Guardar(ctx, &pb.SolicitudGuardar{
		Clave: "usuario123",
		Valor: []byte("datosImportantes"),
	})
	if err != nil {
		log.Fatalf("Error al guardar: %v", err)
	}
	relojActual := decodeVectorLocal(respGuardar.NuevoRelojVector)
	log.Printf("Guardado exitoso. Reloj vectorial: %v", relojActual)

	// Pequeña pausa para replicación
	time.Sleep(300 * time.Millisecond)

	// 3. Obtener el valor
	log.Println("\n=== Operación Obtener (1ra vez) ===")
	respObtener1, err := cliente.Obtener(ctx, &pb.SolicitudObtener{
		Clave: "usuario123",
	})
	if err != nil {
		log.Fatalf("Error al obtener: %v", err)
	}
	
	if !respObtener1.Existe {
		log.Fatal("ERROR: La clave no existe después de guardar")
	}
	relojObtenido := decodeVectorLocal(respObtener1.RelojVector)
	log.Printf("Valor: %s", string(respObtener1.Valor))
	log.Printf("Reloj vectorial: %v", relojObtenido)

	// 4. Eliminar la clave
	log.Println("\n=== Operación Eliminar ===")
	_, err = cliente.Eliminar(ctx, &pb.SolicitudEliminar{
		Clave:       "usuario123",
		RelojVector: respObtener1.RelojVector,
	})
	if err != nil {
		log.Fatalf("Error al eliminar: %v", err)
	}
	log.Println("Eliminación exitosa")

	// Pequeña pausa para replicación
	time.Sleep(300 * time.Millisecond)

	// 5. Verificar eliminación
	log.Println("\n=== Operación Obtener (2da vez - verificación) ===")
	respObtener2, err := cliente.Obtener(ctx, &pb.SolicitudObtener{
		Clave: "usuario123",
	})
	if err != nil {
		log.Fatalf("Error al obtener: %v", err)
	}
	
	if respObtener2.Existe {
		log.Fatal("ERROR: La clave todavía existe después de eliminar")
	}
	log.Println("Eliminación verificada correctamente")
}