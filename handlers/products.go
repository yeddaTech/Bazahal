package handlers

import (
	"log"

	"halalshop/database"
	"halalshop/models"
)

// GetAllProducts legge tutti i prodotti dal database
func GetAllProducts() []models.Product {
	// Usiamo COALESCE per dire a Postgres: "se image_url è vuoto (NULL), passami una stringa vuota"
	// Aggiungiamo anche "ORDER BY id ASC" per avere i prodotti in ordine!
	rows, err := database.DB.Query("SELECT id, name, description, COALESCE(image_url, '') FROM products ORDER BY id ASC")
	if err != nil {
		log.Println("Errore durante la query:", err)
		return nil
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		// Aggiungiamo &p.ImageURL alla fine dello Scan
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImageURL)
		if err != nil {
			log.Println("Errore nello scan:", err)
			continue
		}
		products = append(products, p)
	}

	return products
}

// AddProduct ora accetta anche l'URL dell'immagine!
func AddProduct(name string, description string, imageURL string) error {
	query := "INSERT INTO products (name, description, image_url) VALUES ($1, $2, $3)"

	_, err := database.DB.Exec(query, name, description, imageURL)
	if err != nil {
		log.Println("Errore durante l'inserimento:", err)
		return err
	}

	log.Println("Nuovo prodotto aggiunto:", name)
	return nil
}
