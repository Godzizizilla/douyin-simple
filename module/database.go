package module

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `json:"id,omitempty" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name            string `json:"name,omitempty" gorm:"unique"` // 名字
	FollowCount     uint   `json:"follow_count,omitempty"`       // 关注总数
	FollowerCount   uint   `json:"follower_count,omitempty"`     // 粉丝总数
	Avatar          string `json:"avatar,omitempty"`             // 用户头像
	BackgroundImage string `json:"background_image,omitempty"`   // 用户个人页顶部大图
	Signature       string `json:"signature,omitempty"`          // 个人简介
	TotalFavorited  uint   `json:"total_favorited,omitempty"`    // 获赞数量
	WorkCount       uint   `json:"work_count,omitempty"`         // 作品数
	FavoriteCount   uint   `json:"favorite_count,omitempty"`     // 喜欢数
	IsFollow        bool   `json:"is_follow,omitempty" gorm:"-"` // 是否关注

	EncryptedPassword string `json:"-"` // 哈希密码

	// many to many
	Followings []*User `json:"-" gorm:"many2many:users_follows;foreignKey:ID;joinForeignKey:UserId;References:ID;joinReferences:FollowID"`
	Followers  []*User `json:"-" gorm:"many2many:users_follows;foreignKey:ID;joinForeignKey:FollowID;References:ID;joinReferences:UserId"`
	LikeVideos []Video `json:"-" gorm:"many2many:like_videos"`

	// has many
	MyVideos []Video `json:"-" gorm:"foreignKey:UserID"`
}

type Video struct {
	ID        uint           `json:"id,omitempty" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Title         string `json:"title,omitempty"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount uint   `json:"favorite_count,omitempty"`
	CommentCount  uint   `json:"comment_count,omitempty"`

	UserID uint `json:"-"`
	Author User `gorm:"foreignKey:UserID"`

	IsFavorite bool `json:"is_favorite,omitempty" gorm:"-"`
}

type LikeVideos struct {
	UserID  uint
	VideoID uint
}

type Comment struct {
	Id         uint      `json:"id,omitempty"`
	Content    string    `json:"content,omitempty"`
	CreateDate string    `json:"create_date,omitempty"`
	CreatedAt  time.Time `json:"-"`
	UserID     uint      `json:"-"`
	User       User      `json:"user"`
	VideoID    uint      `json:"-"`
	Video      Video     `json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.EncryptedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.EncryptedPassword = string(hashedPassword)
	return nil
}
