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

func CargarDB() {
	CargarDatos()
}

func GenerarLogicaConsumo() {
	spChequearRechazoLimites()
	spAgregarRechazo()
	spAutorizarCompra()

	spAutorizarCompra()

	spAgregarAlertaRechazo()
	trAgregarAlerta()

	spTestearConsumo()
	trAgregarConsumo()
	
	spSeguridadCompras()
	trSeguridadCompras()
}

func GenerarResumen() {
	spGenerarResumen()
}
