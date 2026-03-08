package api

import (
	"embed"
	"html/template"
	"net/http"

	"halalshop/database"
	"halalshop/handlers"
)

// Questa è la magia: dice a Go di "inglobare" tutti i file HTML nel codice compilato!
//
//go:embed templates/*
var templateFiles embed.FS

var dbInitialized bool

func init() {
	if !dbInitialized {
		database.Connect()
		dbInitialized = true
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	percorso := r.URL.Path

	// 1. ROTTA: Home Page
	if percorso == "/" {
		prodotti := handlers.GetAllProducts()

		// Invece di ParseFiles (che cerca sul disco), usiamo ParseFS (che cerca nella memoria)
		tmpl, err := template.ParseFS(templateFiles, "templates/index.html")
		if err != nil {
			http.Error(w, "Errore nel caricamento della pagina: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, prodotti)
		return
	}

	// 2. ROTTA: Pagina di Upload
	if percorso == "/upload" {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFS(templateFiles, "templates/upload.html")
			if err != nil {
				http.Error(w, "Errore nel caricamento della pagina: "+err.Error(), http.StatusInternalServerError)
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

			// Nota: ho rimosso il salvataggio fisico dell'immagine perché su Vercel non funziona.
			// Per ora inseriamo un URL vuoto o un placeholder, poi lo adatteremo ai link di affiliazione!
			imageURL := ""

			err = handlers.AddProduct(nome, descrizione, imageURL)
			if err != nil {
				http.Error(w, "Errore salvataggio database", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	http.NotFound(w, r)
}
