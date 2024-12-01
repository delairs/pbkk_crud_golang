package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`        // Primary Key
	Name      string `gorm:"type:varchar(100)"` // Nama user
	Email     string `gorm:"uniqueIndex"`       // Email harus unik
	Password  string `gorm:"type:varchar(255)"` // Password dienkripsi
	CreatedAt int64  // Timestamp untuk pencatatan waktu pembuatan
	UpdatedAt int64  // Timestamp untuk pencatatan waktu update
	IsAdmin   bool
	Token     string `gorm:"unique;"`
	ExpiresAt time.Time
}
