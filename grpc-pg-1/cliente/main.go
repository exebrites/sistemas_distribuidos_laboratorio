package main 
import ( 
"context" 
"log" 
"time" 
"grpc-pg-1/proto" 
"google.golang.org/grpc"
"google.golang.org/grpc/credentials/insecure" 
) 
func main() { 
conn, err := grpc.NewClient("localhost:8000", 
grpc.WithTransportCredentials(insecure.NewCredentials())) 
if err != nil { 
log.Fatalf("No se pudo conectar: %v", err) 
} 
defer conn.Close() 
c := proto.NewServicioClient(conn) 
ctx, cancel := context.WithTimeout(context.Background(), 
time.Second) 
defer cancel() 
r, err := c.Hola(ctx, &proto.Requerimiento{Nombre: 
"Claudio"}) 
if err != nil { 
log.Fatalf("Error al llamar al servidro: %v", err) 
} 
log.Printf("Respuesta: %s", r.Mensaje)
}
