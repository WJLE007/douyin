package consts

const (
	TokenPrefix = "token:"
	TokenTTL    = -1

	URLPrefix                 = "http://"+"host:port"
	DefaultAvatarURL          = URLPrefix + "/avatar/default.jpg"
	DefaultBackgroundImageURL = URLPrefix + "/background/default.jpg"

	PlayFilePath   = "D:\\Desktop\\nginx-1.18.0\\html\\douyin\\play"
	CoverFilePath  = "D:\\Desktop\\nginx-1.18.0\\html\\douyin\\cover"
	PlayURLPrefix  = URLPrefix + "/play/"
	CoverURLPrefix = URLPrefix + "/cover/"

	UserIDKey = "userID"

	MaxVideoSize = 1024 * 1024 * 500 // 500MB

	FollowPrefix   = "follow:"
	FollowerPrefix = "follower:"

	VideoFeedCount = 15

	VideoFavorite      = "video_favorite"
	VideoFavoriteCount = "video_favorite_count"
	VideoCommentCount  = "video_comment_count"

	UserFavoriteCount  = "user_favorite_count"
	UserFavoritedCount = "user_favorited_count"

	MessagePrefix = "message:"

	FavoriteActionType   = 1
	UnFavoriteActionType = 2

	PublishCommentActionType = 1
	DeleteCommentActionType  = 2

	MessageActionSend = 1

	MessageReceive = 0
	MessageSend    = 1

	SyncTime = 5
)
