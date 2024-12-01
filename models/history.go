package models

import (
	"time"
)

type History struct {
	History_ID uint      `gorm:"primaryKey"`                // Primary Key untuk tabel History
	User_ID    uint      `gorm:"index"`                     // Foreign Key untuk User
	Game_ID    uint      `gorm:"index"`                     // Foreign Key untuk Game
	Timestamp  time.Time `gorm:"default:current_timestamp"` // Waktu tindakan dilakukan

	User      User `gorm:"foreignKey:User_ID;constraint:onUpdate:CASCADE,onDelete:CASCADE"` // Relasi ke tabel User
	Game      Game `gorm:"foreignKey:Game_ID;constraint:onUpdate:CASCADE,onDelete:CASCADE"` // Relasi ke tabel Game
	UpdatedAt time.Time
}
