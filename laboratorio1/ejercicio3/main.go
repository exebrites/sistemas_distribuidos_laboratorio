/*

Ejercicio 3: Conversión de temperaturas

Implemente una estructura Alumno con nombre, una lista de notas y un método Promedio() 
que devuelve el promedio de notas. Escriba un programa que permita obtener el promedio 
de varios alumnos. Sugerencia: no es necesario que cargue los datos de los alumnos, puede 
definirlos al crearlos.
*/

package main

import (

	"fmt"//imprimir
)

type Alumno struct {
	Nombre string
	Notas  []float64
}

func (a Alumno) Promedio() float64 {
	var suma float64
	for _, nota := range a.Notas {
		suma += nota
	}
	return suma / float64(len(a.Notas))
}

func main() {
	alumnos := []Alumno{
		{"Juan", []float64{7.5, 8.0, 9.0}}, //8.17
		{"Maria", []float64{6.0, 7.0, 8.5}}, //7.17
		{"Pedro", []float64{5.0, 6.5, 7.0}}, //6.17
		{"Exequiel", []float64{9.0, 8.5, 9.5}},//9.0
		{"Federico", []float64{10.0, 9.5, 9.0}},//9.5
		{"Nata", []float64{8.0, 7.5, 8.0}}, //7.83
	}

	for _, alumno := range alumnos {
		fmt.Printf("El promedio de %s es %.2f\n", alumno.Nombre, alumno.Promedio())
	}
}
