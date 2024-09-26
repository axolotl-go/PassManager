package models

import "gorm.io/gorm"

type PasswordEntry struct {
	gorm.Model

	UserID   uint   `gorm:"not null" json:"user_id"`
	Site     string `gorm:"not null" json:"site"`
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
}
