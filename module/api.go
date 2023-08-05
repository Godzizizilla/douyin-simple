package module

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserResponse struct {
	Response
	UserId uint   `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserInfoResponse struct {
	Response
	User User `json:"user"`
}

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

type ApiUser struct {
	ID              uint   `json:"id,omitempty"`                 // 用户ID
	Name            string `json:"name,omitempty"`               // 名字
	FollowCount     uint   `json:"follow_count,omitempty"`       // 关注总数
	FollowerCount   uint   `json:"follower_count,omitempty"`     // 粉丝总数
	Avatar          string `json:"avatar,omitempty"`             // 用户头像
	IsFollow        bool   `json:"is_follow,omitempty" gorm:"-"` // 是否关注
	BackgroundImage string `json:"background_image,omitempty"`   // 用户个人页顶部大图
	Signature       string `json:"signature,omitempty"`          // 个人简介
	TotalFavorited  uint   `json:"total_favorited,omitempty"`    // 获赞数量
	WorkCount       uint   `json:"work_count,omitempty"`         // 作品数
	FavoriteCount   uint   `json:"favorite_count,omitempty"`     // 喜欢数
}

type ApiVideo struct {
	ID            uint     `json:"id,omitempty"`
	Author        *ApiUser `json:"author,omitempty" gorm:"-"`
	PlayUrl       string   `json:"play_url,omitempty"`
	CoverUrl      string   `json:"cover_url,omitempty"`
	FavoriteCount uint     `json:"favorite_count,omitempty"`
	CommentCount  uint     `json:"comment_count,omitempty"`
	IsFavorite    bool     `json:"is_favorite,omitempty" gorm:"-"`
	Title         string   `json:"title,omitempty"`
}
