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

	// 4. UPLOAD PAGE (BLINDATA CON COOKIE)
	if percorso == "/upload" {
		// 🛑 IL BUTTAFUORI: Controlla se l'utente ha il Cookie VIP
		cookie, err := r.Cookie("admin_session")
		if err != nil || cookie.Value != "loggato_con_successo" {
			// Niente cookie o cookie sbagliato? Fuori! Rimbalzato alla home.
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Se ha il cookie, lo facciamo passare alle funzionalità di upload
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
				http.Error(w, "Errore nella registrazione. Forse lo username esiste già?", http.StatusBadRequest)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	// 6. LOGIN (GENERA IL COOKIE)
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

			if handlers.LoginUser(username, password) {
				// 🔥 CREIAMO IL COOKIE "VIP"
				cookie := &http.Cookie{
					Name:     "admin_session",
					Value:    "loggato_con_successo",
					Path:     "/",
					HttpOnly: true,  // Lo nasconde agli script JS (sicurezza extra)
					MaxAge:   86400, // Dura 24 ore (in secondi)
				}
				http.SetCookie(w, cookie) // Incolla il cookie sul browser

				// Reindirizza al quartier generale
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

	// 8. PAGINA MACELLERIE E RISTORANTI HALAL
	if percorso == "/macelleriehalal" {
		tmpl, err := template.ParseFS(embeddedFiles, "templates/macelleriehalal.html")
		if err != nil {
			http.Error(w, "Errore caricamento: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	// 9. LOGOUT (DISTRUGGE IL COOKIE)
	if percorso == "/logout" {
		cookie := &http.Cookie{
			Name:     "admin_session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1, // Il valore negativo distrugge istantaneamente il cookie
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 10. PAGINA 404 PERSONALIZZATA
	w.WriteHeader(http.StatusNotFound) // Dice al browser che la pagina è un errore 404 reale
	tmpl, err := template.ParseFS(embeddedFiles, "templates/404.html")
	if err == nil {
		tmpl.Execute(w, nil)
	} else {
		// Fallback di sicurezza se per caso elimini il file 404.html
		http.NotFound(w, r)
	}
}
