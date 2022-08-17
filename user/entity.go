package user

import (
	"time"
)

type User struct {
	ID         int       `json:"id" gorm:"primary_key"`
	Name       string    `json:"name"`
	Occupation string    `json:"occupation"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Avatar     string    `json:"avatar"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
