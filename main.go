package main

import (
	"fmt"
)

func main() {
	running := true
	var opcion int

	mostrarMenu()

	for running {
		if ret, _ := fmt.Scanln(&opcion); ret == 1 { //Scanea y guarda 1 en ret si el dato que leyo es del tipo de opcion. Esto restringe el scan a ints
			running = manejarOpciones(opcion)
		}
	}
}

func mostrarMenu() {
	fmt.Println("Seleccione la opción deseada y presione enter")
	fmt.Println("1) Crear base de datos")
	fmt.Println("2) Crear tablas")
	fmt.Println("3) Crear PK's & FK's")
	fmt.Println("4) Borrar PK's & FK's")
	fmt.Println("5) Cargar todos los datos")
	fmt.Println("6) Salir")
}
func manejarOpciones(opcion int) bool {
	switch {
	case opcion == 1:
		fmt.Println("Base de datos creada")
	case opcion == 2:
		fmt.Println("Tablas creadas")
	case opcion == 3:
		fmt.Println("PK's & FK's creadas")
	case opcion == 4:
		fmt.Println("PK's & FK's borradas")
	case opcion == 5:
		fmt.Println("Todos los datos cargados")
	case opcion == 6:
		return false
	default:
		fmt.Println("Ingrese un numero valido")
	}
	return true
}
