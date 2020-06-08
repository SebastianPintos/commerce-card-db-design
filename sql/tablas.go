package sql

func crearTablas() {
	_, err = db.Exec(`DROP SCHEMA public CASCADE`)
	logErr(err)

	_, err = db.Exec(`CREATE SCHEMA public`)
	logErr(err)

	_, err = db.Exec(`CREATE TABLE cliente (nrocliente int,
											nombre text,
											apellido text,
											domicilio text,
											telefono char(12));
					CREATE TABLE tarjeta (nrotarjeta char(16),
											nrocliente int,
											validadesde char(6),
											validahasta char(6),
											codseguridad char(4),
											limitecompra decimal(8,2),
											estado char(10));
					CREATE TABLE comercio (nrocomercio int,
											nombre text,
											domicilio text,
											codigopostal text,
											telefono char(12));
					CREATE TABLE compra (nrooperacion serial,
											nrotarjeta char(16),
											nrocomercio int,
											fecha timestamp,
											monto decimal(7,2),
											pagado bool);
					CREATE TABLE rechazo (nrorechazo serial,
											nrotarjeta char(16),
											nrocomercio int,
											fecha timestamp,
											monto decimal(7,2),
											motivo text
											);
					CREATE TABLE cierre (año int,
											mes int,
											terminacion int,
											fechainicio date,
											fechacierre date,
											fechavto date
											);
					CREATE TABLE cabecera(nroresumen serial,
											nombre text,
											apellido text,
											domicilio text,
											nrotarjeta char(16),
											desde date,
											hasta date,
											vence date,
											total decimal(8,2)
											);
					CREATE TABLE detalle(nroresumen serial,
											nrolinea int,
											fecha date,
											nombrecomercio text,
											monto decimal(7,2)
											);
					CREATE TABLE alerta (nroalerta serial,
											nrotarjeta char(16),
											fecha timestamp,
											nrorechazo int,
											codalerta int,
											descripcion text
											);
					CREATE TABLE consumo(nrotarjeta char(16),
											codseguridad char(4),
											nrocomercio int,
											monto decimal(7,2)
											)`)
	logErr(err)
}

func crearPK() {
	_, err = db.Exec(`ALTER TABLE cliente ADD CONSTRAINT cliente_pk PRIMARY KEY (nrocliente);
					  ALTER TABLE tarjeta ADD CONSTRAINT tarjeta_pk PRIMARY KEY (nrotarjeta);
					  ALTER TABLE comercio ADD CONSTRAINT comercio_pk PRIMARY KEY (nrocomercio);
	                  ALTER TABLE compra ADD CONSTRAINT compra_pk PRIMARY KEY (nrooperacion);
	                  ALTER TABLE rechazo ADD CONSTRAINT rechazo_pk PRIMARY KEY (nrorechazo);
	                  ALTER TABLE cierre ADD CONSTRAINT cierre_pk PRIMARY KEY (año, mes, terminacion);
	                  ALTER TABLE cabecera ADD CONSTRAINT cabecera_pk PRIMARY KEY (nroresumen);
	                  ALTER TABLE detalle ADD CONSTRAINT detalle_pk PRIMARY KEY (nroresumen, nrolinea);
					  ALTER TABLE alerta ADD CONSTRAINT alerta_pk PRIMARY KEY (nroalerta);`)
	logErr(err)
}

func crearFK() {
	_, err = db.Exec(`ALTER TABLE tarjeta ADD CONSTRAINT tarjeta_nrocliente_fk FOREIGN KEY (nrocliente) REFERENCES cliente(nrocliente);
					  ALTER TABLE compra ADD CONSTRAINT compra_nrotarjeta_fk FOREIGN KEY (nrotarjeta) REFERENCES tarjeta(nrotarjeta);
					  ALTER TABLE compra ADD CONSTRAINT compra_nrocomercio_fk FOREIGN KEY (nrocomercio) REFERENCES comercio(nrocomercio);
					  ALTER TABLE rechazo ADD CONSTRAINT rechazo_nrotarjeta_fk FOREIGN KEY (nrotarjeta) REFERENCES tarjeta(nrotarjeta);
					  ALTER TABLE rechazo ADD CONSTRAINT rechazo_nrocomercio_fk FOREIGN KEY (nrocomercio) REFERENCES comercio(nrocomercio);
					  ALTER TABLE cabecera ADD CONSTRAINT cabecera_nrotarjeta_fk FOREIGN KEY (nrotarjeta) REFERENCES tarjeta(nrotarjeta);
					  ALTER TABLE detalle ADD CONSTRAINT detalle_cabecera_fk FOREIGN KEY (nroresumen) REFERENCES cabecera(nroresumen);
					  ALTER TABLE alerta ADD CONSTRAINT alerta_nrotarjeta_fk FOREIGN KEY (nrotarjeta) REFERENCES tarjeta(nrotarjeta);
					`)
	logErr(err)
}

func eliminarPK() {
	_, err = db.Exec(`ALTER TABLE cliente DROP CONSTRAINT cliente_pk;
					  ALTER TABLE tarjeta DROP CONSTRAINT tarjeta_pk;
					  ALTER TABLE comercio DROP CONSTRAINT comercio_pk;
	                  ALTER TABLE compra DROP CONSTRAINT compra_pk;
	                  ALTER TABLE rechazo DROP CONSTRAINT rechazo_pk;
	                  ALTER TABLE cierre DROP CONSTRAINT cierre_pk;
	                  ALTER TABLE cabecera DROP CONSTRAINT cabecera_pk;
	                  ALTER TABLE detalle DROP CONSTRAINT detalle_pk;
	                  ALTER TABLE alerta DROP CONSTRAINT alerta_pk;
	                `)
	logErr(err)
}

func eliminarFK() {
	_, err = db.Exec(`ALTER TABLE tarjeta DROP CONSTRAINT tarjeta_nrocliente_fk;
					  ALTER TABLE compra DROP CONSTRAINT compra_nrotarjeta_fk;
					  ALTER TABLE compra DROP CONSTRAINT compra_nrocomercio_fk;
					  ALTER TABLE rechazo DROP CONSTRAINT rechazo_nrotarjeta_fk;
					  ALTER TABLE rechazo DROP CONSTRAINT rechazo_nrocomercio_fk;
					  ALTER TABLE cabecera DROP CONSTRAINT cabecera_nrotarjeta_fk;
					  ALTER TABLE detalle DROP CONSTRAINT detalle_cabecera_fk;
					  ALTER TABLE alerta DROP CONSTRAINT alerta_nrotarjeta_fk;`)
	logErr(err)
}
