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
											telefono char(12));
		create table tarjeta (nrotarjeta char(16),
											nrocliente int,
											validadesde char(6),
											validahasta char(6),
											codseguridad char(4),
											limitecompra decimal(8,2),
											estado char(10));
		create table comercio (nrocomercio int,
											nombre text,
											domicilio text,
											codigopostal text,
											telefono char(12));
		create table compra (nrooperacion int,
											nrotarjeta char(16),
											nrocomercio int,
											fecha timestamp,
											monto decimal(7,2),
											pagado bool);
		create table rechazo (nrorechazo int,
											nrotarjeta char(16),
											nrocomercio int,
											fecha timestamp,
											monto decimal(7,2),
											motivo text
											);
		create table cierre (año int,
											mes int,
											terminacion int,
											fechainicio date,
											fechacierre date,
											fechavto date
											);
		create table cabecera(nroresumen int,
											nombre text,
											apellido text,
											domicilio text,
											nrotarjeta char(16),
											desde date,
											hasta date,
											vence date,
											total decimal(8,2)
											);
		create table detalle(nroresumen int,
											nrolinea int,
											fecha date,
											nombrecomercio text,
											monto decimal(7,2)
											);
		create table alerta (nroalerta int,
											nrotarjeta char(16),
											fecha timestamp,
											nrorechazo int,
											codalerta int,
											descripcion text
											);
		create table consumo(nrotarjeta char(16),
											codseguridad char(4),
											nrocomercio int,
											monto decimal(7,2)
											)`)
	if err != nil {
		log.Fatal(err)
	}
}

func CrearPKyFK(){
	crearPK();
	crearFK();
}

func EliminarPKyFK(){
	eliminarFK();
	eliminarPK();
}

func crearPK(){
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

func crearFK() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`alter table tarjeta add constraint tarjeta_nrocliente_fk foreign key (nrocliente) references cliente(nrocliente);
					  alter table compra add constraint compra_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table compra add constraint compra_nrocomercio_fk foreign key (nrocomercio) references comercio(nrocomercio);
					  alter table compra add constraint rechazo_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table compra add constraint rechazo_nrocomercio_fk foreign key (nrocomercio) references comercio(nrocomercio);
					  alter table compra add constraint cabecera_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table compra add constraint alerta_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table compra add constraint consumo_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table compra add constraint consumo_nrocomercio_fk foreign key (nrocomercio) references comercio(nrocomercio);`)	

    if err != nil {
        log.Fatal(err)
    }
	
}

func eliminarPK(){
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`alter table cliente drop constraint cliente_pk;
					  alter table tarjeta drop constraint tarjeta_pk;
					  alter table comercio drop constraint comercio_pk;
	                  alter table compra drop constraint compra_pk;
	                  alter table rechazo drop constraint rechazo_pk;
	                  alter table cierre drop constraint cierre_pk;
	                  alter table cabecera drop constraint cabecera_pk;
	                  alter table detalle drop constraint detalle_pk;
	                  alter table alerta drop constraint alerta_pk;
	                  alter table consumo drop constraint consumo_pk;`)	

	if err != nil {
		log.Fatal(err)
	}
}

func eliminarFK() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`alter table tarjeta drop constraint tarjeta_nrocliente_fk;
					  alter table compra drop constraint compra_nrotarjeta_fk;
					  alter table compra drop constraint compra_nrocomercio_fk;
					  alter table compra drop constraint rechazo_nrotarjeta_fk;
					  alter table compra drop constraint rechazo_nrocomercio_fk;
					  alter table compra drop constraint cabecera_nrotarjeta_fk;
					  alter table compra drop constraint alerta_nrotarjeta_fk;
					  alter table compra drop constraint consumo_nrotarjeta_fk;
					  alter table compra drop constraint consumo_nrocomercio_fk;`)	

    if err != nil {
        log.Fatal(err)
    }
	
}



