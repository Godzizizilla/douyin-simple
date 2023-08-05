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

type RelationAction string

const (
	Follow   RelationAction = "1"
	UnFollow RelationAction = "2"
)
