/*
Escriba un programa que mediante el uso de un mutex global, escriba dos funciones donde:
la función a() debe bloquear el mutex, invocar a la función b() y desbloquear el mutex; la
función b() debe bloquear el mutex, imprimir “Hola mundo” y desbloquear el mutex. La
gorutina principal debe invocar a la función a(). Explica que sucede al ejecutar el programa.

*/

package main

import (
	"fmt"
	"sync"
)

var mutex sync.Mutex

func a() {
	mutex.Lock()         // Bloquea el mutex.
	defer mutex.Unlock() // Asegura que se libere al final.
	b()                  // Llama a b() (sin necesidad de bloquear otra vez).
}

func b() {
	fmt.Println("Hola mundo") // Imprime sin bloquear el mutex.
}

func main() {
	a()
}