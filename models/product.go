package models

type Product struct {
	ID            int
	Name          string
	Description   string
	ImageURL      string
	AffiliateLink string // 💰 Il nuovo campo per le affiliazioni!
}
