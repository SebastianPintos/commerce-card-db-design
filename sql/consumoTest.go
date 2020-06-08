package sql

func spTestConsumoRechazo() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION testConsumoRechazo(rechazo text default null) returns bool as $$  
			-- Rechazos: 
			-- 	null consumo sin rechazo
			-- 'TarjetaInvalida'
			-- 'CodigoInvalido'
			-- 'LimiteSuperado' 
			-- 'TarjetaExpirada'
			DECLARE
				_nroTarjeta char(16);
				_codigoSeguridad char(4);
				_comercio int;
				_monto decimal(8,2);
				_disponible decimal(8,2);
	
			BEGIN 
		
			IF(rechazo is not null) THEN
				IF(rechazo not in('TarjetaInvalida', 'CodigoInvalido', 'LimiteSuperado', 'TarjetaExpirada')) THEN
					RETURN false;
				END IF;
			END IF;
		   
		-- Comercio aleatorio
			SELECT nrocomercio INTO _comercio 
			FROM comercio 
			order by random() limit 1;
	
		-- Genero monto al azar
			_monto = random()*(1000-100)+100;
	
			-- Genera caso con rechazo tarjeta invalida
			if(rechazo = 'TarjetaInvalida') THEN
				SELECT nrotarjeta, codseguridad, obtenerDisponible(nrotarjeta) as disponible
				into _nroTarjeta, _codigoSeguridad 
				FROM tarjeta 
				WHERE estado != 'vigente' 
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
	
				if(rechazo = 'CodigoInvalido') THEN
					_codigoSeguridad = (_codigoSeguridad::int  + 1)::char(4);	
				END IF;
	
				if(rechazo = 'LimiteSuperado') THEN
					_monto = _disponible  + 100;
				END IF;
	
				if(rechazo = 'TarjetaExpirada') THEN
					SELECT tf.nrotarjeta, tf.codseguridad
					INTO _nroTarjeta, _codigoSeguridad
					FROM (SELECT ((TO_DATE(validahasta ||'01','YYYYMMDD')) +  interval '1 month') as fechavence, nrotarjeta, codseguridad FROM tarjeta) as tf
					WHERE tf.fechavence < current_date
					ORDER BY random();
				END IF;
	
			END IF;
	
			IF(_nroTarjeta is null) THEN
				RETURN false;
			ELSE
			-- 	Codigo Original
				INSERT INTO consumo VALUES(_nroTarjeta,_codigoSeguridad,_comercio,_monto);	
				RETURN true;	
			END IF;
	
		END;
	$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func spTestConsumoAlerta() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION testConsumoAlerta(alerta int) returns bool as $$  
		-- codalerta: 
		-- 1: compra 1min
		-- 5: compra 5min
		-- 32: límite
		DECLARE
			_nroTarjeta char(16);
			_codigoSeguridad char(4);
			_comercio int;
			_monto decimal(8,2);
			_codPostal char(8);
	
		BEGIN 
			IF (alerta not in (1,5,32)) THEN
				Return false;
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
				_monto = random()*(1000-100)+100;
				
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
	
			RETURN true;	
			END;
		$$ LANGUAGE PLPGSQL;`)
	logErr(err)
}

func CorrerTest() {
	_, err = db.Exec(`select testConsumoRechazo();
					  select testConsumoRechazo('TarjetaInvalida');
					  select testConsumoRechazo('CodigoInvalido');
					  select testConsumoRechazo('LimiteSuperado');
					  select testConsumoRechazo('TarjetaExpirada');
					  select testConsumoAlerta(1);
					  select testConsumoAlerta(5);
					  select testConsumoAlerta(32);`)
	logErr(err)
}

func consumoValidoTest() {

	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION consumoValidoTest()returns boolean as $$
		DECLARE
		cantidad int;
		BEGIN 
		
		INSERT INTO consumo VALUES(4032002134557009,070,520,100);
		SELECT INTO cantidad COUNT(*) FROM compra WHERE nrotarjeta = '4032002134557009';
		
		if(cantidad = 1) THEN
			raise notice 'Test consumo válido: True';
			return True;
		END IF;
		
		raise notice 'Test consumo válido: False';
		return False;
		END;
		$$ LANGUAGE PLPGSQL;`)
	logErr(err)

}

func consumoTarjetaInvalidaTest() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION consumoTarjetaInvalidaTest()returns boolean as $$
		DECLARE
		cantidad int;
		BEGIN
		
		INSERT INTO consumo VALUES(4037055274760805,109,520,100);
		
		SELECT INTO cantidad COUNT(*) FROM rechazo, alerta WHERE rechazo.nrotarjeta = '4037055274760805' 	
		AND alerta.nrotarjeta = '4037055274760805' 
		AND motivo='Tarjeta no válida' AND codalerta = 0;

		if(cantidad = 1) THEN
			raise notice 'Test consumo tarjeta inválida y alerta rechazo: True';
			return True;
		END IF;

		raise notice 'Test consumo tarjeta inválida y alerta rechazo: False';
		return False;
		END;
		$$ LANGUAGE PLPGSQL;
	
		`)
	logErr(err)
}

func consumoCodSeguridadInvalidoTest() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION consumoCodSeguridadInvalidoTest()returns boolean as $$
		DECLARE
		cantidad int;
		
		BEGIN
		INSERT INTO consumo VALUES(4000001234567899,732,520,100);
		SELECT INTO cantidad COUNT(*) FROM rechazo, alerta WHERE alerta.nrotarjeta = '4000001234567899'
		AND rechazo.nrotarjeta = '4000001234567899' 
		AND motivo='Código de seguridad inválido' AND codalerta = 0;
		
		if(cantidad = 1) THEN
			raise notice 'Test consumo código de seguridad inválido y alerta rechazo: True';
			return True;
		END IF;


		raise notice 'Test consumo código de seguridad inválido y alerta rechazo: False';
		return False;
		END;
		$$ LANGUAGE PLPGSQL;
	
		`)
	logErr(err)
}

