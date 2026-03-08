package api

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"halalshop/database"
	"halalshop/handlers"
)

// Magia: Inglobiamo SIA i templates HTML che la cartella static (CSS)!
//
//go:embed templates/* static/*
var embeddedFiles embed.FS

var dbInitialized bool

func init() {
	if !dbInitialized {
		database.Connect()
		dbInitialized = true
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	percorso := r.URL.Path

	// 1. ROTTA: File Statici (Il nostro amato CSS!)
	if strings.HasPrefix(percorso, "/static/") {
		http.FileServer(http.FS(embeddedFiles)).ServeHTTP(w, r)
		return
	}

	// 2. ROTTA: Home Page
	if percorso == "/" {
		prodotti := handlers.GetAllProducts()
		tmpl, err := template.ParseFS(embeddedFiles, "templates/index.html")
		if err != nil {
			http.Error(w, "Errore caricamento: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, prodotti)
		return
	}

	// 3. ROTTA: Pagina di Upload
	if percorso == "/upload" {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFS(embeddedFiles, "templates/upload.html")
			if err != nil {
				http.Error(w, "Errore caricamento: "+err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			err := r.ParseMultipartForm(10 << 20)
			if err != nil {
				http.Error(w, "Errore modulo", http.StatusBadRequest)
				return
			}

			nome := r.FormValue("name")
			descrizione := r.FormValue("description")
			imageURL := "" // Ricorda: Niente upload fisico di immagini su Vercel

			err = handlers.AddProduct(nome, descrizione, imageURL)
			if err != nil {
				http.Error(w, "Errore database", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	http.NotFound(w, r)
}
