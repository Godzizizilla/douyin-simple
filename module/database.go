package module

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name            string `gorm:"unique" json:"name"` // 名字
	FollowCount     uint   `json:"follow_count"`       // 关注总数
	FollowerCount   uint   `json:"follower_count"`     // 粉丝总数
	Avatar          string `json:"avatar"`             // 用户头像
	BackgroundImage string `json:"background_image"`   // 用户个人页顶部大图
	Signature       string `json:"signature"`          // 个人简介
	TotalFavorited  uint   `json:"total_favorited"`    // 获赞数量
	WorkCount       uint   `json:"work_count"`         // 作品数
	FavoriteCount   uint   `json:"favorite_count"`     // 喜欢数

	EncryptedPassword string `json:"-"` // 哈希密码

	// many to many
	Followings []*User `gorm:"many2many:users_follows;foreignKey:ID;joinForeignKey:UserId;References:ID;joinReferences:FollowID"`
	Followers  []*User `gorm:"many2many:users_follows;foreignKey:ID;joinForeignKey:FollowID;References:ID;joinReferences:UserId"`
	LikeVideos []Video `gorm:"many2many:like_videos"`

	// has many
	MyVideos []Video `gorm:"foreignKey:UserID"`
}

type Video struct {
	gorm.Model
	Title         string
	PlayUrl       string
	CoverUrl      string
	FavoriteCount uint
	CommentCount  uint

	UserID uint
	User   User
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.EncryptedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.EncryptedPassword = string(hashedPassword)
	return nil
}
