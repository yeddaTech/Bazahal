package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Driver per PostgreSQL
)

// DB è una variabile globale esportata che gli altri pacchetti userentanno
var DB *sql.DB

// Connect inizializza la connessione al database
func Connect() {
	connStr := "user=postgres password=aicha dbname=halalshop sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Errore di connessione a Postgres:", err)
	}

	// Verifica che la connessione sia effettivamente attiva
	if err = DB.Ping(); err != nil {
		log.Fatal("Impossibile raggiungere il database:", err)
	}

	log.Println("Connesso al database con successo!")
}
