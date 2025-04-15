/*
EJERCICIO 4

Escriba una función que convierta de °C a °F y otra de °F a °C. 
Luego realice un menú para elegir qué conversión hacer y pida los datos por teclado*/

// *SUPOSICION* los valores ingresados son enteros
// 1. FUNCION DE °C A °F
// 2. FUNCION DE °F A °C
// 3. MENU -> loop en infinito
// 4. ELEGIR LA CONVERSION
// 5. Tomar datos de consola
// 6. Mostrar por pantalla

package main

import (
	"bufio"//lectura del buffer
	"os"//entrada 
	"fmt"//imprimir
	// "strconv"// parsear
	"strings" //manejo de cadenas
)

func main(){

	//lectura de strings
	// lector := bufio.NewReader(os.Stdin);
	// fmt.Println("ingrese su nombre");
	// nombre, _ := lector.ReadString('\n');
	// fmt.Println(nombre);

	//lectura de enteros y conversion
	// lector := bufio.NewReader(os.Stdin);
	// intString, _ := lector.ReadString('\n');
	// entero,_:= strconv.Atoi(strings.TrimSpace(intString));
	// fmt.Println("Int:", entero);

	//    floatVal, _ := strconv.ParseFloat(strings.TrimSpace(floatStr), 64)

	// valor:=100.0;
	// conversion := conversionCelsiusFahrenheit(valor);

	// fmt.Println("conversion celsius a fahrenheit", conversion);

	// fmt.Println("conversion fahrenheit a celsius", conversionFahrenheitCelsius(valor));

	// MENU
	fmt.Println("BIENVENIDO AL MENU");
	fmt.Println("En caso de salir presionar 'c'")
	lector := bufio.NewReader(os.Stdin);

	a:="c";//parametro de control
	b:="a";//recibir de consola
	for {
		if strings.EqualFold(a, b) {
			break
		}
		// código aquí

		fmt.Println("1. Conversión de °C a °F");
		fmt.Println("2. Conversión de °F a °C");
		fmt.Println("Ingrese la opcion a realizar");
		opcion, _ := lector.ReadString('\n');
		opcion=strings.TrimSpace(opcion);
		switch opcion {
		case "1":
			// código para "a"
			fmt.Println("Ingrese la temperatura en °C");
		case "2":
			// código para "b"
			fmt.Println("Ingrese la temperatura en °F");
		default:
			fmt.Println("En caso de salir presionar 'c'")
		}

		b, _ = lector.ReadString('\n');
		b = strings.TrimSpace(b);
	}
}
// func conversionCelsiusFahrenheit(celsius float64) float64{
// 	return celsius * 9/5 + 32;
// }

// func conversionFahrenheitCelsius(fahrenheit float64) float64{
// 	return (fahrenheit - 32) * 5/9;
// }
