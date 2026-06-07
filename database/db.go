package database

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv" // Il pacchetto per leggere il .env
	_ "github.com/lib/pq"
)

var DB *sql.DB

// bypass vercel cache
// forza aggiornamento vercel
func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: Nessun file .env trovato")
	}

	// --- SONDA DI DEBUG ---
	log.Println("--- INIZIO CHECK VARIABILI VERCEL ---")
	for _, env := range os.Environ() {
		// Separiamo chiave e valore (stampiamo SOLO la chiave per non leakare la password)
		chiave := strings.Split(env, "=")[0]

		// Filtriamo per non intaccare troppo i log
		if strings.Contains(chiave, "DATABASE") {
			log.Printf("Trovata chiave simile a DATABASE: [%s]", chiave)
		}
	}
	log.Println("--- FINE CHECK ---")
	// ----------------------

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("ERRORE FATALE: Variabile DATABASE_URL mancante!")
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
