package sql

func spGenerarCierres() {
	_, err = db.Query(
		`create or replace function generarCierres(anio int) returns void as $$
		Declare
			fdesde date;
			fhasta date;
			fvto date;
		BEGIN
		
			FOR tarjeta IN 0 .. 9 BY 1
			LOOP
				SELECT into fdesde to_date((anio - 1)::text || '12' || (select 23 + trunc(random() * 4))::text, 'YYYYMMDD');
				SELECT into fhasta fdesde::date + cast((select 29 + trunc(random() * 2))::text || ' days' as interval);
				SELECT into fvto fhasta::date + cast('10 days' as interval);
				
				FOR mes IN 1 .. 12 BY 1
				LOOP			
					insert into cierre values(anio,mes,tarjeta,fdesde,fhasta, fvto);
					
					SELECT into fdesde fhasta::date + cast('1 days' as interval);
					SELECT into fhasta fdesde::date + cast((select 29 + trunc(random() * 2))::text || ' days' as interval);
					SELECT into fvto fhasta::date + cast('10 days' as interval);
				END LOOP;
			END LOOP;
		END;
		$$ language plpgsql;`)

	logErr(err)
}

func spGenerarResumen() {
	_, err = db.Query(
		`create or replace function generarResumen(cliente int, anioR int, mesR int) returns bool as $$
		declare
			   idResumen int;
			   totalPagar decimal(7,2) := 0;
			   _linea record;
			   tarjeta char(16);
		   
				begin		
				
				FOR tarjeta IN 
					Select nrotarjeta
					From Tarjeta
					Where nrocliente = cliente
					and estado = 'vigente'
				LOOP
				
				-- 	Generar Cabecera
					INSERT INTO cabecera (nombre, apellido, domicilio, nrotarjeta, desde, hasta, vence) 
					SELECT cli.nombre, cli.apellido, cli.domicilio, t.nrotarjeta, c.fechainicio, c.fechacierre, c.fechavto
						FROM public.tarjeta t, public.cierre c, public.cliente cli
						WHERE SUBSTRING (t.nrotarjeta, LENGTH(t.nrotarjeta), 1)::int = c.terminacion
						and cli.nrocliente = t.nrocliente
						and t.nrotarjeta = tarjeta
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
					and t.nrotarjeta = tarjeta
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
					set total = COALESCE(NULLIF(totalPagar, 0), 0)
					WHERE nroresumen = idResumen;	
					
				--Cambiar pagado a True
					FOR _linea in SELECT * FROM public.tarjeta t, public.cierre c, public.compra co, public.comercio com
						WHERE SUBSTRING (t.nrotarjeta, LENGTH(t.nrotarjeta), 1)::int = c.terminacion
						and co.nrotarjeta = t.nrotarjeta
						and com.nrocomercio = co.nrocomercio
						and t.nrotarjeta = tarjeta 
					LOOP
						UPDATE compra set pagado = True where nrotarjeta=_linea.nrotarjeta and monto=_linea.monto;									
					END LOOP;
				
				END LOOP;
				
				return True;
				
				   end;
		$$ language plpgsql;`)
	logErr(err)
}

func spGenerarResumenesPeriodo() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION generarResumenesPeriodo(aniOR int, mesR int) returns void as $$
		DECLARE
			   _linea record;
				
		BEGIN
		FOR _linea in 
			SELECT cli.nrocliente
			FROM compra co , cliente cli, tarjeta tr
			WHERE co.nrotarjeta = tr.nrotarjeta
			  AND tr.nrocliente = cli.nrocliente
			  AND co.pagado = False
			GROUP BY cli.nrocliente
		LOOP
				perform generarResumen(_linea.nrocliente, anioR, mesR);
		END LOOP;
		END;
		$$ LANGUAGE PLPGSQL;`)

	logErr(err)

	_, err = db.Query(
		`SELECT generarResumenesPeriodo(2020,06);`)

	logErr(err)
}

func spObtenerDisponible() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION obtenerDisponible(_nrotarjeta char(16))returns decimal(8,2) as $$
		DECLARE
			_limite decimal(8,2);
			_consumos decimal(8,2);
			
		BEGIN	
			SELECT coalesce(sum(monto), 0) INTO _consumos 
			FROM compra
			WHERE nrotarjeta =_nrotarjeta 
			  and pagado = False;
			
			SELECT limitecompra INTO _limite 
			FROM tarjeta 
			WHERE nrotarjeta = _nrotarjeta;
			
			RETURN _limite - _consumos;
		END;	
		
	$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func spChequearRechazoLimites() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION chequearRechazoLimites(numerOR int) returns void as $$
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

func spAutorizarCompra() {
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
				PERFORM agregarrechazo(CAST(_nrotarjeta as char(16)),CAST(_nrocomercio as int),CAST(current_timestamp as timestamp),CAST(_monto as decimal(7,2)),CAST('Código de seguridad inválido' as text));
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

func spAgregarRechazo() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION agregarrechazo(_nrotarjeta char(16),_nrocomercio int, _fecha timestamp,_monto decimal(7,2),_motivo text) returns void as $$
		DECLARE
			numerorechazo int;

		BEGIN

		INSERT INTO rechazo(nrotarjeta, nrocomercio, fecha, monto, motivo) VALUES( _nrotarjeta, _nrocomercio, current_timestamp, _monto, _motivo)
		RETURNING nrorechazo INTO numerorechazo;

		--mover INSERT rechazo
		PERFORM chequearRechazoLimites(numerorechazo);

		END;

	$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func spTestearConsumo() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION testear_consumo() returns trigger as $$
		BEGIN
		
		PERFORM autorizarcompra(new.nrotarjeta,new.codseguridad, new.nrocomercio, new.monto);
		return new;
		END;

	$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func spAgregarAlertaRechazo() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION agregar_alerta() returns trigger as $$
		BEGIN

		INSERT INTO alerta(nrotarjeta,fecha,nrorechazo,codalerta,descripcion) VALUES(new.nrotarjeta, new.fecha, new.nrorechazo, 0 , new.motivo);

		return new;
		END;

	$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func spSeguridadCompras() {
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