func consumoExcedeLimiteTest() {
	_, err = db.Query(
		`
		CREATE OR REPLACE FUNCTION consumoExcedeLimiteTest()returns boolean as $$
		DECLARE
		cantidad int;

		BEGIN
		INSERT INTO consumo VALUES(4000006877865030,115,520,99999);
		SELECT INTO cantidad COUNT(*) FROM rechazo, alerta WHERE alerta.nrotarjeta = '4000001234567899' 
		AND rechazo.nrotarjeta = '4000006877865030' 
		AND motivo='Supera límite de tarjeta' AND codalerta = 0;
		
		if(cantidad = 1) THEN
			raise notice 'Test consumo excede límite de tarjeta y alerta rechazo: True';
			return True;
		END IF;


		raise notice 'Test consumo excede límite de tarjeta y alerta rechazo: False';
		return False;
		END;
		$$ LANGUAGE PLPGSQL;
	
		
		`)
	logErr(err)
}

func consumoTarjetaExpiradaTest() {
	_, err = db.Query(
		`
		CREATE OR REPLACE FUNCTION consumoTarjetaExpiradaTest()returns boolean as $$
		DECLARE
		cantidad int;
		
		BEGIN
		
		INSERT INTO consumo VALUES(4037001554363655,332,520,100);
		SELECT INTO cantidad COUNT(*) FROM rechazo, alerta WHERE alerta.nrotarjeta = '4037001554363655' 
		AND rechazo.nrotarjeta = '4037001554363655' 
		AND motivo='Plazo de vigencia expirado' AND codalerta = 0;
		
		if(cantidad = 1) THEN
			return True;
			raise notice 'Test consumo tarjeta expirada y alerta rechazo: True';
		END IF;
		

		raise notice 'Test consumo tarjeta expirada y alerta rechazo: False';
		return False;
		END;
		$$ LANGUAGE PLPGSQL;
	`)
	logErr(err)
}

func consumoTarjetaSuspendidaTest() {
	_, err = db.Query(
		`
		CREATE OR REPLACE FUNCTION consumoTarjetaSuspendidaTest()returns boolean as $$
		DECLARE
		cantidad int;
		
		BEGIN
		INSERT INTO consumo VALUES(4032002224865843,327,520,100);
		SELECT INTO cantidad COUNT(*) FROM rechazo, alerta WHERE rechazo.nrotarjeta = '4032002224865843' 
		AND alerta.nrotarjeta = '4032002224865843' 
		AND motivo='La tarjeta se encuentra suspendida' AND codalerta = 0;
		
		if(cantidad = 1) THEN
			return True;
		raise notice 'Test consumo tarjeta suspendida y alerta rechazo: True';	
		END IF;

		raise notice 'Test consumo tarjeta suspendida y alerta rechazo: False';
		return False;
		END;
		$$ LANGUAGE PLPGSQL;
		`)
	logErr(err)
}

func consumoAlerta1Test() {
	_, err = db.Query(
		`
		CREATE OR REPLACE FUNCTION consumoAlerta1Test()returns boolean as $$
		DECLARE
		cantidad int;

		BEGIN
		INSERT INTO consumo VALUES(4000001234567899,733,501,100);
		INSERT INTO consumo VALUES(4000001234567899,733,513,100);

		SELECT INTO cantidad COUNT(*) FROM alerta WHERE nrotarjeta = '4000001234567899' 
		AND codalerta = 1;

		if(cantidad = 1) THEN
			raise notice 'Test consumo alerta 1: True';
			return True;
		END IF;


		raise notice 'Test consumo alerta 1: False';
		return False;
		END;
		$$ LANGUAGE PLPGSQL;		
		`)
	logErr(err)
}

func consumoAlerta5Test() {
	_, err = db.Query(
		`
		CREATE OR REPLACE FUNCTION consumoAlerta5Test()returns boolean as $$
		DECLARE
		cantidad int;
		
		BEGIN
		INSERT INTO consumo VALUES(4000001234567899,733,501,100);
		INSERT INTO consumo VALUES(4000001234567899,733,564,100);
		SELECT INTO cantidad COUNT(*) FROM alerta WHERE nrotarjeta = '4000001234567899' 
		AND codalerta = 5;
		
		if(cantidad = 1) THEN
			raise notice 'Test consumo alerta 5: True';
			return True;
		END IF;

		raise notice 'Test consumo alerta 5: False';
		return False;
		END;
		$$ LANGUAGE PLPGSQL;
		`)
	logErr(err)
}

func consumoAlerta32Test() {
	_, err = db.Query(
		`
		CREATE OR REPLACE FUNCTION consumoAlerta32Test()returns boolean as $$
		DECLARE
		cantidad int;
		
		BEGIN
		INSERT INTO consumo VALUES(4037001554363655,332,520,99999);
		INSERT INTO consumo VALUES(4037001554363655,332,520,99999);
		SELECT INTO cantidad COUNT(*) FROM alerta, tarjeta WHERE alerta.nrotarjeta = '4037001554363655'
		AND tarjeta.nrotarjeta = '4037001554363655' AND tarjeta.estado = 'suspendida' 
		AND alerta.codalerta = 32;
		
		if(cantidad = 1) THEN
			raise notice 'Test consumo alerta 32: True';
			return True;
		END IF;
		
		raise  notice 'Test consumo alerta 32: False';
		return False;
		END;
		$$ LANGUAGE PLPGSQL;
		`)
	logErr(err)
}
