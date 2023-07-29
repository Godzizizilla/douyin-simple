package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	db.AutoMigrate(&User{}, &Video{})
}

func IsUserExists(name string) (bool, string) {
	var user User
	err := DB.Where("name = ?", name).First(&user).Error
	if err == nil {
		return true, user.PasswordHash // 用户名已存在
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, "" // 用户名不存在
	}
	// TODO 处理数据库查询错误
	return false, ""
}

func AddUser(user *User) (int64, error) {
	if err := DB.Create(user).Error; err != nil {
		return 0, errors.New("add user failed")
	}
	return user.ID, nil
}

func GetUserIdByName(name string) (int64, error) {
	var user User
	err := DB.Where("name = ?", name).First(&user).Error
	if err == nil {
		return user.ID, nil
	} else {
		return 0, errors.New("get user id failed")
	}
}

func GetUserInfoByID(id int64) (*User, error) {
	var user User
	err := DB.First(&user, id).Error
	if err == nil {
		return &user, nil
	} else {
		return nil, errors.New("get user info by id failed")
	}
}

func AddVideo(video *Video) error {
	if err := DB.Create(video).Error; err != nil {
		return errors.New("add video failed")
	}
	return nil
}

func GetVideosByUserID(userID int64, videos *[]Video) {
	DB.Where("user_id = ?", userID).Find(&videos)
}
