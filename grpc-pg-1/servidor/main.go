package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"grpc-pg-1/proto"

	"net"

	"google.golang.org/grpc"
)

// servidor es la estructura que implementa el servicio
type servidor struct {
	proto.UnimplementedServicioServer
	proto.UnimplementedSaludoServiceServer
	listaSaludados []string
	mutex          sync.Mutex
}

// Hola es el metodo que se encarga de responder a la peticion de Hola
func (s *servidor) Hola(ctx context.Context, req *proto.Requerimiento) (*proto.Respuesta, error) {
	// loguea la peticion recibida
	log.Printf("Recibido: %s", req.Nombre)
	// devuelve una respuesta con un mensaje que saluda al usuario
	return &proto.Respuesta{Mensaje: "Hola " + req.Nombre}, nil
}

func (s *servidor) Saludar(ctx context.Context, req *proto.Saludo) (*proto.Respuesta, error) {
	s.mutex.Lock()
	s.listaSaludados = append(s.listaSaludados, req.Nombre)
	s.mutex.Unlock()

	log.Printf("Saludar llamado con: %s", req.Nombre)
	return &proto.Respuesta{Mensaje: "Saludos " + req.Nombre}, nil
}

// Método para SaludoService.ListadoPersonas
func (s *servidor) ListadoPersonas(ctx context.Context, _ *proto.Vacio) (*proto.Lista, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return &proto.Lista{Personas: s.listaSaludados}, nil
}
func main() {
	// abre el puerto 8000 para escuchar peticiones
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		// si hubo un error, loguea el error y termina el programa
		log.Fatalf("Error al escuchar: %v", err)
	}
	// crea un servidor grpc
	s := grpc.NewServer()
	// registra el servidor como implementacion del servicio
	proto.RegisterServicioServer(s, &servidor{})
	proto.RegisterSaludoServiceServer(s, &servidor{})
	// imprime un mensaje para indicar que el servidor esta listo
	fmt.Println("Servidor escuchando en :8000")
	// servidor grpc inicia a escuchar peticiones
	if err := s.Serve(lis); err != nil {
		// si hubo un error, loguea el error y termina el programa
		log.Fatalf("Error al servir: %v", err)
	}
}
