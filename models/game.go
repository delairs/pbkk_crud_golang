package models

type Game struct {
	Game_ID     uint   `gorm:"primaryKey"`
	Game_Title  string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:varchar(255)"`
	Link_Embed  string `gorm:"type:varchar(255)"`
	Path_Image  string `gorm:"type:varchar(255)"`
	UpdatedAt   int64
}
