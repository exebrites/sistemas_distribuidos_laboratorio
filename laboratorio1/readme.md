# Laboratorio 1 – Introducción y concurrencia en Go

## Objetivos

- Familiarizarse con el proceso de escritura de aplicaciones secuenciales y concurrentes
  usando el lenguaje de programación Go.
- Usar gorutinas, canales, cerrojos y grupos de espera.
- Fomentar el trabajo grupal / en equipo

## Tareas

1. Escriba una función SumarPares que reciba un slice de enteros y devuelva la suma de los
   números pares. Implemente un programa que demuestre su funcionamiento
2. Desarrolle un programa que lea una línea de texto desde la entrada estándar y cuente e
   imprima cuántas palabras tiene. Busque ayuda en los paquetes strings y fmt
3. Implemente una estructura Alumno con nombre, una lista de notas y un método Promedio()
   que devuelve el promedio de notas. Escriba un programa que permita obtener el promedio
   de varios alumnos. Sugerencia: no es necesario que cargue los datos de los alumnos, puede
   definirlos al crearlos
4. Escriba una función que convierta de °C a °F y otra de °F a °C. Luego realice un menú para
   elegir qué conversión hacer y pida los datos por teclado
5. Escriba un programa que reciba como argumento un nombre de archivo y muestre por
   consola su contenido. Si el archivo no existe se debe mostrar al usuario un mensaje y
   terminar
6. Diseñar un sistema en anillo donde cinco nodos, representados por goroutines, se envían
   mensajes de heartbeat cada 1 segundo entre sí de manera cíclica a través de canales. El
   sistema debe funcionar por 1 minuto y terminar
7. Simular el acceso concurrente a un log compartido donde 10 goroutines, cada una
   representando un nodo, intentan registrar un evento crítico cada 0.5 segundos (por ejemplo:
   "nodo-3: temperatura alta" o "nodo-7: pérdida de conexión"); cada evento debe escribirse en
   un slice de strings protegido por sync.Mutex para garantizar la integridad del log.
8. Simular un sistema de monitoreo donde 1 goroutina, envían un ping a una lista de 3 nodos
   (nodo-1, nodo-2, nodo-3) cada 2 segundos; cada ping simula una latencia aleatoria entre 100
   y 500 milisegundos, la gorutina debe guardar el nodo con menor latencia de respuesta de
   cada ronda en un slice. Luego de 10 rondas se debe imprimir los resultados y terminar.
   Consulta: se debe proteger el slice con mutex.
9. Simular un middleware donde un único publicador envía cada 1 segundo un mensaje (por
   ejemplo: "evento-X") a 3 suscriptores. Cada suscriptor está representado por una goroutine
   que escucha su propio canal y muestra los eventos recibidos. El sistema debe permitir que
   todos los suscriptores reciban el mismo mensaje simultáneamente y en igual orden.
10. Desarrolle un programa que tenga una variable global x iniciada en 0 (cero) y una función
    incrementar() que incremente en 5 a la variable x. La gorutina principal debe lanzar 100
    gorutinas que invoquen a la función incrementar() y luego imprimir el valor de x. Ejecute su
    programa usando la bandera -race para detectar si hay una carrera de datos. Además, el valor
    final de x debe ser 500, pero es posible que observe que a veces es 490 o 495 u otros
    valores. Usando WaitGroup y Mutexes, corrija su programa para que imprima el valor
    correcto y no tenga una carrera de datos.
11. Escriba un programa que mediante el uso de un mutex global, escriba dos funciones donde:
    la función a() debe bloquear el mutex, invocar a la función b() y desbloquear el mutex; la
    función b() debe bloquear el mutex, imprimir “Hola mundo” y desbloquear el mutex. La
    gorutina principal debe invocar a la función a(). Explica que sucede al ejecutar el programa
