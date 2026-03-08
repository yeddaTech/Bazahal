package handlers

import (
	"halalshop/database"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUser cifra la password e salva l'utente nel database
func RegisterUser(username, password string) error {
	// Generiamo un hash inviolabile della password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Salviamo l'utente
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2)`
	_, err = database.DB.Exec(query, username, string(hash))
	return err
}

// LoginUser controlla se l'utente esiste e se la password combacia
func LoginUser(username, password string) bool {
	var hash string

	// Cerchiamo l'utente nel DB
	query := `SELECT password_hash FROM users WHERE username = $1`
	err := database.DB.QueryRow(query, username).Scan(&hash)
	if err != nil {
		return false // Utente non trovato
	}

	// Confrontiamo la password inserita con l'hash salvato
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // Se err è nil, la password è giusta!
}
