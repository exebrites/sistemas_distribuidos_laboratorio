package main

import (
	"context"
	"fmt"
	"grpc-pg-1/proto"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// // Creamos un cliente de grpc con la direccion del servidor
	conn, err := grpc.Dial("localhost:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// Si hubo un error, lo logueamos y terminamos el programa
		log.Fatalf("No se pudo conectar: %v", err)
	}
	// Cerramos la conexion al finalizar el programa
	defer conn.Close()

	// Creamos un contexto con un tiempo de espera de 1 segundo
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second)
	// Cerramos el contexto al finalizar el programa
	defer cancel()

	fmt.Println("Menu:")
	fmt.Println("1. Saludar a Claudio")
	fmt.Println("2. Saludar a varias personas")
	var opcion int
	fmt.Print("Ingrese la opcion: ")
	fmt.Scanln(&opcion)
	switch opcion {
	case 1:
		saludarAClaudio(ctx, conn)
	case 2:
		llamarPersonas(ctx, conn)
	default:
		fmt.Println("Opcion no valida")
		return
	}
}

// saludarAClaudio llama al metodo Hola del servicio con el nombre "Claudio" y
// loguea la respuesta del servidor
func saludarAClaudio(ctx context.Context, conn *grpc.ClientConn) {
	cServicio := proto.NewServicioClient(conn)
	// // Llamamos al metodo Hola del servicio con el nombre "Claudio"
	r, err := cServicio.Hola(ctx, &proto.Requerimiento{Nombre: "Claudio"})
	// // Logueamos la respuesta del servidor
	log.Printf("Respuesta: %s", r.Mensaje)
	if err != nil {
		// Si hubo un error, lo logueamos y terminamos el programa
		log.Fatalf("Error al llamar al servidor: %v", err)
	}
}

// llamarPersonas llama al metodo Saludar del servicio con 10 personas distintas
// y loguea el listado de personas saludadas
func llamarPersonas(ctx context.Context, conn *grpc.ClientConn) {
	cSaludoServicio := proto.NewSaludoServiceClient(conn)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			nombre := fmt.Sprintf("Persona_%d", i)
			cSaludoServicio.Saludar(ctx, &proto.Saludo{Nombre: nombre})
		}(i)
	}

	wg.Wait()

	lista, _ := cSaludoServicio.ListadoPersonas(context.Background(), &proto.Vacio{})
	fmt.Println("Personas saludadas:")
	for _, nombre := range lista.Personas {
		fmt.Println(nombre)
	}

}
