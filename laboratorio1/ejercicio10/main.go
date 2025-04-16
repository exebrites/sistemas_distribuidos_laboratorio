/*
Ejercicio10
Desarrolle un programa que tenga

una variable global x iniciada en 0 (cero) y una función incrementar() que incremente en 5 a la variable x. La gorutina principal debe lanzar 100
gorutinas que invoquen a la función incrementar() y luego imprimir el valor de x. Ejecute su
programa usando la bandera -race para detectar si hay una carrera de datos. Además, el valor
final de x debe ser 500, pero es posible que observe que a veces es 490 o 495 u otros
valores. Usando WaitGroup y Mutexes, corrija su programa para que imprima el valor
correcto y no tenga una carrera de datos.
*/
package main

import (
	"fmt"
	"sync"
	// "sync"
)

type Valor struct {
	x       int
	cerrojo sync.Mutex
}

func (a *Valor) Incrementar() {
	defer wg.Done()
	a.cerrojo.Lock()
	a.x += 5 // a.x = a.x + 5
	a.cerrojo.Unlock()
}

var wg sync.WaitGroup
var wg1 sync.WaitGroup

func main() {
	// fmt.Print("ejercicio 10")

	valor := &Valor{}
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go valor.Incrementar()

		}

		// wg.Wait()

	}()

	wg1.Wait() // espera a que la goroutine termine
	wg.Wait()  // Ahora sí esperamos a que terminen todas las goroutines
	fmt.Println("Ambas gorutinas han finalizado!")
	fmt.Print(valor.x)
}
