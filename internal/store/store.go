package store

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type DbReceipt struct {
	Id      string
	Details string
	Score   int
}

func NewDbReceipt(id string, details string, score int) DbReceipt {
	return DbReceipt{
		Id:      id,
		Details: details,
		Score:   score,
	}
}

type StoreInterface interface {
	PutReceipt(DbReceipt) error
	GetReceipt(string) (DbReceipt, error)
}

type Store struct {
	StoreInterface
	db *sql.DB
}

func NewStore(dbFile string) Store {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		slog.Error("Error opening file",
			"error", err,
			"file", dbFile,
		)
	}

	if err := migrations(db); err != nil {
		slog.Error("Error running migrations",
			"error", err,
		)
	}

	return Store{
		db: db,
	}
}

func (s *Store) PutReceipt(dbReceipt DbReceipt) error {
	return putReceipt(dbReceipt, s.db)
}

func putReceipt(dbReceipt DbReceipt, db *sql.DB) error {
	if _, err := db.Exec(
		`INSERT INTO receipt VALUES(?,?,?);`,
		dbReceipt.Id,
		dbReceipt.Details,
		dbReceipt.Score,
	); err != nil {
		return err
	}
	return nil
}

func (s *Store) GetReceipt(id string) (DbReceipt, error) {
	return getReceipt(id, s.db)
}

func getReceipt(id string, db *sql.DB) (DbReceipt, error) {
	dbReceipt := DbReceipt{}

	row := db.QueryRow(
		`SELECT * FROM receipt WHERE id=?;`,
		id,
	)
	if err := row.Scan(
		&dbReceipt.Id,
		&dbReceipt.Details,
		&dbReceipt.Score,
	); err != nil {
		return dbReceipt, err
	}
	return dbReceipt, nil
}

func migrations(db *sql.DB) error {
	const schema string = `
        CREATE TABLE IF NOT EXISTS receipt (
            id      TEXT PRIMARY KEY,
            details TEXT NOT NULL,
            score   INTEGER NOT NULL
        );
    `
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	return nil
}
