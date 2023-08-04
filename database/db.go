package database

import (
	"errors"
	"fmt"
	"github.com/Godzizizilla/douyin-simple/module"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	DB *gorm.DB
)

func InitDB() {
	// Replace 'your_username', 'your_password', 'your_database' with your MySQL credentials
	dsn := "root:mysql2023@tcp(127.0.0.1:3306)/douyin_simple_db?charset=utf8mb4&parseTime=True&loc=Local"

	// Connect to the MySQL database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	fmt.Println("Connected to the database!")

	db.AutoMigrate(&module.User{}, &module.Video{})
}

func GetUserInfoByID(userID uint) (*module.ApiUser, error) {
	var user module.ApiUser
	if err := DB.Model(&module.User{}).First(&user, userID).Error; err != nil {
		return nil, errors.New("get user info by id failed")
	}
	return &user, nil
}

func GetUserByUsername(username string) (userID uint, encryptedPassword string, err error) {
	var result struct {
		ID                uint
		EncryptedPassword string
	}
	err = DB.Model(&module.User{}).Select("id", "encrypted_password").First(&result, "name = ?", username).Error
	if err != nil {
		return 0, "", errors.New("user not found")
	}
	return result.ID, result.EncryptedPassword, nil
}

func AddUser(userName string, password string) (userID uint, err error) {
	var user = module.User{
		Name:              userName,
		EncryptedPassword: password,
	}
	if err := DB.Create(&user).Error; err != nil {
		return 0, errors.New("add user failed")
	}
	return user.ID, nil
}

func AddVideo(video *module.Video) error {
	if err := DB.Create(video).Error; err != nil {
		return errors.New("add video failed")
	}
	return nil
}

func GetUserVideosByID(userID uint) (*[]module.ApiVideo, error) {
	var videos []module.ApiVideo
	if err := DB.Model(&module.Video{}).Find(&videos, "user_id = ?", userID).Error; err != nil {
		return nil, errors.New("该用户未发布视频")
	}
	return &videos, nil
}

func GetVideosBeforeTimestamp(timestamp time.Time) (*[]module.ApiVideo, error) {
	/*
		var video module.Video
		db.Model(&module.Video{}).Joins("User").Where("user_id = ?", 4).Find(&video)
		fmt.Println(video.User.Name)
	*/

	var videos []module.ApiVideo
	if err := DB.Model(&module.Video{}).Where("created_at < ?", timestamp).Find(&videos).Error; err != nil {
		return nil, errors.New("没有在这之前的视频")
	}
	return &videos, nil
}
