package models

import "gorm.io/gorm"

// Tabel User
type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Name     string
	GoogleID string `gorm:"uniqueIndex"`
	Links    []Link `gorm:"foreignKey:UserID"` // Relasi ke Link
}

// Tabel Link
type Link struct {
	gorm.Model
	LongURL   	string `gorm:"not null"`
	ShortCode 	string `gorm:"uniqueIndex;not null"`
	ShortURL  	string `json:"short_url"`
	QRURL     	string `json:"qr_url"`
	ClickCount	int    `gorm:"default:0" json:"click_count"`
	UserID    	uint   // Pemilik link (Foreign Key)
}