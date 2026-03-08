package main

import (
	"fmt"
	"log"
	"net/http"

	"halalshop/api" // Importiamo la cartella api dove c'è la nostra logica
)

func main() {
	// Diciamo a Go: "Qualsiasi cosa l'utente cerchi, passala al nostro Handler di Vercel"
	http.HandleFunc("/", api.Handler)

	fmt.Println("🚀 Server locale avviato! Apri il browser e vai su: http://localhost:8080")

	// Accendiamo il server sulla porta 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Errore del server:", err)
	}
}
