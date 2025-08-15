package domain

import "time"

type Role string

const (
	RoleBuyer  Role = "buyer"
	RoleSeller Role = "seller"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Product struct {
	ID             int64     `json:"id"`
	SellerID       int64     `json:"seller_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	PriceCents     int64     `json:"price_cents"` // BIGINT
	Stock          int       `json:"stock"`
	CoverPictureID *int64    `json:"cover_picture_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Picture struct {
	ID        int64     `json:"id"`
	MIMEType  string    `json:"mime_type"`
	SizeBytes int64     `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
	Position  int       `json:"position,omitempty"` // позиция в рамках продукта
}
