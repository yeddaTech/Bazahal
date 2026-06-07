package database

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: Nessun file .env trovato")
	}

	// --- SONDA DI DEBUG (lasciamola per vedere se Vercel si sveglia) ---
	log.Println("--- INIZIO CHECK VARIABILI VERCEL ---")
	for _, env := range os.Environ() {
		chiave := strings.Split(env, "=")[0]
		if strings.Contains(chiave, "DATABASE") {
			log.Printf("Trovata chiave simile a DATABASE: [%s]", chiave)
		}
	}
	log.Println("--- FINE CHECK ---")
	// ----------------------

	// 1. COMMENTIAMO LA CHIAMATA A VERCEL
	// connStr := os.Getenv("DATABASE_URL")
	// if connStr == "" {
	//     log.Fatal("ERRORE FATALE: Variabile DATABASE_URL mancante!")
	// }

	connStr := os.Getenv("NEON_DB_URL") // <-- CAMBIA DA DATABASE_URL A NEON_DB_URL

	if connStr == "" {
		log.Fatal("ERRORE FATALE: Variabile NEON_DB_URL mancante!")
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
