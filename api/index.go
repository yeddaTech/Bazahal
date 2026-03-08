package api

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"halalshop/database"
	"halalshop/handlers"
)

// Usiamo init() per connetterci al database solo quando Vercel "sveglia" la funzione
var dbInitialized bool

func init() {
	if !dbInitialized {
		database.Connect()
		dbInitialized = true
	}
}

// Handler è la funzione magica che Vercel chiama per OGNI richiesta
func Handler(w http.ResponseWriter, r *http.Request) {
	percorso := r.URL.Path

	// 1. ROTTA: File Statici (CSS, immagini già presenti)
	if strings.HasPrefix(percorso, "/static/") {
		fs := http.StripPrefix("/static/", http.FileServer(http.Dir("../static")))
		fs.ServeHTTP(w, r)
		return
	}

	// 2. ROTTA: Home Page
	if percorso == "/" {
		prodotti := handlers.GetAllProducts()
		// Vercel lavora dalla cartella principale, quindi cerchiamo i template lì
		tmpl, err := template.ParseFiles("../templates/index.html")
		if err != nil {
			http.Error(w, "Errore nel caricamento della pagina: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, prodotti)
		return
	}

	// 3. ROTTA: Pagina di Upload
	if percorso == "/upload" {
		// Se l'utente vuole solo VEDERE il modulo (GET)
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("../templates/upload.html")
			if err != nil {
				http.Error(w, "Errore nel caricamento della pagina", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
			return
		}

		// Se l'utente ha premuto "SALVA PRODOTTO" (POST)
		if r.Method == http.MethodPost {
			err := r.ParseMultipartForm(10 << 20)
			if err != nil {
				http.Error(w, "Errore modulo", http.StatusBadRequest)
				return
			}

			nome := r.FormValue("name")
			descrizione := r.FormValue("description")
			imageURL := ""

			// ⚠️ ATTENZIONE: Su Vercel questo blocco per le immagini fisiche non funzionerà bene
			// perché Vercel non permette di salvare file nella cartella static/.
			// In futuro, dovrai passare un semplice link testuale dell'immagine (es. URL di Amazon)!
			file, header, err := r.FormFile("image")
			if err == nil {
				defer file.Close()
				os.MkdirAll("../static/uploads", os.ModePerm) // Tentativo di creare la cartella
				percorsoFisico := "../static/uploads/" + header.Filename
				out, err := os.Create(percorsoFisico)
				if err == nil {
					defer out.Close()
					io.Copy(out, file)
					imageURL = "/static/uploads/" + header.Filename
				}
			}

			// Salviamo nel database
			err = handlers.AddProduct(nome, descrizione, imageURL)
			if err != nil {
				http.Error(w, "Errore salvataggio", http.StatusInternalServerError)
				return
			}

			// Riportiamo l'utente alla home
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	// Se qualcuno cerca una pagina che non esiste
	http.NotFound(w, r)
}
