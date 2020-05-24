package sql

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func CrearDB() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`create database test`)
	if err != nil {
		log.Fatal(err)
	}
}

func CrearTablas() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`DROP SCHEMA public CASCADE`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE SCHEMA public`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create table cliente (nrocliente int,
											nombre text,
											apellido text,
											domicilio text,
											telefono char(12))`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create table tarjeta (nrotarjeta char(16),
											nrocliente int,
											validadesde char(6),
											validahasta char(6),
											codseguridad char(4),
											limitecompra decimal(8,2),
											estado char(10))`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create table comercio (nrocomercio int,
											nombre text,
											domicilio text,
											codigopostal text,
											telefono char(12))`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create table compra (nrooperacion int,
											nrotarjeta char(16),
											nrocomercio int,
											fecha timestamp,
											monto decimal(7,2),
											pagado bool)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create table rechazo (nrorechazo int,
											nrotarjeta char(16),
											nrocomercio int,
											fecha timestamp,
											monto decimal(7,2),
											motivo text
											)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create table cierre (año int,
											mes int,
											terminacion int,
											fechainicio date,
											fechacierre date,
											fechavto date
											)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create table cabecera(nroresumen int,
											nombre text,
											apellido text,
											domicilio text,
											nrotarjeta char(16),
											desde date,
											hasta date,
											vence date,
											total decimal(8,2)
											)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create table detalle(nroresumen int,
											nrolinea int,
											fecha date,
											nombrecomercio text,
											monto decimal(7,2)
											)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create table alerta (nroalerta int,
											nrotarjeta char(16),
											fecha timestamp,
											nrorechazo int,
											codalerta int,
											descripcion text
											)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create table consumo(nrotarjeta char(16),
											codseguridad char(4),
											nrocomercio int,
											monto decimal(7,2)
											)`)
	if err != nil {
		log.Fatal(err)
	}
}

func CrearPK(){
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`alter table cliente add constraint cliente_pk primary key (nrocliente);
					  alter table tarjeta add constraint tarjeta_pk primary key (nrotarjeta);
					  alter table comercio add constraint comercio_pk primary key (nrocomercio);
					  alter table compra add constraint compra_pk primary key (nrooperacion);
					  alter table rechazo add constraint rechazo_pk primary key (nrorechazo);
					  alter table cierre add constraint cierre_pk primary key (año, mes, terminacion);
					  alter table cabecera add constraint cabecera_pk primary key (nroresumen);
					  alter table detalle add constraint detalle_pk primary key (nroresumen, nrolinea);
					  alter table alerta add constraint alerta_pk primary key (nroalerta);
					  alter table consumo add constraint consumo_pk primary key (nrotarjeta);`)	

    if err != nil {
        log.Fatal(err)
    }
	
}

func CrearFK(){
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}

