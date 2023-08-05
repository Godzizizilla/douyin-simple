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

	db.AutoMigrate(&module.User{}, &module.Video{}, &module.Comment{})

	/*var video module.Video
	db.Model(&module.Video{}).
		Joins("LEFT JOIN like_videos ON videos.id = like_videos.video_id").
		Where("like_videos.user_id = ?", 4).
		Select("videos.*, EXISTS(SELECT 1 FROM like_videos WHERE like_videos.user_id = ? AND like_videos.video_id = videos.id) AS is_favorite", 4).
		First(&video, 2)
	fmt.Println(video.IsFavorite)*/
}

func GetUserInfoByID(userID uint, currentUserID uint) (*module.User, error) {
	var user module.User
	if err := DB.Model(&module.User{}).First(&user, userID).Error; err != nil {
		return nil, errors.New("get user info by id failed")
	}
	checkIsFollow(&user, currentUserID)
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
	DB.Model(&module.User{}).Where("id = ?", video.UserID).Update("work_count", gorm.Expr("work_count + ?", 1))
	return nil
}

func GetPublishVideosByUserID(userID uint) (*[]module.Video, error) {
	var videos []module.Video
	// 由于都是同一个用户的视频, 所以不使用.Joins("Author"), 只查询一次用户的信息, 再批量添加, 该逻辑需要有调用本接口的函数处理
	if err := DB.Model(&module.Video{}).Find(&videos, "user_id = ?", userID).Error; err != nil {
		return nil, errors.New("该用户未发布视频")
	}
	checkIsFavorite(&videos, userID)
	return &videos, nil
}

func GetFavoriteVideosByUserID(userID uint, currentUserID uint) *[]module.Video {
	var videos []module.Video
	DB.Model(&module.Video{}).
		Joins("Author").
		Where("videos.id in (SELECT video_id  FROM like_videos WHERE user_id = ?)", userID).
		Find(&videos)
	checkIsFavorite(&videos, currentUserID)
	return &videos
}

func GetVideosBeforeTimestamp(timestamp time.Time, userID uint) (*[]module.Video, error) {
	/*
		var video module.Video
		db.Model(&module.Video{}).Joins("User").Where("user_id = ?", 4).Find(&video)
		fmt.Println(video.User.Name)
	*/

	/*	// 填充is_favorite
		db.Table("videos").
		Select("videos.*", "EXISTS(SELECT 1 FROM likes WHERE likes.user_id = ? AND likes.video_id = videos.id) AS is_favorite", userId).
		LeftJoin("likes", "videos.id = likes.video_id").
		Scan(&videos)
	*/

	/*  // 填充Author
	var video module.Video
	db.Model(&module.Video{}).Joins("Author").Where("user_id = ?", 4).Find(&video)
	fmt.Println(video.Author.Name)
	*/

	var videos []module.Video
	// 由于更可能是不同用户的视频, 所以使用.Joins("Author")
	if err := DB.Model(&module.Video{}).Joins("Author").Where("videos.created_at < ?", timestamp).Find(&videos).Error; err != nil {
		return nil, errors.New("没有在这之前的视频")
	}

	// 判断是否点赞
	checkIsFavorite(&videos, userID)
	// 判断是否关注
	checkIsFollow(&videos, userID)

	return &videos, nil
}

