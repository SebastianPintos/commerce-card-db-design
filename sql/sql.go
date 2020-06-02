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

	_, _err = _db.Exec(`create database test2`)
	logErr(_err)
}

func CrearTablas() {
	_, err = db.Exec(`DROP SCHEMA public CASCADE`)
	logErr(err)

	_, err = db.Exec(`CREATE SCHEMA public`)
	logErr(err)

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
		create table compra (nrooperacion serial,
											nrotarjeta char(16),
											nrocomercio int,
											fecha timestamp,
											monto decimal(7,2),
											pagado bool);
		create table rechazo (nrorechazo serial,
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
		create table cabecera(nroresumen serial,
											nombre text,
											apellido text,
											domicilio text,
											nrotarjeta char(16),
											desde date,
											hasta date,
											vence date,
											total decimal(8,2)
											);
		create table detalle(nroresumen serial,
											nrolinea int,
											fecha date,
											nombrecomercio text,
											monto decimal(7,2)
											);
		create table alerta (nroalerta serial,
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
	_, err = db.Exec(`alter table cliente add constraint cliente_pk primary key (nrocliente);
					  alter table tarjeta add constraint tarjeta_pk primary key (nrotarjeta);
					  alter table comercio add constraint comercio_pk primary key (nrocomercio);
	                  alter table compra add constraint compra_pk primary key (nrooperacion);
	                  alter table rechazo add constraint rechazo_pk primary key (nrorechazo);
	                  alter table cierre add constraint cierre_pk primary key (año, mes, terminacion);
	                  alter table cabecera add constraint cabecera_pk primary key (nroresumen);
	                  alter table detalle add constraint detalle_pk primary key (nroresumen, nrolinea);
					  alter table alerta add constraint alerta_pk primary key (nroalerta);`)
	logErr(err)
}

func crearFK() {
	_, err = db.Exec(`alter table tarjeta add constraint tarjeta_nrocliente_fk foreign key (nrocliente) references cliente(nrocliente);
					  alter table compra add constraint compra_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table compra add constraint compra_nrocomercio_fk foreign key (nrocomercio) references comercio(nrocomercio);
					  alter table rechazo add constraint rechazo_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table rechazo add constraint rechazo_nrocomercio_fk foreign key (nrocomercio) references comercio(nrocomercio);
					  alter table cabecera add constraint cabecera_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					  alter table detalle add constraint detalle_cabecera_fk foreign key (nroresumen) references cabecera(nroresumen);
					  alter table alerta add constraint alerta_nrotarjeta_fk foreign key (nrotarjeta) references tarjeta(nrotarjeta);
					`)
	logErr(err)
}

func eliminarPK() {
	_, err = db.Exec(`alter table cliente drop constraint cliente_pk;
					  alter table tarjeta drop constraint tarjeta_pk;
					  alter table comercio drop constraint comercio_pk;
	                  alter table compra drop constraint compra_pk;
	                  alter table rechazo drop constraint rechazo_pk;
	                  alter table cierre drop constraint cierre_pk;
	                  alter table cabecera drop constraint cabecera_pk;
	                  alter table detalle drop constraint detalle_pk;
	                  alter table alerta drop constraint alerta_pk;
	                `)
	logErr(err)
}

func eliminarFK() {
	_, err = db.Exec(`alter table tarjeta drop constraint tarjeta_nrocliente_fk;
					  alter table compra drop constraint compra_nrotarjeta_fk;
					  alter table compra drop constraint compra_nrocomercio_fk;
					  alter table rechazo drop constraint rechazo_nrotarjeta_fk;
					  alter table rechazo drop constraint rechazo_nrocomercio_fk;
					  alter table cabecera drop constraint cabecera_nrotarjeta_fk;
					  alter table detalle drop constraint detalle_cabecera_fk;
					  alter table alerta drop constraint alerta_nrotarjeta_fk;`)
	logErr(err)
}

func CargarDatos() {
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

					  insert into tarjeta values(4000001234567899,11348773, 201508, 202008, 733 ,50000,'vigente');
					  insert into tarjeta values(4037001554363655,12349972, 201507, 202007, 332 ,55000,'vigente');
					  insert into tarjeta values(4000001355435322,22648991, 201507, 202007, 201 ,60000,'vigente');
					  insert into tarjeta values(4032011233774494,11341003, 201509, 202009, 204 ,120000,'vigente');
					  insert into tarjeta values(4035055234867402,51558783, 201510, 202010, 108 ,150000,'vigente');
					  insert into tarjeta values(4060001234507040,21347800, 201510, 202010, 909 ,110000,'vigente');
					  insert into tarjeta values(4040071730767070,11448979, 201704, 202204, 810 ,57000,'vigente');
					  insert into tarjeta values(4032002224865843,44349773, 201704, 202204, 327 ,64000,'suspendida');
					  insert into tarjeta values(4034006634262869,33348679, 201708, 202208, 097 ,100000,'suspendida');
					  insert into tarjeta values(4034001232557669,25348533, 201708, 202208, 653 ,140000,'suspendida');
					  insert into tarjeta values(4032002134557009,12228777, 201801, 202301, 070 ,150000,'vigente');
					  insert into tarjeta values(4033002233062344,32680014, 201801, 202301, 202,90000,'anulada');
					  insert into tarjeta values(4000006877865030,21545800, 201801, 202301, 115 ,80000,'vigente');
					  insert into tarjeta values(4000001223567822,23679022, 201604, 202104, 559 ,70000,'vigente');
					  insert into tarjeta values(4000001244532899,12795452, 201604, 202104, 842 ,59000,'vigente');
					  insert into tarjeta values(4032003238867044,11732790, 201602, 202102, 379 ,73000,'vigente');
					  insert into tarjeta values(4000002440217199,29546643, 201601, 202101, 794 ,62000,'vigente');
					  insert into tarjeta values(4032000435566909,18397552, 201701, 202201, 621 ,59000,'suspendida');
					  insert into tarjeta values(4037055274760805,13348765, 201712, 202212, 109 ,69000,'anulada');
					  insert into tarjeta values(4000632234361811,13348765, 201709, 202209, 195 ,53000,'suspendida');
					  insert into tarjeta values(4000000203465800,14348789, 201808, 202308, 290 ,78000,'anulada');
					  insert into tarjeta values(4003300224374894,14348789, 201809, 202309, 284 ,84000,'anulada');

					  `)
	logErr(err)

	_generarCierres()
}

func _generarCierres() {
	generarCierres()

	_, err = db.Query(
		`select generarCierres(2020);`)
	logErr(err)
}

func generarCierres() {
	_, err = db.Query(
		`
		create or replace function generarCierres(año int)returns void as $$
		declare
		  fechainicio text;
		  fechafin text;
		  fechavto text;
		  _mes int;
		begin
		for terminacion in 0..9 loop
			for mes in 1..12 loop
				_mes=mes+1;
				if(mes=12) then
					_mes=1;
				end if;
				if(mes<10 and _mes<10) then
					fechainicio=concat(cast(año as text),'0',cast(mes as text),'01');
					fechafin=concat(cast(año as text),'0',cast(_mes as text),'01');
					fechavto=concat(cast(año as text),'0',cast(_mes as text),'15');
				end if;
				if(mes>=10 and _mes>=10) then
					fechainicio=concat(cast(año as text),cast(mes as text),'01');
					fechafin=concat(cast(año as text),cast(_mes as text),'01');
					fechavto=concat(cast(año as text),cast(_mes as text),'15');
				end if;
				if(mes>=10 and _mes<10) then
					fechainicio=concat(cast(año as text),cast(mes as text),'01');
					fechafin=concat(cast(año as text),cast(_mes as text),'01');
					fechavto=concat(cast(año as text),'0',cast(_mes as text),'15');
				end if;

				insert into cierre values(año, mes, terminacion, to_date(fechainicio,'YYYYMMDD'), to_date(fechafin,'YYYYMMDD'), to_date(fechavto,'YYYYMMDD'));

			end loop;
		end loop;

		end;

		$$ language plpgsql;`)

	logErr(err)
}

func autorizarCompra() {
	agregarRechazo()

	_, err = db.Query(
		`create or replace function autorizarcompra(_nrotarjeta char(16),_codseguridad char(4),_nrocomercio int, _monto decimal(7,2)) returns bool as $$
		 declare
			totalpendiente decimal(7,2);
			montomaximo decimal(8,2);
			fechaVenceTarjeta int;
			fechaVence date;

		 begin

			perform * from tarjeta where nrotarjeta=_nrotarjeta and estado='suspendida';

			if (found) then
				perform agregarrechazo(cast(_nrotarjeta as char(16)),cast(_nrocomercio as int),cast(current_timestamp as timestamp),cast(_monto as decimal(7,2)),cast('La tarjeta se encuentra suspendida' as text));
				return False;
			end if;

			perform * from tarjeta where nrotarjeta=_nrotarjeta and estado='vigente';

			if (not found) then
				perform agregarrechazo(cast(_nrotarjeta as char(16)),cast(_nrocomercio as int),cast(current_timestamp as timestamp),cast(_monto as decimal(7,2)),cast('Tarjeta no válida' as text));
				return False;
			end if;

			perform * from tarjeta where nrotarjeta=_nrotarjeta and codseguridad=_codseguridad;

			if (not found) then
				perform agregarrechazo(cast(_nrotarjeta as char(16)),cast(_nrocomercio as int),cast(current_timestamp as timestamp),cast(_monto as decimal(7,2)),cast('Número de seguridad inválido' as text));
				return False;
			end if;

			totalpendiente:= (select sum(monto) from compra where nrotarjeta =_nrotarjeta and pagado=False);
			montomaximo:= (select limitecompra from tarjeta where nrotarjeta=_nrotarjeta);

			if(totalpendiente is null and _monto > montomaximo or totalpendiente is not null and totalpendiente + _monto>montomaximo) then
				perform agregarrechazo(cast(_nrotarjeta as char(16)),cast(_nrocomercio as int),cast(current_timestamp as timestamp),cast(_monto as decimal(7,2)),cast('Supera límite de tarjeta' as text));
				return False;
			end if;

			select validahasta into fechaVenceTarjeta from tarjeta where nrotarjeta=_nrotarjeta;

			select into FechaVence to_date(fechaVenceTarjeta ||'01','YYYYMMDD');
			select into FechaVence (FechaVence +  interval '1 month')::date;

			if (FechaVence < current_date) then
			perform agregarrechazo(cast(_nrotarjeta as char(16)),cast(_nrocomercio as int),cast(current_timestamp as timestamp),cast(_monto as decimal(7,2)),cast('Plazo de vigencia expirado' as text));
				return False;
			end if;

			insert into compra(nrotarjeta, nrocomercio, fecha, monto, pagado) values( _nrotarjeta, _nrocomercio, current_timestamp, _monto,False);
			return True;

		end;
	$$ language plpgsql;`)
	logErr(err)
}

func GenerarLogicaConsumo() {
	autorizarCompra()
	crearTriggerRechazo()
	crearTriggerConsumo()
	generarConsumo()
}

func generarConsumo() {
	_, err = db.Query(
		`
		create or replace function generarConsumo(cantidad int)returns void as $$
		declare
		  tarjetaAleatoria record;
		  comercioAleatorio int;
		  montoAleatorio decimal(7,2);
		begin

		for _consumo in 0..cantidad-1 loop
			select into montoAleatorio ((random() * (80000 - 100)) + 100) as aleatorio;
			select into comercioAleatorio nrocomercio from comercio order by random() limit 1;
			select into tarjetaAleatoria * from tarjeta order by random() limit 1;
			insert into consumo values(tarjetaAleatoria.nrotarjeta, tarjetaAleatoria.codseguridad, comercioAleatorio, montoAleatorio);
		end loop;
		end;

		$$ language plpgsql;`)
	logErr(err)
}

func crearTriggerConsumo() {
	agregarTestConsumo()

	_, err = db.Query(
		`create trigger agregarconsumo_trg
		before insert on consumo

		for each row
			execute procedure testear_consumo();

		`)
	logErr(err)
}

func agregarTestConsumo() {
	_, err = db.Query(
		`create or replace function testear_consumo() returns trigger as $$
		begin

		perform autorizarcompra(new.nrotarjeta,new.codseguridad, new.nrocomercio,new.monto);
		return new;
		end;

	$$ language plpgsql;`)
	logErr(err)
}

func agregarRechazo() {
	_, err = db.Query(
		`create or replace function agregarrechazo(_nrotarjeta char(16),_nrocomercio int, _fecha timestamp,_monto decimal(7,2),_motivo text) returns void as $$
		declare
			numerorechazo int;

		begin

		insert into rechazo(nrotarjeta, nrocomercio, fecha, monto, motivo) values( _nrotarjeta, _nrocomercio, current_timestamp, _monto, _motivo)
		RETURNING nrorechazo INTO numerorechazo;

		--mover insert rechazo
		select ChequearRechazoLimites(numerorechazo);

		end;

	$$ language plpgsql;`)
	logErr(err)
}

func agregarAlertaRechazo() {
	_, err = db.Query(
		`create or replace function agregar_alerta() returns trigger as $$
		begin

		insert into alerta(nrotarjeta,fecha,nrorechazo,codalerta,descripcion) values(new.nrotarjeta, new.fecha, new.nrorechazo, 0, new.motivo);

		return new;
		end;

	$$ language plpgsql;`)
	logErr(err)
}

func crearTriggerRechazo() {
	agregarAlertaRechazo()

	_, err = db.Query(
		`create trigger agregarrechazo_trg
		before insert on rechazo

		for each row
			execute procedure agregar_alerta();

		`)
	logErr(err)
}
func CrearTriggersSeguridad() {
	seguridadCompras()
	_, err = db.Query(
		`create trigger compras_lapso_tiempo
		before insert on compra

		for each row
			execute procedure compras_lapso_tiempo();
		`)
	logErr(err)
}

func seguridadCompras() {
	_, err = db.Query(
		`create or replace function compras_lapso_tiempo() returns trigger as $$
		declare
			ultimaCompra record;
			difTimestamps decimal;
			codPostalAnterior int;
			codPostalActual int;
		begin
			select * into ultimaCompra from compra where nrotarjeta = new.nrotarjeta order by nrooperacion desc limit 1;

			if(not found) then
				raise notice 'No hay compra anterior';
				return new;
			end if;

			select into difTimestamps extract(epoch from new.fecha - ultimaCompra.fecha) / 60;

			select codigopostal into codPostalAnterior from comercio where nrocomercio = ultimaCompra.nrocomercio;
			select codigopostal into codPostalActual from comercio where nrocomercio = new.nrocomercio;

			if(difTimestamps < 1 and ultimaCompra.nrocomercio != new.nrocomercio and codPostalAnterior = codPostalActual) then
				raise notice 'Alerta compra en menos de 1 minuto en una misma zona';
				return new;
			end if;

			if(difTimestamps < 5 and ultimaCompra.nrocomercio != new.nrocomercio and codPostalAnterior != codPostalActual) then
				raise notice 'Alerta compra en menos de 5 minutos en diferentes zonas';
				return new;
			end if;

			return new;
			end;
		$$ language plpgsql;
	`)
	logErr(err)
}

func GenerarResumen() {
	_, err = db.Query(
		`create or replace function generarResumen(cliente int, anioR int, mesR int) returns bool as $$
		declare
			   idResumen int;
			   totalPagar decimal(7,2) := 0;

				begin
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
				end if;

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
				end if;

		-- Actualizar Resumen
				totalPagar := (SELECT SUM(monto)
							  FROM detalle
							  WHERE nroresumen = idResumen
							  GROUP BY nroresumen);
				UPDATE cabecera
				set total = totalPagar
				WHERE nroresumen = idResumen;

				return True;

				   end;
		$$ language plpgsql;`)
	logErr(err)
}

func ChequearRechazoLimites() {
	_, err = db.Query(
		`create or replace function ChequearRechazoLimites(numeroR int) returns void as $$
		Declare
			tarjetaR char(16);
			fechaR timestamp;
		begin

			SELECT nrotarjeta, fecha INTO tarjetaR, fechaR FROM rechazo where nrorechazo = numeroR;

			perform nrotarjeta
			from rechazo
			where nrotarjeta = tarjetaR
			and fecha = fechaR
			and motivo = 'Supera límite de tarjeta'
			group by nrotarjeta
			having count(*) > 1;

			if (found) then
				insert into alerta(nrotarjeta,fecha,nrorechazo,codalerta,descripcion)
				values (tarjetaR, fechaR, numeroR, 23, 'tarjeta suspendida');

				update tarjeta
				Set estado = 'suspendida'
				where nroTarjeta = tarjetaR;
			end if;

		end;
		$$ language plpgsql;
		`)
	logErr(err)
}
