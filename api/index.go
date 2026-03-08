package api

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"halalshop/database"
	"halalshop/handlers"
)

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

	// File Statici (CSS)
	if strings.HasPrefix(percorso, "/static/") {
		http.FileServer(http.FS(embeddedFiles)).ServeHTTP(w, r)
		return
	}

	// Home Page
	if percorso == "/" {
		prodotti := handlers.GetAllProducts()
		tmpl, err := template.ParseFS(embeddedFiles, "templates/index.html")
		if err != nil {
			http.Error(w, "Errore: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, prodotti)
		return
	}

	// Pagina Aggiungi Prodotto
	if percorso == "/upload" {
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFS(embeddedFiles, "templates/upload.html")
			if err != nil {
				http.Error(w, "Errore: "+err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			// Non usiamo più ParseMultipartForm per i file fisici, solo form testuali
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Errore modulo", http.StatusBadRequest)
				return
			}

			nome := r.FormValue("name")
			descrizione := r.FormValue("description")
			imageURL := r.FormValue("image_url")           // Ora è un link!
			affiliateLink := r.FormValue("affiliate_link") // Il link per guadagnare!

			// Passiamo tutto al database
			err = handlers.AddProduct(nome, descrizione, imageURL, affiliateLink)
			if err != nil {
				http.Error(w, "Errore salvataggio", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	http.NotFound(w, r)
}
