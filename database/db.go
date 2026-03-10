package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv" // Il pacchetto per leggere il .env
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	// 1. Prova a caricare il file .env (questo funzionerà sul tuo PC locale)
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: Nessun file .env trovato, userò le variabili di sistema (tipico su Vercel)")
	}

	// 2. Legge la stringa magica dalla cassaforte (o da Vercel)
	connStr := os.Getenv("DATABASE_URL")

	// 3. Se la cassaforte è vuota, blocchiamo tutto prima di fare danni
	if connStr == "" {
		log.Fatal("ERRORE FATALE: Variabile DATABASE_URL mancante! Controlla il file .env o le impostazioni di Vercel.")
	}

	var dbErr error
	DB, dbErr = sql.Open("postgres", connStr)
	if dbErr != nil {
		log.Fatal("Errore di connessione a Postgres:", dbErr)
	}

	if dbErr = DB.Ping(); dbErr != nil {
		log.Fatal("Impossibile raggiungere il database:", dbErr)
	}

	log.Println("Connesso al database con successo in totale sicurezza!")
}
