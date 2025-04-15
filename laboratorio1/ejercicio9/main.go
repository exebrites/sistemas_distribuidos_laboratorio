/*
Simular un middleware donde un único publicador envía cada 1 segundo un mensaje (por 
ejemplo: "evento-X") a 3 suscriptores. Cada suscriptor está representado por una goroutine 
que escucha su propio canal y muestra los eventos recibidos. El sistema debe permitir que 
todos los suscriptores reciban el mismo mensaje simultáneamente y en igual orden.

// Aclaración: Agregue la función time.AfterFunc para que el programa se ejecute por 1 minuto y
// luego se cierre automáticamente.
*/

package main

import (
	"fmt"
	"time"
)

func main() {
	// Crear un canal para el publicador
	publisher := make(chan string)

	// Crear canales para los suscriptores
	subscriber1 := make(chan string)
	subscriber2 := make(chan string)
	subscriber3 := make(chan string)

	// Función para distribuir mensajes a los suscriptores
	go func() {
		for msg := range publisher {
			// Enviar el mensaje a todos los suscriptores
			subscriber1 <- msg
			subscriber2 <- msg
			subscriber3 <- msg
		}
	}()

	// Función para manejar cada suscriptor
	subscriberHandler := func(id int, ch <-chan string) {
		for msg := range ch {
			fmt.Printf("Suscriptor %d recibió: %s\n", id, msg)
		}
	}

	// Iniciar goroutines para los suscriptores
	go subscriberHandler(1, subscriber1)
	go subscriberHandler(2, subscriber2)
	go subscriberHandler(3, subscriber3)

	// Publicador envía mensajes cada 1 segundo
	go func() {
		counter := 1
		for {
			msg := fmt.Sprintf("evento-%d", counter)
			publisher <- msg
			counter++
			time.Sleep(1 * time.Second)
		}
	}()

	// Ejecutar el programa por un minuto
	time.AfterFunc(1*time.Minute, func() {
		close(publisher)
		close(subscriber1)
		close(subscriber2)
		close(subscriber3)
		fmt.Println("Finalizando el programa después de 1 minuto.")
	})

	// Evitar que el programa termine inmediatamente
	select {}
}
