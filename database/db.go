package database

import (
	"database/sql"
	"log"
	"os" // Aggiunto per leggere le variabili d'ambiente

	_ "github.com/lib/pq" // Driver per PostgreSQL
)

// DB è una variabile globale esportata che gli altri pacchetti useranno
var DB *sql.DB

// Connect inizializza la connessione al database
func Connect() {
	// 1. Prova a leggere l'URL del cloud (Neon) dalle variabili di sistema
	connStr := os.Getenv("DATABASE_URL")

	// 2. Se è vuoto (ovvero sei sul tuo PC in locale), usa il tuo database locale
	if connStr == "" {
		connStr = "user=postgres password=aicha dbname=halalshop sslmode=disable"
	}

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
