package main

import (
	"fmt"
	"log"
	"net/http"
	"os" // Aggiunto per leggere le variabili d'ambiente
	"os/exec"
	"runtime"
	"time"

	"halalshop/api"

	"github.com/joho/godotenv"
)

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		fmt.Println("⚠️ Non sono riuscito ad aprire il browser. Vai manualmente su:", url)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("ℹ️ Nessun file .env trovato (normale se sei su Vercel)")
	}

	http.HandleFunc("/", api.Handler)

	// 1. CHIEDIAMO LA PORTA AL SISTEMA
	port := os.Getenv("PORT")
	isLocal := false

	if port == "" {
		port = "3090" // Fallback per lo sviluppo sul tuo PC
		isLocal = true
	}

	fmt.Printf("🚀 Accendo il motore sulla porta %s...\n", port)

	// 2. APRIAMO IL BROWSER SOLO IN LOCALE
	if isLocal {
		go func() {
			time.Sleep(1 * time.Second)
			openBrowser("http://localhost:" + port)
		}()
	} else {
		fmt.Println("☁️ Ambiente Cloud rilevato: avvio serverless in corso...")
	}

	// 3. ACCENDIAMO IL SERVER SULLA PORTA DINAMICA
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Errore del server:", err)
	}
}
