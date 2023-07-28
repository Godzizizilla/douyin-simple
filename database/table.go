package database

import (
	"gorm.io/gorm"
	"time"
)

type CustomModel struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	CustomModel
	Name            string `gorm:"unique"`
	FollowCount     uint
	FollowerCount   uint
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  uint
	WorkCount       uint
	FavoriteCount   uint

	PasswordHash string
}
