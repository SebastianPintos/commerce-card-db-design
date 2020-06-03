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

func CrearTablas() {
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

func CrearPKyFK() {
	crearPK()
	crearFK()
}

func EliminarPKyFK() {
	eliminarFK()
	eliminarPK()
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

func CargarDatos() {
	_, err = db.Exec(`INSERT INTO cliente VALUES(11348773,'Rocío', 'Losada','Av. Presidente Perón 1530',1151102983);
					  INSERT INTO cliente VALUES(12349972,'María Estela', 'Martínez','Belgrano 1830',1150006655);
					  INSERT INTO cliente VALUES(22648991,'Laura', 'Santos','Italia 220',1153399452);
					  INSERT INTO cliente VALUES(11341003,'Graciela', 'Chasco','Tribulato 1340',1258579091);
					  INSERT INTO cliente VALUES(51558783,'Gabriela', 'Troncoso','Muñoz 1820',112234667);
					  INSERT INTO cliente VALUES(21347800,'Marta', 'Carbajo','San Luis 873',111998340);
					  INSERT INTO cliente VALUES(11448979,'Belén', 'Ferraris','Echeverría 780',113229087);
					  INSERT INTO cliente VALUES(44349773,'Abril', 'Hernández','Av. Sourdeaux 1700',115598342);
					  INSERT INTO cliente VALUES(33348679,'Sofía', 'Godoy','Av. SenadOR Morón 1221',114558004);
					  INSERT INTO cliente VALUES(25348533,'Adriana', 'Golluscio','Misiones 725',112111558);
					  INSERT INTO cliente VALUES(12228777,'Juan Carlos', 'Leguizamon','Serrano 120',1151101182);
					  INSERT INTO cliente VALUES(32680014,'Alberto', 'Ferrero','Pardo 990',1159944558);
					  INSERT INTO cliente VALUES(21545800,'Roberto', 'Ubertalli','Santa Fé 160',110076548);
					  INSERT INTO cliente VALUES(23679022,'Mario', 'Valdéz','Tucumán 550',116690874);
					  INSERT INTO cliente VALUES(12795452,'Lautaro', 'Flores','Río Diamante 186',113678652);
					  INSERT INTO cliente VALUES(11732790,'Bautista', 'Bello','Río Cuarto 191',111451419);
					  INSERT INTO cliente VALUES(29546643,'Diego', 'Fagnani','Av. Gaspar Campos 122',111009070);
					  INSERT INTO cliente VALUES(18397552,'Pedro', 'Tomarello','Av. San Martín 1511',110887547);
					  INSERT INTO cliente VALUES(13348765,'José', 'Mengarelli','Guido Spano 244',110044332);
					  INSERT INTO cliente VALUES(14348789,'Ricardo', 'Llanos','Corrientes 183',119034572);

					  INSERT INTO comercio VALUES(501,'Kevingston', 'Av. Tte. Gral. Ricchieri 965', 1661 ,46666181);
					  INSERT INTO comercio VALUES(523,'47 street', 'Paunero 1575', 1663 ,47597581);
					  INSERT INTO comercio VALUES(513,'Garbarino', 'Av. Bartolomé Mitre 1198', 1661 ,08104440018);
					  INSERT INTO comercio VALUES(521,'Bella Vista Hogar', 'Av. SenadOR Morón 1094', 1661 ,46661544);
					  INSERT INTO comercio VALUES(578,'Panadería y Confitería: La Princesa', 'Av. SenadOR Morón 1200', 1661 ,46681339);
					  INSERT INTO comercio VALUES(564,'FOX', 'Av Pres. Juan Domingo Perón 907', 1663 ,46676777);
					  INSERT INTO comercio VALUES(569,'La Pata Loca', 'Av. Moisés Lebensohn 98', 1661 ,46660861);
					  INSERT INTO comercio VALUES(545,'Frávega', 'Av. Pres. Juan Domingo Perón 1127', 1663 ,44512063);
					  INSERT INTO comercio VALUES(543,'Spit Bella Vista', 'Av. SenadOR Morón 1452', 1661 ,1153519765);
					  INSERT INTO comercio VALUES(527,'Óptica Cristal', 'Av. Dr. Ricardo Balbín 1125', 1663 ,46649400);
					  INSERT INTO comercio VALUES(508,'Óptica Mattaldi', 'Av. Mattaldi 1141', 1661 ,46683911);
					  INSERT INTO comercio VALUES(509,'Estancia San Francisco San Miguel', 'Concejal Tribulato 1265', 1663 ,5446676082);
					  INSERT INTO comercio VALUES(500,'Rabelia heladería', 'San José 972', 1663 ,46649352);
					  INSERT INTO comercio VALUES(520,'Heladería Ciwe', 'San José 785', 1663 ,46646003);
					  INSERT INTO comercio VALUES(588,'Rever Pass', 'Paunero 1447,', 1663 ,44513921);
					  INSERT INTO comercio VALUES(582,'Rapsodia', 'Av. Pres. Arturo Umberto Illia 3770', 1663 ,1160911581);
					  INSERT INTO comercio VALUES(530,'Grimoldi', 'Paunero 1415', 1663 ,44517343);
					  INSERT INTO comercio VALUES(596,'Umma', 'Paunero 1476', 1663 ,44519267);
					  INSERT INTO comercio VALUES(538,'COTO', 'Ohiggins 1280', 1661 ,46682636);
					  INSERT INTO comercio VALUES(553,'Disco', 'Av. SenadOR Morón 960', 1661 ,08107778888);

					  INSERT INTO tarjeta VALUES(4000001234567899,11348773, 201508, 202008, 733 ,50000,'vigente');
					  INSERT INTO tarjeta VALUES(4037001554363655,12349972, 201507, 202007, 332 ,55000,'vigente');
					  INSERT INTO tarjeta VALUES(4000001355435322,22648991, 201507, 202007, 201 ,60000,'vigente');
					  INSERT INTO tarjeta VALUES(4032011233774494,11341003, 201509, 202009, 204 ,120000,'vigente');
					  INSERT INTO tarjeta VALUES(4035055234867402,51558783, 201510, 202010, 108 ,150000,'vigente');
					  INSERT INTO tarjeta VALUES(4060001234507040,21347800, 201510, 202010, 909 ,110000,'vigente');
					  INSERT INTO tarjeta VALUES(4040071730767070,11448979, 201704, 202204, 810 ,57000,'vigente');
					  INSERT INTO tarjeta VALUES(4032002224865843,44349773, 201704, 202204, 327 ,64000,'suspendida');
					  INSERT INTO tarjeta VALUES(4034006634262869,33348679, 201708, 202208, 097 ,100000,'suspendida');
					  INSERT INTO tarjeta VALUES(4034001232557669,25348533, 201708, 202208, 653 ,140000,'suspendida');
					  INSERT INTO tarjeta VALUES(4032002134557009,12228777, 201801, 202301, 070 ,150000,'vigente');
					  INSERT INTO tarjeta VALUES(4033002233062344,32680014, 201801, 202301, 202,90000,'anulada');
					  INSERT INTO tarjeta VALUES(4000006877865030,21545800, 201801, 202301, 115 ,80000,'vigente');
					  INSERT INTO tarjeta VALUES(4000001223567822,23679022, 201604, 202104, 559 ,70000,'vigente');
					  INSERT INTO tarjeta VALUES(4000001244532899,12795452, 201604, 202104, 842 ,59000,'vigente');
					  INSERT INTO tarjeta VALUES(4032003238867044,11732790, 201602, 202102, 379 ,73000,'vigente');
					  INSERT INTO tarjeta VALUES(4000002440217199,29546643, 201601, 202101, 794 ,62000,'vigente');
					  INSERT INTO tarjeta VALUES(4032000435566909,18397552, 201701, 202201, 621 ,59000,'suspendida');
					  INSERT INTO tarjeta VALUES(4037055274760805,13348765, 201712, 202212, 109 ,69000,'anulada');
					  INSERT INTO tarjeta VALUES(4000632234361811,13348765, 201709, 202209, 195 ,53000,'suspendida');
					  INSERT INTO tarjeta VALUES(4000000203465800,14348789, 201808, 202308, 290 ,78000,'anulada');
					  INSERT INTO tarjeta VALUES(4003300224374894,14348789, 201809, 202309, 284 ,84000,'anulada');

					  `)
	logErr(err)

	_generarCierres()
}

func _generarCierres() {
	generarCierres()

	_, err = db.Query(
		`SELECT generarCierres(2020);`)
	logErr(err)
}

func generarCierres() {
	_, err = db.Query(
		`
		CREATE OR REPLACE FUNCTION generarCierres(año int)returns void as $$
		DECLARE
		  fechainicio text;
		  fechafin text;
		  fechavto text;
		  _mes int;
		BEGIN
		FOR terminacion in 0..9 LOOP
			FOR mes in 1..12 LOOP
				_mes=mes+1;
				if(mes=12) THEN
					_mes=1;
				END IF;
				if(mes<10 and _mes<10) THEN
					fechainicio=concat(CAST(año as text),'0',CAST(mes as text),'01');
					fechafin=concat(CAST(año as text),'0',CAST(_mes as text),'01');
					fechavto=concat(CAST(año as text),'0',CAST(_mes as text),'15');
				END IF;
				if(mes>=10 and _mes>=10) THEN
					fechainicio=concat(CAST(año as text),CAST(mes as text),'01');
					fechafin=concat(CAST(año as text),CAST(_mes as text),'01');
					fechavto=concat(CAST(año as text),CAST(_mes as text),'15');
				END IF;
				if(mes>=10 and _mes<10) THEN
					fechainicio=concat(CAST(año as text),CAST(mes as text),'01');
					fechafin=concat(CAST(año as text),CAST(_mes as text),'01');
					fechavto=concat(CAST(año as text),'0',CAST(_mes as text),'15');
				END IF;

				INSERT INTO cierre VALUES(año, mes, terminacion, TO_DATE(fechainicio,'YYYYMMDD'), TO_DATE(fechafin,'YYYYMMDD'), TO_DATE(fechavto,'YYYYMMDD'));

			END LOOP;
		END LOOP;

		END;

		$$ LANGUAGE PLPGSQL;`)

	logErr(err)
}

func autorizarCompra() {
	agregarRechazo()

	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION autorizarcompra(_nrotarjeta char(16),_codseguridad char(4),_nrocomercio int, _monto decimal(7,2)) returns bool as $$
		 DECLARE
			totalpendiente decimal(8,2);
			montomaximo decimal(8,2);
			fechaVenceTarjeta int;
			fechaVence date;

		 BEGIN

			PERFORM * FROM tarjeta WHERE nrotarjeta=_nrotarjeta and estado='suspendida';

			if (found) THEN
				PERFORM agregarrechazo(CAST(_nrotarjeta as char(16)),CAST(_nrocomercio as int),CAST(current_timestamp as timestamp),CAST(_monto as decimal(7,2)),CAST('La tarjeta se encuentra suspendida' as text));
				return False;
			END IF;

			PERFORM * FROM tarjeta WHERE nrotarjeta=_nrotarjeta and estado='vigente';

			if (not found) THEN
				PERFORM agregarrechazo(CAST(_nrotarjeta as char(16)),CAST(_nrocomercio as int),CAST(current_timestamp as timestamp),CAST(_monto as decimal(7,2)),CAST('Tarjeta no válida' as text));
				return False;
			END IF;

			PERFORM * FROM tarjeta WHERE nrotarjeta=_nrotarjeta and codseguridad=_codseguridad;

			if (not found) THEN
				PERFORM agregarrechazo(CAST(_nrotarjeta as char(16)),CAST(_nrocomercio as int),CAST(current_timestamp as timestamp),CAST(_monto as decimal(7,2)),CAST('Número de seguridad inválido' as text));
				return False;
			END IF;

			totalpendiente:= (SELECT sum(monto) FROM compra WHERE nrotarjeta =_nrotarjeta and pagado=False);
			montomaximo:= (SELECT limitecompra FROM tarjeta WHERE nrotarjeta=_nrotarjeta);

			if(totalpendiente is null and _monto > montomaximo OR totalpendiente is not null and totalpendiente + _monto>montomaximo) THEN
				PERFORM agregarrechazo(CAST(_nrotarjeta as char(16)),CAST(_nrocomercio as int),CAST(current_timestamp as timestamp),CAST(_monto as decimal(7,2)),CAST('Supera límite de tarjeta' as text));
				return False;
			END IF;

			SELECT validahasta INTO fechaVenceTarjeta FROM tarjeta WHERE nrotarjeta=_nrotarjeta;

			SELECT INTO FechaVence TO_DATE(fechaVenceTarjeta ||'01','YYYYMMDD');
			SELECT INTO FechaVence (FechaVence +  interval '1 month')::date;

			if (FechaVence < current_date) THEN
			PERFORM agregarrechazo(CAST(_nrotarjeta as char(16)),CAST(_nrocomercio as int),CAST(current_timestamp as timestamp),CAST(_monto as decimal(7,2)),CAST('Plazo de vigencia expirado' as text));
				return False;
			END IF;

			INSERT INTO compra(nrotarjeta, nrocomercio, fecha, monto, pagado) VALUES( _nrotarjeta, _nrocomercio, current_timestamp, _monto,False);
			return True;

		END;
	$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func GenerarLogicaConsumo() {
	autorizarCompra()
	generarConsumo()
	crearTriggerConsumo()
}

func GenerarLogicaAlertas() {
	crearTriggerRechazo()
	crearTriggersSeguridad()
}

func generarConsumo() {
	_, err = db.Query(
		`
		CREATE OR REPLACE FUNCTION generarConsumo(cantidad int)returns void as $$
		DECLARE
		  tarjetaAleatoria record;
		  comercioAleatorio int;
		  montoAleatorio decimal(7,2);
		BEGIN

		FOR _consumo in 0..cantidad-1 LOOP
			montoAleatorio = 999 + random()*99000;
			PERFORM TRUNC(montoAleatorio,2);
			SELECT INTO comercioAleatorio nrocomercio FROM comercio ORDER BY random() LIMIT 1;
			SELECT INTO tarjetaAleatoria * FROM tarjeta ORDER BY random() LIMIT 1;
			INSERT INTO consumo VALUES(tarjetaAleatoria.nrotarjeta, tarjetaAleatoria.codseguridad, comercioAleatorio, CAST(montoAleatorio as decimal(7,2)));
		END LOOP;
		END;

		$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func crearTriggerConsumo() {
	agregarTestConsumo()

	_, err = db.Query(
		`CREATE trigger agregarconsumo_trg
		BEFORE INSERT ON consumo

		FOR EACH ROW
			EXECUTE PROCEDURE testear_consumo();

		`)
	logErr(err)
}

func agregarTestConsumo() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION testear_consumo() returns trigger as $$
		BEGIN
		
		PERFORM autorizarcompra(new.nrotarjeta,new.codseguridad, new.nrocomercio, new.monto);
		return new;
		END;

	$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func agregarRechazo() {
	chequearRechazoLimites()
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION agregarrechazo(_nrotarjeta char(16),_nrocomercio int, _fecha timestamp,_monto decimal(7,2),_motivo text) returns void as $$
		DECLARE
			numerorechazo int;

		BEGIN

		INSERT INTO rechazo(nrotarjeta, nrocomercio, fecha, monto, motivo) VALUES( _nrotarjeta, _nrocomercio, current_timestamp, _monto, _motivo)
		RETURNING nrorechazo INTO numerorechazo;

		--mover INSERT rechazo
		PERFORM ChequearRechazoLimites(numerorechazo);

		END;

	$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func agregarAlertaRechazo() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION agregar_alerta() returns trigger as $$
		BEGIN

		INSERT INTO alerta(nrotarjeta,fecha,nrorechazo,codalerta,descripcion) VALUES(new.nrotarjeta, new.fecha, new.nrorechazo, 0 , new.motivo);

		return new;
		END;

	$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func crearTriggerRechazo() {
	agregarAlertaRechazo()

	_, err = db.Query(
		`CREATE trigger agregarrechazo_trg
		BEFORE INSERT ON rechazo

		FOR EACH ROW
			EXECUTE PROCEDURE agregar_alerta();

		`)
	logErr(err)
}
func crearTriggersSeguridad() {
	seguridadCompras()
	_, err = db.Query(
		`CREATE trigger compras_lapso_tiempo
		BEFORE INSERT ON compra

		FOR EACH ROW
			EXECUTE PROCEDURE compras_lapso_tiempo();
		`)
	logErr(err)
}

func seguridadCompras() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION compras_lapso_tiempo() returns trigger as $$
		DECLARE
			ultimaCompra record;
			difTimestamps decimal;
			codPostalAnterior int;
			codPostalActual int;
		BEGIN
			SELECT * INTO ultimaCompra FROM compra WHERE nrotarjeta = new.nrotarjeta ORDER BY nrooperacion DESC LIMIT 1;

			if(not found) THEN
				return new;
			END IF;

			SELECT INTO difTimestamps EXTRACT(EPOCH FROM new.fecha - ultimaCompra.fecha) / 60;

			SELECT codigopostal INTO codPostalAnterior FROM comercio WHERE nrocomercio = ultimaCompra.nrocomercio;
			SELECT codigopostal INTO codPostalActual FROM comercio WHERE nrocomercio = new.nrocomercio;

			if(difTimestamps < 1 and ultimaCompra.nrocomercio != new.nrocomercio and codPostalAnterior = codPostalActual) THEN
				INSERT INTO alerta(nrotarjeta,fecha,nrorechazo,codalerta,descripcion) VALUES(new.nrotarjeta, new.fecha, -1, 1 , 'Compra en menos de 1 minuto en una misma zona');
				return new;
			END IF;

			if(difTimestamps < 5 and ultimaCompra.nrocomercio != new.nrocomercio and codPostalAnterior != codPostalActual) THEN
				INSERT INTO alerta(nrotarjeta,fecha,nrorechazo,codalerta,descripcion) VALUES(new.nrotarjeta, new.fecha, -1, 5 , 'Compra en menos de 5 minutos en diferentes zonas');
				return new;
			END IF;

			return new;
		END;
		$$ LANGUAGE PLPGSQL;
	`)
	logErr(err)
}

func GenerarResumen() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION generarResumen(cliente int, aniOR int, mesR int) returns bool as $$
		DECLARE
			   idResumen int;
			   totalPagar decimal(8,2) := 0;
			   _linea record;
				
		BEGIN
		-- 	Generar Cabecera
		INSERT INTO cabecera (nombre, apellido, domicilio, nrotarjeta, desde, hasta, vence)
		SELECT cli.nombre, cli.apellido, cli.domicilio, t.nrotarjeta, c.fechainicio, c.fechacierre, c.fechavto
		FROM public.tarjeta t, public.cierre c, public.cliente cli
		WHERE SUBSTRING (t.nrotarjeta, LENGTH(t.nrotarjeta), 1)::int = c.terminacion
		and cli.nrocliente = t.nrocliente
		and t.nrocliente = cliente
		and c.año = anioR
		and c.mes = mesR
		RETURNING nroresumen INTO idResumen;

		if (idResumen is null) then
					raise 'No se pudo generar el resumen, Cliente inexistente';
					return False;
		END IF;

		-- Generar detalle
		INSERT INTO detalle (nroresumen, nrolinea, fecha, nombrecomercio, monto)
		SELECT idResumen, ROW_NUMBER () OVER (ORDER BY t.nrotarjeta) as nrolinea, co.fecha, com.nombre , co.monto
			FROM public.tarjeta t, public.cierre c, public.compra co, public.comercio com
			WHERE SUBSTRING (t.nrotarjeta, LENGTH(t.nrotarjeta), 1)::int = c.terminacion
			and co.nrotarjeta = t.nrotarjeta
			and com.nrocomercio = co.nrocomercio
			and t.nrocliente = cliente
			and c.año = anioR
			and c.mes = mesR
			and co.fecha >= c.fechainicio
			and co.fecha <= c.fechacierre;

		if (lastval() is NULL) then
				raise 'No se pudo generar el resumen';
				return False;
		END IF;

		-- Actualizar Resumen
		totalPagar := (SELECT SUM(monto)
							  FROM detalle
							  WHERE nroresumen = idResumen
							  GROUP BY nroresumen);
		UPDATE cabecera
			SET total = totalPagar WHERE nroresumen = idResumen;

		--Cambiar pagado a True
		for _linea in SELECT * FROM public.tarjeta t, public.cierre c, public.compra co, public.comercio com
			WHERE SUBSTRING (t.nrotarjeta, LENGTH(t.nrotarjeta), 1)::int = c.terminacion
			and co.nrotarjeta = t.nrotarjeta
			and com.nrocomercio = co.nrocomercio
			and t.nrocliente = cliente loop

					UPDATE compra set pagado = True where nrotarjeta=_linea.nrotarjeta and monto=_linea.monto;				
				
		END LOOP;	
				
		return True;

		END;
		$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func chequearRechazoLimites() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION ChequearRechazoLimites(numerOR int) returns void as $$
		DECLARE
			tarjetaR char(16);
			fechaR timestamp;
		BEGIN

			SELECT nrotarjeta, fecha INTO tarjetaR, fechaR FROM rechazo WHERE nrorechazo = numeroR;

			PERFORM nrotarjeta
			FROM rechazo
			WHERE nrotarjeta = tarjetaR
			and fecha = fechaR
			and motivo = 'Supera límite de tarjeta'
			GROUP BY nrotarjeta
			HAVING COUNT(*) > 1;

			if (found) THEN
				INSERT INTO alerta(nrotarjeta,fecha,nrorechazo,codalerta,descripcion)
				VALUES (tarjetaR, fechaR, numeroR, 32, 'Tarjeta suspendida');

				UPDATE tarjeta
				SET estado = 'suspendida'
				WHERE nroTarjeta = tarjetaR;
			END IF;

		END;
		$$ LANGUAGE PLPGSQL;
		`)
	logErr(err)
}
