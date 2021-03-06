package noSQL

import (
	"encoding/json"
	"log"
	"strconv"

	bolt "github.com/coreos/bbolt"
)

type Cliente struct {
	Nrocliente int
	Nombre     string
	Apellido   string
	Domicilio  string
	Telefono   string
}

type Tarjeta struct {
	Nrotarjeta   string
	Nrocliente   int
	Validadesde  string
	Validahasta  string
	Codseguridad string
	Limitecompra int
	Estado       string
}

type Comercio struct {
	Nrocomercio  int
	Nombre       string
	Domicilio    string
	Codigopostal string
	Telefono     string
}

type Compra struct {
	Nrooperacion int
	Nrotarjeta   string
	Nrocomercio  int
	Fecha        string
	Monto        int
	Pagado       bool
}

var db *bolt.DB
var err error

func dbConnection() {
	db, err = bolt.Open("./no-sql/boltDB/test.db", 0600, nil)
	logErr(err)
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CargaDatosNoDB() {

	dbConnection()

	cargarCliente(11348773, "Rocío", "Losada", "Av. Presidente Perón 1530", "1151102983")
	cargarCliente(12349972, "María Estela", "Martínez", "Belgrano 1830", "1150006655")
	cargarCliente(22648991, "Laura", "Santos", "Italia 220", "1153399452")

	cargarTarjeta("4000001234567899", 11348773, "201508", "202008", "733", 50000, "vigente")
	cargarTarjeta("4037001554363655", 12349972, "201507", "202007", "332", 55000, "vigente")
	cargarTarjeta("4000001355435322", 22648991, "201507", "202007", "201", 60000, "vigente")

	cargarComercio(501, "Kevingston", "Av. Tte. Gral. Ricchieri 965", "1661", "46666181")
	cargarComercio(523, "47 street", "Paunero 1575", "1663", "47597581")
	cargarComercio(513, "Garbarino", "Av. Bartolomé Mitre 1198", "1661", "08104440018")

	cargarCompra(1, "4000001234567899", 501, "2020-04-25 00:00:00", 1500.00, true)
	cargarCompra(2, "4000001234567899", 513, "2020-04-27 00:00:00", 4500.00, true)
	cargarCompra(3, "4000001234567899", 523, "2020-04-30 00:00:00", 850.00, true)
}

func cargarCliente(nrocliente int, nombre string, apellido string, domicilio string, telefono string) {
	cliente := Cliente{nrocliente, nombre, apellido, domicilio, telefono}

	data, err := json.Marshal(cliente)
	if err != nil {
		log.Fatal(err)
	}

	CreateUpdate(db, "Cliente", []byte(strconv.Itoa(cliente.Nrocliente)), data)
}

func cargarTarjeta(nrotarjeta string, nrocliente int, validadesde string, validahasta string, codseguridad string, limitecompra int, estado string) {
	tarjeta := Tarjeta{nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado}

	data, err := json.Marshal(tarjeta)
	if err != nil {
		log.Fatal(err)
	}

	CreateUpdate(db, "Tarjeta", []byte(strconv.Itoa(tarjeta.Nrocliente)), data)

	// consulta, err := ReadUnique(db, "Tarjeta", []byte(strconv.Itoa(tarjeta.Nrocliente)))
	// fmt.Printf("%s\n", consulta)
}

func cargarComercio(nrocomercio int, nombre string, domicilio string, codigopostal string, telefono string) {
	comercio := Comercio{nrocomercio, nombre, domicilio, codigopostal, telefono}

	data, err := json.Marshal(comercio)
	if err != nil {
		log.Fatal(err)
	}

	CreateUpdate(db, "Comercio", []byte(strconv.Itoa(comercio.Nrocomercio)), data)

	// consulta, err := ReadUnique(db, "Comercio", []byte(strconv.Itoa(comercio.Nrocomercio)))
	// fmt.Printf("%s\n", consulta)
}

func cargarCompra(nrooperacion int, nrotarjeta string, nrocomercio int, fecha string, monto int, pagado bool) {
	compra := Compra{nrooperacion, nrotarjeta, nrocomercio, fecha, monto, pagado}

	data, err := json.Marshal(compra)
	if err != nil {
		log.Fatal(err)
	}

	CreateUpdate(db, "Compra", []byte(strconv.Itoa(compra.Nrooperacion)), data)

	// consulta, err := ReadUnique(db, "Compra", []byte(strconv.Itoa(compra.Nrooperacion)))
	// fmt.Printf("%s\n", consulta)
}

func CreateUpdate(db *bolt.DB, bucketName string, key []byte, value []byte) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))

	err = b.Put(key, value)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func ReadUnique(db *bolt.DB, bucketName string, key []byte) ([]byte, error) {
	var buf []byte

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		buf = b.Get(key)
		return nil
	})

	return buf, err
}
