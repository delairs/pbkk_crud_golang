package models

import (
	"time"
)

type History struct {
	History_ID uint      `gorm:"primaryKey"`                // Primary Key untuk tabel History
	UserID     uint      `gorm:"index"`                     // Foreign Key untuk User
	GameID     uint      `gorm:"index"`                     // Foreign Key untuk Game
	Timestamp  time.Time `gorm:"default:current_timestamp"` // Waktu tindakan dilakukan

	User      User `gorm:"foreignKey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE"` // Relasi ke tabel User
	Game      Game `gorm:"foreignKey:GameID;constraint:onUpdate:CASCADE,onDelete:CASCADE"` // Relasi ke tabel Game
	UpdatedAt time.Time
}
