package sql

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func DbConnection() {
	db, err = sql.Open("postgres", "user=postgres host=localhost dbname=test2 sslmode=disable")
	logErr(err)
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CrearDB() {
	_db, _err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	logErr(_err)

	defer _db.Close()

	_, _err = _db.Exec(`CREATE DATABASE test2`)
	logErr(_err)
}

func cargarDB() {
	CrearTablas()
	crearSPs()
	crearTriggers()
}

func crearSPs() {
	spGenerarCierres()
	spGenerarResumen()
	generarLogicaConsumo()
	generarLogicaAlertas()
}

func generarLogicaConsumo() {
	//autorizarCompra()
	//generarConsumo()
	//crearTriggerConsumo()
	spChequearRechazoLimites()
	spAgregarRechazo()
	spAutorizarCompra()
	spTestearConsumo()
}

func generarLogicaAlertas() {
	//crearTriggerRechazo()
	//crearTriggersSeguridad()
	spAgregarAlerta()
	spSeguridadCompras()
}

func crearTriggers() {
	trAgregarConsumo()
	trAgregarRechazo()
	trSeguridadCompras()
}
