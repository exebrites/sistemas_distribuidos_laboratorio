# grpc-pg-1-grupo01

## Integrantes
- Exequiel Brites
- Marcos Natanael Da Silva
- Federico Javier Ramos
- Thiago Schiaffino Alejandro
## Instrucciones

### 1. Compilar el archivo .proto

```bash
protoc --go_out=. --go-grpc_out=. proto/servicio.proto
```
2. Ejecutar el servidor


```bash
go run servidor/main.go > listado_servidor.txt
```


3. Ejecutar el cliente


```bash
go run cliente/main.go > salida_cliente.txt

```