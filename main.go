package main

import (
	"fmt"

	noSQL "./no-sql"
	"./sql"
)

func main() {

	sql.DbConnection()

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
	fmt.Println("6) Agregar lógica de consumo virtual y alertas")
	fmt.Println("7) Ejecutar test de consumo")
	fmt.Println("8) Cargar Datos noSql en BoltDB")
	fmt.Println("9) Generar Resumenes")
	fmt.Println("10) Salir ")
}
func manejarOpciones(opcion int) bool {
	switch {
	case opcion == 1:
		sql.CrearDB()
		fmt.Println("Base de datos creada")
	case opcion == 2:
		sql.CrearTablas()
		fmt.Println("Tablas creadas")
	case opcion == 3:
		sql.CrearPKyFK()
		fmt.Println("PK's & FK's creadas")
	case opcion == 4:
		sql.EliminarPKyFK()
		fmt.Println("PK's & FK's borradas")
	case opcion == 5:
		sql.CargarDB()
		fmt.Println("Todos los datos cargados")
	case opcion == 6:
		sql.GenerarLogicaConsumo()
		fmt.Println("Se agregó lógica de consumo y alertas")
	case opcion == 7:
		sql.TestearConsumo()
		fmt.Println("Test de consumo ejecutados")
	case opcion == 8:
		noSQL.CargaDatosNoDB()
		fmt.Println("Todos los datos noSQL DB cargados")
	case opcion == 9:
		sql.GenerarResumen()
		fmt.Println("Resumenes Generados")
	case opcion == 10:
		return false
	default:
		fmt.Println("Ingrese un numero valido")
	}
	return true
}
