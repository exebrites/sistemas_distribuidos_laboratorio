package main

import (
	"fmt"
	"os"
)

func main() {
	// Verificar que se haya pasado un argumento (el nombre del archivo)

	/*os.Args: Es un slice que contiene los argumentos pasados al programa.
	os.Args[0] es el nombre del programa, y os.Args[1] es el primer argumento.*/
	if len(os.Args) < 2 { //Si no se proporciona un argumento (len(os.Args) < 2), se muestra un mensaje de error y el programa termina.
		fmt.Println("Error: Debes proporcionar un nombre de archivo como argumento.")
		return
	}

	// Obtener el nombre del archivo desde los argumentos
	nombreArchivo := os.Args[1]

	// Leer el contenido del archivo
	//.ReadFile: Lee el archivo completo y devuelve su contenido como un slice de bytes ([]byte) y un error (err).
	contenido, err := os.ReadFile(nombreArchivo)

	if err != nil { //Si el archivo no existe o no se puede leer, err no será nil
		fmt.Printf("Error: El archivo '%s' no existe o no se puede leer.\n", nombreArchivo)
		return
	}

	// Mostrar el contenido del archivo
	fmt.Println("Contenido del archivo:")
	fmt.Println(string(contenido)) //Convierte el slice de bytes ([]byte) a una cadena legible.
}
