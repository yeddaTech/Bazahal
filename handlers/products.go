package handlers

import (
	"halalshop/database"
	"halalshop/models"
	"log"
)

// Ora accetta 4 parametri!
func AddProduct(name, description, imageURL, affiliateLink string) error {
	// Aggiunto affiliate_link alla query SQL
	query := `INSERT INTO products (name, description, image_url, affiliate_link) VALUES ($1, $2, $3, $4)`

	_, err := database.DB.Exec(query, name, description, imageURL, affiliateLink)
	if err != nil {
		log.Println("Errore durante l'inserimento del prodotto:", err)
		return err
	}
	return nil
}

func GetAllProducts() []models.Product {
	var products []models.Product

	// Usiamo COALESCE per evitare crash se nel database l'immagine o il link sono vuoti (NULL)
	query := `SELECT id, name, description, COALESCE(image_url, ''), COALESCE(affiliate_link, '') FROM products ORDER BY id DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Println("Errore durante la lettura dei prodotti:", err)
		return products
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		// Leggiamo 5 campi, incluso l'AffiliateLink
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImageURL, &p.AffiliateLink)
		if err != nil {
			log.Println("Errore nello scan della riga:", err)
			continue
		}
		products = append(products, p)
	}
	return products
}
