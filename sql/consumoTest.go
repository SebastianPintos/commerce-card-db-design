package sql

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
