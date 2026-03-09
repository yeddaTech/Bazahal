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

	// 1. File Statici (CSS)
	if strings.HasPrefix(percorso, "/static/") {
		http.FileServer(http.FS(embeddedFiles)).ServeHTTP(w, r)
		return
	}

	// 2. HOME PAGE
	if percorso == "/" {
		tmpl, _ := template.ParseFS(embeddedFiles, "templates/index.html")
		tmpl.Execute(w, nil)
		return
	}

	// 3. SHOP PAGE
	if percorso == "/shop" {
		prodotti := handlers.GetAllProducts()
		tmpl, _ := template.ParseFS(embeddedFiles, "templates/shop.html")
		tmpl.Execute(w, prodotti)
		return
	}

	// 4. UPLOAD PAGE
	if percorso == "/upload" {
		if r.Method == http.MethodGet {
			tmpl, _ := template.ParseFS(embeddedFiles, "templates/upload.html")
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			r.ParseForm()
			handlers.AddProduct(r.FormValue("name"), r.FormValue("description"), r.FormValue("image_url"), r.FormValue("affiliate_link"))
			http.Redirect(w, r, "/shop", http.StatusSeeOther)
			return
		}
	}

	// 5. REGISTRAZIONE
	if percorso == "/register" {
		if r.Method == http.MethodGet {
			tmpl, _ := template.ParseFS(embeddedFiles, "templates/register.html")
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			r.ParseForm()
			username := r.FormValue("username")
			password := r.FormValue("password")

			err := handlers.RegisterUser(username, password)
			if err != nil {
				// Se lo username esiste già, diamo errore
				http.Error(w, "Errore nella registrazione. Forse lo username esiste già?", http.StatusBadRequest)
				return
			}
			// Se va tutto bene, lo mandiamo al login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	// 6. LOGIN (ORA È VERO!)
	if percorso == "/login" {
		if r.Method == http.MethodGet {
			tmpl, _ := template.ParseFS(embeddedFiles, "templates/login.html")
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			r.ParseForm()
			username := r.FormValue("username")
			password := r.FormValue("password")

			// Usiamo il database per controllare!
			if handlers.LoginUser(username, password) {
				http.Redirect(w, r, "/upload", http.StatusSeeOther)
				return
			} else {
				http.Error(w, "Credenziali errate. Riprova.", http.StatusUnauthorized)
				return
			}
		}
	}

	// 7. PAGINA ORARI PREGHIERA
	if percorso == "/orari" {
		tmpl, err := template.ParseFS(embeddedFiles, "templates/orari.html")
		if err != nil {
			http.Error(w, "Errore caricamento: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	http.NotFound(w, r)
}
