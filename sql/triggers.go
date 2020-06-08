package sql

func trAgregarConsumo() {
	_, err = db.Query(
		`CREATE trigger agregarconsumo_trg
		BEFORE INSERT ON consumo

		FOR EACH ROW
			EXECUTE PROCEDURE testear_consumo();

		`)
	logErr(err)
}

func trAgregarAlerta() {
	_, err = db.Query(
		`CREATE trigger agregaralerta_trg
		BEFORE INSERT ON rechazo

		FOR EACH ROW
			EXECUTE PROCEDURE agregar_alerta();

		`)
	logErr(err)
}

func trSeguridadCompras() {
	_, err = db.Query(
		`CREATE trigger compras_lapsotiempo
		BEFORE INSERT ON compra

		FOR EACH ROW
			EXECUTE PROCEDURE compras_lapso_tiempo();
		`)
	logErr(err)
}
