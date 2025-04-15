/*
 	Diseñar un sistema en anillo donde cinco nodos, representados por goroutines, se envían
mensajes de heartbeat cada 1 segundo entre sí de manera cíclica a través de canales. El
sistema debe funcionar por 1 minuto y terminar.
*/

package main

import (
	"fmt"
	"time"
)

// node simula un nodo en el anillo. Recibe mensajes del canal 'in',
// los procesa y envía un mensaje de latido (heartbeat) al canal 'out'.
func node(id int, in <-chan string, out chan<- string) {
	for {
		select {
		case msg := <-in: // Recibe un mensaje del canal 'in'
			fmt.Printf("Nodo %d recibió: %s\n", id, msg)
			time.Sleep(1 * time.Second) // Simula un retraso en el procesamiento
			// Envía un mensaje de latido (heartbeat) al canal 'out'
			out <- fmt.Sprintf("Latido del Nodo %d", id)
		}
	}
}

func main() {
	const numNodes = 5           // Número de nodos en el anillo
	const duration = 1 * time.Minute // Duración durante la cual el sistema funcionará

	// Crea un slice de canales para conectar los nodos
	channels := make([]chan string, numNodes)
	for i := 0; i < numNodes; i++ {
		channels[i] = make(chan string) // Inicializa cada canal
	}

	// Crea y lanza los nodos como goroutines
	for i := 0; i < numNodes; i++ {
		// Cada nodo recibe mensajes de su canal 'in' y los envía al canal 'in' del siguiente nodo
		go node(i, channels[i], channels[(i+1)%numNodes])
	}

	// Inicia el flujo de mensajes enviando el primer mensaje al primer nodo
	go func() {
		channels[0] <- "Inicio"
	}()

	// Permite que el sistema funcione durante la duración especificada
	time.Sleep(duration)
	fmt.Println("Apagando el sistema...")
}