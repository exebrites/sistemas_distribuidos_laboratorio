/*
EJERCICIO 7
Simular el acceso concurrente a un log compartido donde 10 goroutines, cada una representando un nodo
,intentan registrar un evento crítico cada 0.5 segundos
(por ejemplo: "nodo-3: temperatura alta" o "nodo-7: pérdida de conexión");

cada evento debe escribirse en un slice de strings protegido por sync.Mutex para garantizar la integridad del log.
*/

// 1. crear gorutina
// 2. crear slice  de strings -> archivo log
// 3. crear un mutex
// 4. definir una estructura

package main

import (
	"fmt"
	"sync"
	"time"
)

type Log struct {
	eventos []string
	cerrojo sync.Mutex
}

func (e *Log) RegistrarEvento(evento string) {
	e.cerrojo.Lock()
	e.eventos = append(e.eventos, evento)
	e.cerrojo.Unlock()
}

func main() {
	// fmt.Print("ejercicio 7")
	// definir nodos
	// registrar eventso
	// imprimir log

	log := Log{}
	for i := 1; i <= 10; i++ {
		//  log.RegistrarEvento(fmt.Sprintf("nodo-%d: temperatura alta", i))
		time.Sleep(500 * time.Millisecond)
		go log.RegistrarEvento(fmt.Sprintf("nodo-%d: perdida de conexion", i))

	}
	fmt.Print("\n")
	fmt.Println(log.eventos)

	// fmt.Println("fin")

}
