/* 2. Desarrolle un programa que lea una línea de texto desde la entrada estándar y cuente e
imprima cuántas palabras tiene. Busque ayuda en los paquetes strings y fmt. */

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Crear un lector para la entrada estándar (teclado)
	lector := bufio.NewReader(os.Stdin)

	fmt.Println("Ingrese una línea de texto:")

	// Leer una línea de texto (hasta que el usuario presione Enter)
	texto, _ := lector.ReadString('\n')

	// Eliminar espacios en blanco al inicio/final (incluido el '\n')
	texto = strings.TrimSpace(texto)

	// Dividir el texto en palabras (los espacios son separadores)
	palabras := strings.Fields(texto)

	// Contar las palabras e imprimir el resultado
	fmt.Printf("Número de palabras: %d\n", len(palabras))
}
