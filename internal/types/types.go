// Code generated by goctl. DO NOT EDIT.
package types

type BaseResponse struct {
	Code    int64  `json:"status_code"`
	Message string `json:"status_msg,omitempty"`
}

type RelationActionRequest struct {
	ToUserID   int64 `form:"to_user_id" validate:"required"`
	ActionType int64 `form:"action_type" validate:"required,oneof=1 2"`
}

type RelationActionResponse struct {
	BaseResponse
}

type RelationFollowListRequest struct {
	UserID int64 `form:"user_id" validate:"required"`
}

type RelationFollowListResponse struct {
	BaseResponse
	UserList []*User `json:"user_list"`
}

type RelationFollowerListRequest struct {
	UserID int64 `form:"user_id" validate:"required"`
}

type RelationFollowerListResponse struct {
	BaseResponse
	UserList []*User `json:"user_list"`
}

type Video struct {
	ID            int64  `json:"id"`
	Author        *User  `json:"author"`
	Title         string `json:"title"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
}

type VideoActionRequest struct {
	Title string `form:"title"`
}

type VideoActionResponse struct {
	BaseResponse
}

type VideoFeedRequest struct {
	LatestTime int64  `form:"latest_time,optional"`
	Token      string `form:"token,optional"`
}

type VideoFeedResponse struct {
	BaseResponse
	VideoList []*Video `json:"video_list"`
}

type VideoFavoriteActionRequest struct {
	VideoID    int64 `form:"video_id"`
	ActionType int64 `form:"action_type"`
}

type VideoFavoriteActionResponse struct {
	BaseResponse
}

type VideoListRequest struct {
	UserID int64 `form:"user_id,optional"`
}

type VideoListResponse struct {
	BaseResponse
	VideoList []*Video `json:"video_list"`
}

type VideoFavoriteListRequest struct {
	UserID int64 `form:"user_id,optional"`
}

type VideoFavoriteListResponse struct {
	BaseResponse
	VideoList []*Video `json:"video_list"`
}

type Comment struct {
	ID          int64  `json:"id"`
	User        *User  `json:"user"`
	Content     string `json:"content"`
	Create_date string `json:"create_date"`
}

type VideoCommentActionRequest struct {
	VideoID     int64  `form:"video_id"`
	ActionType  int64  `form:"action_type" validate:"required,oneof=1 2"`
	CommentText string `form:"comment_text,optional"`
	CommentID   int64  `form:"comment_id,optional"`
}

type VideoCommentActionResponse struct {
	BaseResponse
	Comment *Comment `json:"comment"`
}

type VideoCommentListRequest struct {
	VideoID int64 `form:"video_id"`
}

type VideoCommentListResponse struct {
	BaseResponse
	CommentList []*Comment `json:"comment_list"`
}

type User struct {
	ID              int64  `json:"id,string"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
	FollowCount     int64  `json:"follow_count"`
	TotalFavorited  int64  `json:"total_favorited"`
	Signature       string `json:"signature"`
	BackgroundImage string `json:"background_image"`
	FollowerCount   int64  `json:"follower_count"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
	IsFollow        bool   `json:"is_follow"`
}

type UserRegisteRequest struct {
	Username string `form:"username" validate:"required,min=0,max=32"`
	Password string `form:"password" validate:"required,min=0,max=32"`
}

type UserRegisterResponse struct {
	BaseResponse
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserLoginRequest struct {
	Username string `form:"username" validate:"required,min=0,max=32"`
	Password string `form:"password" validate:"required,min=0,max=32"`
}

type UserLoginResponse struct {
	BaseResponse
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfoRequest struct {
	UserID int64 `form:"user_id" validate:"required"`
}

type UserInfoResponse struct {
	BaseResponse
	User *User `json:"user"`
}

type Message struct {
	ID         int64  `json:"id"`
	FromUserID int64  `json:"from_user_id"`
	ToUserID   int64  `json:"to_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

type FriendUser struct {
	User
	MsgType int64  `json:"msg_type"  validate:"oneof=0 1"`
	Message string `json:"message,optional"`
}

type MessageActionRequest struct {
	ToUserID   int64  `form:"to_user_id"`
	ActionType int64  `form:"action_type" validate:"required,oneof=1"`
	Content    string `form:"content"`
}

type MessageActionResponse struct {
	BaseResponse
}

type MessageListRequest struct {
	ToUserID   int64 `form:"to_user_id"`
	PreMsgTime int64 `form:"pre_msg_time,optional"`
}

type MessageListResponse struct {
	BaseResponse
	MessageList []*Message `json:"message_list"`
}

type FriendListRequest struct {
	UserID int64 `form:"user_id"`
}

type FriendListResponse struct {
	BaseResponse
	FriendList []*FriendUser `json:"user_list"`
}
