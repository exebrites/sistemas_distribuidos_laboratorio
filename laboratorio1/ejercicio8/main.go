/* 8. Simular un sistema de monitoreo donde 1 goroutina, envían un ping a una lista de 3 nodos
(nodo-1, nodo-2, nodo-3) cada 2 segundos; cada ping simula una latencia aleatoria entre 100
y 500 milisegundos, la gorutina debe guardar el nodo con menor latencia de respuesta de
cada ronda en un slice. Luego de 10 rondas se debe imprimir los resultados y terminar.
Consulta: se debe proteger el slice con mutex. */

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Almacena el número de ronda, el nodo con menor latencia y su valor.
type Resultado struct {
	Ronda    int
	Nodo     string
	Latencia int
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Inicializar semilla aleatoria

	var (
		wg         sync.WaitGroup
		mu         sync.Mutex  // Mutex para proteger el slice
		resultados []Resultado // Slice compartido
		nodos      = []string{"nodo-1", "nodo-2", "nodo-3"}
		rondas     = 10
		intervalo  = 2 * time.Second
	)

	for ronda := 1; ronda <= rondas; ronda++ {
		wg.Add(1)
		go func(ronda int) {
			defer wg.Done()

			// Canal para recibir latencias de los nodos
			ch := make(chan struct {
				Nodo     string
				Latencia int
			}, len(nodos))

			// Enviar ping a cada nodo (en paralelo)
			for _, nodo := range nodos {
				go func(n string) {
					latencia := rand.Intn(401) + 100 // 100-500ms
					time.Sleep(time.Duration(latencia) * time.Millisecond)
					ch <- struct {
						Nodo     string
						Latencia int
					}{n, latencia}
				}(nodo)
			}

			// Buscar el nodo con menor latencia en esta ronda
			minLatencia := 500
			var mejorNodo string
			for i := 0; i < len(nodos); i++ {
				respuesta := <-ch
				if respuesta.Latencia < minLatencia {
					minLatencia = respuesta.Latencia
					mejorNodo = respuesta.Nodo
				}
			}

			// Guardar resultado (protegido por Mutex)
			mu.Lock()
			resultados = append(resultados, Resultado{ronda, mejorNodo, minLatencia})
			mu.Unlock()

			fmt.Printf("Ronda %d: Mejor nodo = %s (latencia: %dms)\n", ronda, mejorNodo, minLatencia)
			time.Sleep(intervalo) // Esperar 2 segundos entre rondas
		}(ronda)
	}

	wg.Wait() // Esperar a que terminen todas las goroutines

	// Imprimir resultados finales
	fmt.Println("\nResultados finales:")
	for _, res := range resultados {
		fmt.Printf("Ronda %d: %s (%dms)\n", res.Ronda, res.Nodo, res.Latencia)
	}
}
