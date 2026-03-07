package models

type Product struct {
	ID          int
	Name        string
	Description string
	ImageURL    string // <-- Abbiamo aggiunto questa riga!
}
