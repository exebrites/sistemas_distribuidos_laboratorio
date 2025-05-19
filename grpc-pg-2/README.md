# grpc-pg-2-grupo01

## Integrantes
- Exequiel Brites
- Marcos Natanael Da Silva
- Federico Javier Ramos
- Thiago Schiaffino Alejandro

## Instrucciones

#### 1. Compilar el archivo .proto

```bash
protoc --go_out=. --go-grpc_out=. proto/monitor.proto
```
#### 2. Ejecutamos el Servidor

```bash
go run servidor/main.go
```

#### 3. Ejecutamos 5 Clientes (nodos) en terminales diferentes

```bash
go run cliente/main.go nodo1
go run cliente/main.go nodo2
go run cliente/main.go nodo3
go run cliente/main.go nodo4
go run cliente/main.go nodo5
```
### Pregunta

❓¿Qué implica aumentar o disminuir la frecuencia de envío de heartbeats?

La frecuencia de envío de heartbeats (es decir, cada cuánto tiempo los clientes notifican al servidor que siguen vivos) afecta directamente la rapidez con la que se detectan fallas y el consumo de red del sistema.

🔺 Si AUMENTAMOS la frecuencia (por ejemplo, cada 1 segundo):

✅ El servidor detecta fallos más rápido.

❌ Se generan más mensajes, lo que puede saturar la red si hay muchos nodos.

🔻 Si DISMINUIMOS la frecuencia (por ejemplo, cada 30 segundos):

✅ Se reduce el tráfico de red y la carga del sistema.

❌ El servidor puede tardar más en detectar que un nodo ha fallado.

❌ Hay mayor riesgo de falsos positivos si una red está lenta (parece que el nodo cayó pero solo hubo demora).