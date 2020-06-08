package sql

func spTestConsumoRechazo() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION testConsumoRechazo(rechazo text default null) returns char(16) as $$  
			-- Rechazos: 
			-- 	'' consumo sin rechazo, se realiza la compra
			-- 'Tarjeta no válida'
			-- 'Código de seguridad inválido'
			-- 'Supera límite de tarjeta' 
			-- 'Plazo de vigencia expirado'
			-- 'La tarjeta se encuentra suspendida'
			DECLARE
				_nroTarjeta char(16);
				_codigoSeguridad char(4);
				_comercio int;
				_monto decimal(7,2);
				_disponible decimal(8,2);
	
			BEGIN 
		
			IF(rechazo is not null) THEN
				IF(rechazo not in('Tarjeta no válida', 'Código de seguridad inválido', 'Supera límite de tarjeta' , 'Plazo de vigencia expirado','La tarjeta se encuentra suspendida','')) THEN
					RETURN '';	
				END IF;
			END IF;
		   
		-- Comercio aleatorio
			SELECT nrocomercio INTO _comercio 
			FROM comercio 
			order by random() limit 1;
	
		-- Genero monto al azar
			_monto = 100 + random()*900;
			perform trunc(_monto,2);
	
			-- Genera caso con rechazo tarjeta invalida
			
			IF (rechazo = 'La tarjeta se encuentra suspendida') THEN
					raise notice 'Tarjeta suspendida';
					SELECT nrotarjeta, codseguridad, obtenerDisponible(nrotarjeta) as disponible
					into _nroTarjeta, _codigoSeguridad 
					FROM tarjeta 
					WHERE estado = 'suspendida' 
					and obtenerDisponible(nrotarjeta) > 1000
					ORDER BY random()
					LIMIT 1;
			
			
			ELSE
				IF(rechazo = 'Tarjeta no válida') THEN
					SELECT nrotarjeta, codseguridad, obtenerDisponible(nrotarjeta) as disponible
					into _nroTarjeta, _codigoSeguridad 
					FROM tarjeta 
					WHERE estado != 'vigente' AND estado != 'suspendida' 
					and obtenerDisponible(nrotarjeta) > 1000
					ORDER BY random()
					LIMIT 1;
				ELSE
					SELECT nrotarjeta, codseguridad, obtenerDisponible(nrotarjeta)
					into _nroTarjeta, _codigoSeguridad, _disponible
					FROM tarjeta 
					WHERE estado = 'vigente' 
			  	  	and obtenerDisponible(nrotarjeta) > 1000
					ORDER BY random()
					LIMIT 1;
				END IF;
				
				IF(rechazo = 'Código de seguridad inválido') THEN
					_codigoSeguridad = (_codigoSeguridad::int  + 1)::char(4);	
				END IF;
	
				if(rechazo = 'Supera límite de tarjeta' ) THEN
					_monto = _disponible  + 100;
				END IF;
	
				if(rechazo = 'Plazo de vigencia expirado') THEN
					SELECT tf.nrotarjeta, tf.codseguridad
					INTO _nroTarjeta, _codigoSeguridad
					FROM (SELECT ((TO_DATE(validahasta ||'01','YYYYMMDD')) +  interval '1 month') as fechavence, nrotarjeta, codseguridad FROM tarjeta) as tf
					WHERE tf.fechavence < current_date
					ORDER BY random();
				END IF;
	
			END IF;

	
			IF(_nroTarjeta is null) THEN
				RETURN '';
			ELSE
			-- 	Codigo Original
				INSERT INTO consumo VALUES(_nroTarjeta,_codigoSeguridad,_comercio,cast(_monto as decimal(7,2)));	
				--RETURN true;
				RETURN _nroTarjeta;	
			END IF;
	--	RETURN _nroTarjeta;
		END;
	$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func spTestConsumoAlerta() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION testConsumoAlerta(alerta int) returns char(16) as $$  
		-- codalerta: 
		-- 1: compra 1min
		-- 5: compra 5min
		-- 32: límite
		DECLARE
			_nroTarjeta char(16);
			_codigoSeguridad char(4);
			_comercio int;
			_monto decimal(7,2);
			_codPostal char(8);
	
		BEGIN 
			IF (alerta not in (1,5,32)) THEN
				Return '';
			END IF;
		
		-- Selecciono una tarjeta al azar	
			SELECT nrotarjeta, codseguridad
			INTO _nroTarjeta, _codigoSeguridad
			FROM tarjeta 
			WHERE estado = 'vigente' 
			  and obtenerDisponible(nrotarjeta) > 1000
			ORDER BY random()
			LIMIT 1;	
	
		-- Comercio aleatorio inicial
			SELECT nrocomercio, codigopostal INTO _comercio, _codPostal 
			FROM comercio 
			order by random() limit 1;
	
			FOR tarjeta IN 1 .. 2 BY 1
			LOOP
				-- Genero monto al azar
				_monto = 100 + random()* 900;
				perform trunc(_monto,2);
				
				IF(alerta = 1) THEN
					SELECT nrocomercio, codigopostal 
					INTO _comercio, _codPostal 
					FROM comercio 
					WHERE codigopostal = _codPostal 
					and nrocomercio != _comercio
					order by random() 
					limit 1;
				END IF;
	
				IF(alerta = 5) THEN
					SELECT nrocomercio, codigopostal 
					INTO _comercio, _codPostal 
					FROM comercio 
					WHERE codigopostal != _codPostal 
					order by random() 
					limit 1;
				END IF;
	
				IF(alerta = 32) THEN
					_monto = 99999;
				END IF;
	
				INSERT INTO consumo VALUES(_nroTarjeta,_codigoSeguridad,_comercio,_monto);			
			END LOOP;		
	
			RETURN _nroTarjeta;	
			END;
		$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func correrTest() {
	spObtenerDisponible();
	spTestConsumoRechazo();
	spTestConsumoAlerta();
	 consumoTest();
	 alertaTest();
	 _, err = db.Query(
		`   SELECT consumoTest(''),
			consumoTest('Tarjeta no válida'),
			consumoTest('La tarjeta se encuentra suspendida'),
			consumoTest('Plazo de vigencia expirado'),
			consumoTest('Código de seguridad inválido'),
			consumoTest('Supera límite de tarjeta'),
			alertaTest(1),
			alertaTest(5),
			alertaTest(32);
			`)
	logErr(err);
}

