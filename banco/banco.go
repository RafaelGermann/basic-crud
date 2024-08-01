package banco

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func Conectar() (*sql.DB, error) {
	stringConnection := "host=localhost port=5432 user=postgres password=admin123 dbname=Banco sslmode=disable"
	db, erro := sql.Open("postgres", stringConnection)
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil
}
