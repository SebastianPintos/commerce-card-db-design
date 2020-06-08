package sql

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func DbConnection() {
	db, err = sql.Open("postgres", "user=postgres host=localhost dbname=tarjeta sslmode=disable")
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

	_, _err = _db.Exec(`CREATE DATABASE tarjeta`)
	logErr(_err)
}

func CargarDB() {
	cargarDatos()
	cargarCierres()
}

func CrearTablas(){
	crearTablas();
}

func CrearPKyFK() {
	crearPK()
	crearFK()
}

func EliminarPKyFK() {
	eliminarFK()
	eliminarPK()
}

func GenerarLogicaConsumo() {
	spObtenerDisponible()

	spChequearRechazoLimites()
	spAgregarRechazo()
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

func TestearConsumo() {
	spTestConsumoRechazo()
	spTestConsumoAlerta()

	CorrerTest()

	// consumoValidoTest()
	// consumoTarjetaInvalidaTest()
	// consumoCodSeguridadInvalidoTest()
	// consumoExcedeLimiteTest()
	// consumoTarjetaExpiradaTest()
	// consumoTarjetaSuspendidaTest()
	// consumoAlerta1Test()
	// consumoAlerta5Test()
	// consumoAlerta32Test()

	/*Para ejecutar todos los test
	_, err = db.Query(
	`   SELECT consumoValidoTest(),
		consumoTarjetaInvalidaTest(),
		consumoCodSeguridadInvalidoTest(),
		consumoExcedeLimiteTest(),
		consumoTarjetaExpiradaTest(),
		consumoTarjetaSuspendidaTest(),
		consumoAlerta1Test(),
		consumoAlerta5Test(),
		consumoAlerta32Test();
		`)
	logErr(err)*/
}