func consumoTest() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION consumoTest(_motivo text)returns boolean as $$
		DECLARE
		ntarjeta char(16);
		BEGIN
		
		SELECT into ntarjeta * FROM testConsumoRechazo(_motivo);
		
		IF (_motivo = '') THEN
			PERFORM * FROM compra WHERE nrotarjeta = ntarjeta;
			if(found) THEN
			return True;
			END IF;
		ELSE
			PERFORM * FROM rechazo, alerta WHERE alerta.nrotarjeta = ntarjeta 
			AND rechazo.nrotarjeta = ntarjeta
			AND rechazo.motivo = _motivo
			AND alerta.codalerta = 0 ;
			if(found) THEN
			return True;
			END IF;
		END IF;
		
		return False;
		END;
		$$ LANGUAGE PLPGSQL;
	
		`)
	logErr(err)
}

func alertaTest() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION alertaTest(_codalerta int)returns boolean as $$
		DECLARE
		ntarjeta char(16);
		BEGIN
		
		SELECT into ntarjeta * FROM testConsumoAlerta(_codalerta);
		
		IF (ntarjeta = '') THEN
			raise notice 'Inserte un código válido';
			return False;
		
		ELSE
			PERFORM * FROM alerta WHERE nrotarjeta = ntarjeta AND codalerta = _codalerta;
			if(found) THEN
			return True;
			END IF;
		END IF;
		
		return False;
		END;
		$$ LANGUAGE PLPGSQL;
	
		`)
	logErr(err)
}
