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
					  alter table rechazo add constraint rechazo_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table rechazo add constraint rechazo_nrocomercio_fk foreign key (nrocomercio) references comercio(nrocomercio);
					  alter table cabecera add constraint cabecera_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table alerta add constraint alerta_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table consumo add constraint consumo_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table consumo add constraint consumo_nrocomercio_fk foreign key (nrocomercio) references comercio(nrocomercio);`)	

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

func CargarDatos() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`insert into cliente values(11348773,'Rocío', 'Losada','Av. Presidente Perón 1530',1151102983);
					  insert into cliente values(12349972,'María Estela', 'Martínez','Belgrano 1830',1150006655);
					  insert into cliente values(22648991,'Laura', 'Santos','Italia 220',1153399452);
					  insert into cliente values(11341003,'Graciela', 'Chasco','Tribulato 1340',1258579091);
					  insert into cliente values(51558783,'Gabriela', 'Troncoso','Muñoz 1820',112234667);
					  insert into cliente values(21347800,'Marta', 'Carbajo','San Luis 873',111998340);
					  insert into cliente values(11448979,'Belén', 'Ferraris','Echeverría 780',113229087);
					  insert into cliente values(44349773,'Abril', 'Hernández','Av. Sourdeaux 1700',115598342);
					  insert into cliente values(33348679,'Sofía', 'Godoy','Av. Senador Morón 1221',114558004);
					  insert into cliente values(25348533,'Adriana', 'Golluscio','Misiones 725',112111558);
					  insert into cliente values(12228777,'Juan Carlos', 'Leguizamon','Serrano 120',1151101182);
					  insert into cliente values(32680014,'Alberto', 'Ferrero','Pardo 990',1159944558);
					  insert into cliente values(21545800,'Roberto', 'Ubertalli','Santa Fé 160',110076548);
					  insert into cliente values(23679022,'Mario', 'Valdéz','Tucumán 550',116690874);
					  insert into cliente values(12795452,'Lautaro', 'Flores','Río Diamante 186',113678652);
					  insert into cliente values(11732790,'Bautista', 'Bello','Río Cuarto 191',111451419);
					  insert into cliente values(29546643,'Diego', 'Fagnani','Av. Gaspar Campos 122',111009070);
					  insert into cliente values(18397552,'Pedro', 'Tomarello','Av. San Martín 1511',110887547);
					  insert into cliente values(13348765,'José', 'Mengarelli','Guido Spano 244',110044332);
					  insert into cliente values(14348789,'Ricardo', 'Llanos','Corrientes 183',119034572);
					  
					  insert into comercio values(501,'Kevingston', 'Av. Tte. Gral. Ricchieri 965', 1661 ,46666181);
					  insert into comercio values(523,'47 street', 'Paunero 1575', 1663 ,47597581);
					  insert into comercio values(513,'Garbarino', 'Av. Bartolomé Mitre 1198', 1661 ,08104440018);
					  insert into comercio values(521,'Bella Vista Hogar', 'Av. Senador Morón 1094', 1661 ,46661544);
					  insert into comercio values(578,'Panadería y Confitería: La Princesa', 'Av. Senador Morón 1200', 1661 ,46681339);
					  insert into comercio values(564,'FOX', 'Av Pres. Juan Domingo Perón 907', 1663 ,46676777);
					  insert into comercio values(569,'La Pata Loca', 'Av. Moisés Lebensohn 98', 1661 ,46660861);
					  insert into comercio values(545,'Frávega', 'Av. Pres. Juan Domingo Perón 1127', 1663 ,44512063);
					  insert into comercio values(543,'Spit Bella Vista', 'Av. Senador Morón 1452', 1661 ,1153519765);
					  insert into comercio values(527,'Óptica Cristal', 'Av. Dr. Ricardo Balbín 1125', 1663 ,46649400);
					  insert into comercio values(508,'Óptica Mattaldi', 'Av. Mattaldi 1141', 1661 ,46683911);
					  insert into comercio values(509,'Estancia San Francisco San Miguel', 'Concejal Tribulato 1265', 1663 ,5446676082);
					  insert into comercio values(500,'Rabelia heladería', 'San José 972', 1663 ,46649352);
					  insert into comercio values(520,'Heladería Ciwe', 'San José 785', 1663 ,46646003);
					  insert into comercio values(588,'Rever Pass', 'Paunero 1447,', 1663 ,44513921);
					  insert into comercio values(582,'Rapsodia', 'Av. Pres. Arturo Umberto Illia 3770', 1663 ,1160911581);
					  insert into comercio values(530,'Grimoldi', 'Paunero 1415', 1663 ,44517343);
					  insert into comercio values(596,'Umma', 'Paunero 1476', 1663 ,44519267);
					  insert into comercio values(538,'COTO', 'Ohiggins 1280', 1661 ,46682636);
					  insert into comercio values(553,'Disco', 'Av. Senador Morón 960', 1661 ,08107778888);
					  
					  `)	

    if err != nil {
        log.Fatal(err)
    }
	
}

