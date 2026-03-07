package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"halalshop/database"
	"halalshop/handlers"
)

func main() {
	database.Connect()
	defer database.DB.Close()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		prodotti := handlers.GetAllProducts()
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Errore nel caricamento della pagina", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, prodotti)
	})

	// Rotta per la pagina di Upload
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

		// 1. Se l'utente vuole solo VEDERE il modulo (GET)
		if r.Method == http.MethodGet {
			tmpl, err := template.ParseFiles("templates/upload.html")
			if err != nil {
				http.Error(w, "Errore nel caricamento della pagina", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
			return
		}

		// 2. Se l'utente ha premuto "SALVA PRODOTTO" (POST)
		if r.Method == http.MethodPost {
			err := r.ParseMultipartForm(10 << 20) // Limite 10 MB
			if err != nil {
				http.Error(w, "Errore modulo", http.StatusBadRequest)
				return
			}

			nome := r.FormValue("name")
			descrizione := r.FormValue("description")
			imageURL := "" // Se non carica nulla, l'immagine sarà vuota

			// RECUPERIAMO IL FILE DAL MODULO
			file, header, err := r.FormFile("image")
			if err == nil {
				defer file.Close()

				// Creiamo il file fisico nella cartella static/uploads/
				percorsoFisico := "static/uploads/" + header.Filename
				out, err := os.Create(percorsoFisico)
				if err == nil {
					defer out.Close()
					io.Copy(out, file) // Copiamo il contenuto

					// Questo è l'URL che salveremo nel Database!
					imageURL = "/static/uploads/" + header.Filename
				}
			}

			// Salviamo tutto nel database (compresa l'immagine!)
			err = handlers.AddProduct(nome, descrizione, imageURL)
			if err != nil {
				http.Error(w, "Errore salvataggio", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	})

	fmt.Println("🚀 Server avviato! Apri il browser e vai su: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
