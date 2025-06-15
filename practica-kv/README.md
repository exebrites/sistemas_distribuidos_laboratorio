# Sistema Clave-Valor Multi-Master con Coordinador y Réplicas

Implementación de un sistema distribuido de almacenamiento clave-valor usando Go y gRPC, con soporte para replicación multi-master y manejo de conflictos mediante relojes vectoriales.

## Requisitos Previos

- Go 1.16 o superior
- Protocol Buffers (protoc)
- Plugins protoc-gen-go y protoc-gen-go-grpc

## Estructura del Proyecto
practica-kv/  
├── proto/  
│   └── kv.proto # Definición de servicios y mensajes gRPC  
├── coordinador/  
│   └── servidor_coordinador.go # Servidor coordinador  
├── replica/  
│   └── servidor_replica.go # Servidores réplica  
├── cliente/  
│   ├── cliente_ejemplo.go # Cliente básico  
│   └── conflicto/ # Cliente para pruebas de conflicto  
│   └── conflicto.go  
└── go.mod # Definición de módulo Go  


## Compilación

**Generar stubs gRPC:**

```bash
protoc --go_out=. --go-grpc_out=. proto/kv.proto
```


**Compilar componentes del proyecto(opcional):**  

 - **Réplicas**
    ```go
    go build -o bin/replica replica/servidor_replica.go
    ```

 - **Coordinador**
    ```go
    go build -o bin/coordinador coordinador/servidor_coordinador.go
    ```

 - **Clientes**
    ```go
    go build -o bin/cliente cliente/cliente_ejemplo.go
    go build -o bin/cliente-conflicto cliente/conflicto/conflicto.go
    ```

<br>  

## Ejecución

### Paso1 - Replicas
Abrir distintos terminales, iniciar y ejecutar en terminales separadas las siguientes replicas:

```go
// Réplica 0
go run replica/servidor_replica.go 0 :50051 :50052 :50053
```

```go
// Réplica 1
go run replica/servidor_replica.go 1 :50052 :50051 :50053
```

```go
// Réplica 2
go run replica/servidor_replica.go 2 :50053 :50051 :50052
```

### Paso2 - Iniciar el coordinador

```go
go run coordinador/servidor_coordinador.go -listen :6000 :50051 :50052 :50053
```


### Paso3 - Ejecutar cliente de ejemplo
```go
go run cliente/cliente_ejemplo.go
```

La salida esperada es la siguiente:

```bash
=== Operación Guardar ===
Guardado exitoso. Reloj vectorial: [0 1 0]

=== Operación Obtener (1ra vez) ===
Valor: datosImportantes
Reloj vectorial: [0 1 0]

=== Operación Eliminar ===
Eliminación exitosa

=== Operación Obtener (2da vez - verificación) ===
Eliminación verificada correctamente
```

### Paso4 - Prueba de conflicto
También se puede ejecutar el cliente de conflicto para simular un escenario de conflicto entre replicas, haciendo la ejecución de los clientes de conflicto en terminales separadas.  

**Comandos de ejecución**: 

```go
// Ejecutar en terminales separadas de manera simultanea

// Terminal 1
go run cliente/conflicto/conflicto.go cliente1

// Terminal 2
go run cliente/conflicto/conflicto.go cliente2
```

**Salida esperada**:

```bash
// Terminal 1
Cliente cliente1 guardó: valor_cliente1. Reloj resultante: [2 1 3]

// Terminal 2
Cliente cliente2 guardó: valor_cliente2. Reloj resultante: [4 1 3]
```


Logs de réplicas (mostrarán detección y resolución de conflictos):  

```bash
Réplica 0 - CONFLICTO DETECTADO! Reloj local: [2 1 3], Reloj remoto: [4 1 3]
Réplica 0 - Resolución: Aceptando versión remota
```

## Equipo del trabajo: Grupo 01  
- Brites Exequiel
- DaSilva Marcos N.
- Ramos Federico J.
- Schiaffino Alejandro Thiago