func FavoriteAction(userID uint, videoID uint, action module.FavoriteAction) error {
	var isExists bool
	DB.Model(&module.LikeVideos{}).
		Select("EXISTS(SELECT 1 FROM like_videos WHERE user_id = ? AND video_id = ?) AS is_exists", userID, videoID).
		Scan(&isExists)
	if isExists && action == module.Unlike {
		DB.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&module.LikeVideos{})
		DB.Model(&module.Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		DB.Model(&module.User{}).Where("id = (select user_id from videos where id = ?)", videoID).Update("total_favorited", gorm.Expr("total_favorited - ?", 1))
		DB.Model(&module.User{}).Where("id = ?", userID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		return nil
	} else if !isExists && action == module.Like {
		DB.Model(&module.LikeVideos{}).Create(&module.LikeVideos{UserID: userID, VideoID: videoID})
		DB.Model(&module.Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		DB.Model(&module.User{}).Where("id = (select user_id from videos where id = ?)", videoID).Update("total_favorited", gorm.Expr("total_favorited + ?", 1))
		DB.Model(&module.User{}).Where("id = ?", userID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		return nil
	}
	return errors.New("favorite action error")
}

func GetCommentsByVideoID(videoID uint, currentUserID uint) *[]module.Comment {
	var comments []module.Comment
	DB.Model(&module.Comment{}).Joins("User").Where("video_id = ?", videoID).Order("created_at desc").Find(&comments)
	checkIsFollow(&comments, currentUserID)
	return &comments
}

func CommentAction(comment *module.Comment, action module.CommentAction) error {
	if action == module.PublishComment {
		if err := DB.Model(&module.Comment{}).Create(comment).Error; err != nil {
			return errors.New("添加评论失败")
		}
		DB.Model(&module.Comment{}).Joins("User").First(comment)
		DB.Model(&module.Video{}).Where("id = ?", comment.VideoID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
	} else {
		if err := DB.Delete(&module.Comment{}, comment.Id).Error; err != nil {
			return errors.New("删除评论失败")
		}
		DB.Model(&module.Video{}).Where("id = ?", comment.VideoID).Update("comment_count", gorm.Expr("comment_count - ?", 1))
	}

	return nil
}

func RelationAction(toUserID uint, currentUserID uint, action module.RelationAction) error {
	var isExists bool
	DB.Table("users_follows").
		Select("EXISTS(SELECT 1 FROM users_follows WHERE user_id = ? AND follow_id = ?) AS is_exists", currentUserID, toUserID).
		Scan(&isExists)
	if isExists && action == module.UnFollow {
		DB.Where("user_id = ? AND follow_id = ?", currentUserID, toUserID).Delete(&module.UsersFollows{})
		DB.Model(&module.User{}).Where("id = ?", toUserID).Update("follower_count", gorm.Expr("follower_count - ?", 1))
		DB.Model(&module.User{}).Where("id = ?", currentUserID).Update("follow_count", gorm.Expr("follow_count - ?", 1))
		return nil
	} else if !isExists && action == module.Follow {
		DB.Create(&module.UsersFollows{UserID: currentUserID, FollowID: toUserID})
		DB.Model(&module.User{}).Where("id = ?", toUserID).Update("follower_count", gorm.Expr("follower_count + ?", 1))
		DB.Model(&module.User{}).Where("id = ?", currentUserID).Update("follow_count", gorm.Expr("follow_count + ?", 1))
		return nil
	}
	return errors.New("favorite action error")
}

func GetFollowingListByUserID(userID uint) *[]module.User {
	var followings []module.User
	DB.Model(&module.User{}).
		Where("users.id in (SELECT follow_id FROM users_follows WHERE user_id = ?)", userID).
		Find(&followings)
	for i := 0; i < len(followings); i++ {
		followings[i].IsFollow = true
	}
	return &followings
}

func GetFollowerListByUserID(userID uint, currentUserID uint) *[]module.User {
	var followings []module.User
	DB.Model(&module.User{}).
		Where("users.id in (SELECT user_id FROM users_follows WHERE follow_id = ?)", userID).
		Find(&followings)
	checkIsFollow(&followings, currentUserID)
	return &followings
}

func GetFriendListByUserID(userID uint, currentUserID uint) *[]module.User {
	var friends []module.User

	DB.Table("users_follows").
		Select("u.*").
		Joins("JOIN users AS u ON users_follows.follow_id = u.id").
		Where("users_follows.user_id = ?", userID).
		Scan(&friends)
	for i := 0; i < len(friends); i++ {
		friends[i].IsFollow = true
	}
	return &friends
}

func checkIsFavorite(videos *[]module.Video, userID uint) {
	if userID == 0 {
		return
	}

	for i := 0; i < len(*videos); i++ {
		DB.Table("like_videos").
			Select("EXISTS(SELECT 1 FROM like_videos WHERE user_id = ? AND video_id = ?) AS is_favorite", userID, (*videos)[i].ID).
			Scan(&(*videos)[i].IsFavorite)
	}
}

func checkIsFollow(videos any, userID uint) {
	if userID == 0 {
		return
	}

	switch x := videos.(type) {
	case *[]module.Video:
		for i := 0; i < len(*x); i++ {
			DB.Table("users_follows").
				Select("EXISTS(SELECT 1 FROM users_follows WHERE user_id = ? AND follow_id = ?) AS is_follow", userID, (*x)[i].Author.ID).
				Scan(&(*x)[i].Author.IsFollow)
		}

	case *module.User:
		if userID == (*x).ID {
			return
		}
		DB.Table("users_follows").
			Select("EXISTS(SELECT 1 FROM users_follows WHERE user_id = ? AND follow_id = ?) AS is_follow", userID, (*x).ID).
			Scan(&(*x).IsFollow)
	case *[]module.User:
		for i := 0; i < len(*x); i++ {
			DB.Table("users_follows").
				Select("EXISTS(SELECT 1 FROM users_follows WHERE user_id = ? AND follow_id = ?) AS is_follow", userID, (*x)[i].ID).
				Scan(&(*x)[i].IsFollow)
		}
	case *[]module.Comment:
		for i := 0; i < len(*x); i++ {
			DB.Table("users_follows").
				Select("EXISTS(SELECT 1 FROM users_follows WHERE user_id = ? AND follow_id = ?) AS is_follow", userID, (*x)[i].User.ID).
				Scan(&(*x)[i].User.IsFollow)
		}
	}
}
