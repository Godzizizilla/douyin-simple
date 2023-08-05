package module

type FavoriteAction string

const (
	Like   FavoriteAction = "1"
	Unlike FavoriteAction = "2"
)

type CommentAction string

const (
	PublishComment CommentAction = "1"
	DeleteComment  CommentAction = "2"
)
