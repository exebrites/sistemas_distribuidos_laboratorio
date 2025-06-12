package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "practica-kv/proto" // Asegúrate que coincida con tu módulo
)

func main() {
	// 1. Conectar al coordinador
	conn, err := grpc.Dial(":6000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No se pudo conectar al coordinador: %v", err)
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
	log.Printf("Guardado exitoso. Reloj vectorial: %v", respGuardar.NuevoRelojVector)

	// 3. Obtener el valor
	log.Println("\n=== Operación Obtener (1ra vez) ===")
	respObtener1, err := cliente.Obtener(ctx, &pb.SolicitudObtener{
		Clave: "usuario123",
	})
	if err != nil {
		log.Fatalf("Error al obtener: %v", err)
	}
	if !respObtener1.Existe {
		log.Println("Clave no existe (esto no debería ocurrir)")
	} else {
		log.Printf("Valor: %s", string(respObtener1.Valor))
		log.Printf("Reloj vectorial: %v", respObtener1.RelojVector)
	}

	// 4. Eliminar la clave
	log.Println("\n=== Operación Eliminar ===")
	_, err = cliente.Eliminar(ctx, &pb.SolicitudEliminar{
		Clave:       "usuario123",
		RelojVector: respObtener1.RelojVector, // Enviamos el reloj que obtuvimos
	})
	if err != nil {
		log.Fatalf("Error al eliminar: %v", err)
	}
	log.Println("Eliminación exitosa")

	// 5. Verificar eliminación
	log.Println("\n=== Operación Obtener (2da vez - verificación) ===")
	respObtener2, err := cliente.Obtener(ctx, &pb.SolicitudObtener{
		Clave: "usuario123",
	})
	if err != nil {
		log.Fatalf("Error al obtener: %v", err)
	}
	if !respObtener2.Existe {
		log.Println("Clave ya no existe (eliminación verificada)")
	} else {
		log.Println("¡Error! La clave todavía existe")
	}
}