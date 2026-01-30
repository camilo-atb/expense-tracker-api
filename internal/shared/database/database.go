package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func DatabaseConnection() (*sql.DB, error) {
	driver := "pgx"
	host := "localhost"
	port := 5432
	usuario := "postgres"
	contrasenia := os.Getenv("postgresql_password")
	nombreDB := os.Getenv("postgresql_name_db")

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		usuario,
		contrasenia,
		nombreDB,
	)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

/*
CREATE TYPE status AS ENUM ('active', 'removed');
CREATE TYPE category AS ENUM ('comida', 'transporte', 'ocio', 'entretenimiento', 'estudios');
CREATE TYPE transaction_type AS ENUM ('income', 'expense');


CREATE TABLE IF NOT EXISTS tasks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    description TE) NOT NULL,
    amount TEXT NOT NULL,
    category estado NOT NULL DEFAULT 'pending',
    type TIMESTAMP NOT NULL DEFAULT NOW(),
    maximum_term DATE NOT NULL
);
*/
