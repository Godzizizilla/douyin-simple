package database

//type User struct {
//	gorm.Model
//	Name            string `gorm:"unique" json:"name"` // 名字
//	FollowCount     uint   `json:"follow_count"`       // 关注总数
//	FollowerCount   uint   `json:"follower_count"`     // 粉丝总数
//	Avatar          string `json:"avatar"`             // 用户头像
//	BackgroundImage string `json:"background_image"`   // 用户个人页顶部大图
//	Signature       string `json:"signature"`          // 个人简介
//	TotalFavorited  uint   `json:"total_favorited"`    // 获赞数量
//	WorkCount       uint   `json:"work_count"`         // 作品数
//	FavoriteCount   uint   `json:"favorite_count"`     // 喜欢数
//
//	Videos       []Video `json:"-"` // 作品
//	PasswordHash string  `json:"-"` // 哈希密码
//}
//
//type Video struct {
//}

/*type CustomModel struct {
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
	Videos          []Video

	PasswordHash string
}

type Video struct {
	CustomModel
	UserID        int64
	Title         string
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
}*/

/*
type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}
*/
