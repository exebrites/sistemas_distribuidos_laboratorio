/*
Escriba una función SumarPares que reciba un slice de enteros y devuelva la suma de los
números pares. Implemente un programa que demuestre su funcionamiento
*/
// opcional cargar por consola los valores. 
// 1. escribir funcion SumarPares
// 2. recibir un slice
// 3. recorrer el slice
// 4. si el numero es par, sumarlo
// 5. devolver la suma
// 6. mostrar por pantalla la suma


package main

import (
	"fmt"
)

func main() {

	var slice = []int {1,2,3,4};
	sumar:= SumarPares(slice);
	fmt.Println("El valor de la suma de pares es:",sumar);
}

func SumarPares	(a []int )int{
	suma :=0;
 	for _, valor := range a{
		// fmt.Println(valor);
		//determinar si es par
		if valor%2==0 {
			// fmt.Println("es par");
			suma+=valor;
		}
	}
	return suma
 }