package main

import (
	"fmt"
	"log"
	"net/http"
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

	fmt.Println("🚀 Accendo il motore e apro il sito in automatico...")

	// 1. ORA APRIAMO IL BROWSER SULLA PORTA 3000
	go func() {
		time.Sleep(1 * time.Second)
		openBrowser("http://localhost:3090")
	}()

	// 2. ACCENDIAMO IL SERVER SULLA PORTA 3090
	err = http.ListenAndServe(":3090", nil)
	if err != nil {
		log.Fatal("Errore del server:", err)
	}
}
