package logic

import (
	"context"
	"douyin-zero/common/errorx"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"
	"douyin-zero/internal/dal/model"
	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"
	"douyin-zero/util"
	"regexp"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoCommentActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCommentActionLogic {
	return &VideoCommentActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoCommentActionLogic) VideoCommentAction(req *types.VideoCommentActionRequest) (resp *types.VideoCommentActionResponse, err error) {
	// 判断视频是否存在
	_, err = dal.GetVideoByID(l.ctx, l.svcCtx, req.VideoID)
	if err != nil {
		return nil, err
	}

	switch req.ActionType {
	case consts.PublishCommentActionType:
		return l.PublishComment(req)
	case consts.DeleteCommentActionType:
		return l.DeleteComment(req)
	default:
		return nil, nil
	}
}

func (l *VideoCommentActionLogic) PublishComment(req *types.VideoCommentActionRequest) (resp *types.VideoCommentActionResponse, err error) {

	content := l.removeTag(req.CommentText)
	if content == "" {
		return nil, errorx.NewDefaultError("输入内容不合法哦")
	}

	// 添加评论
	userID := l.ctx.Value(consts.UserIDKey).(int64)
	createdAt := time.Now()
	comment := &model.Comment{
		UserID:    userID,
		VideoID:   req.VideoID,
		Content:   req.CommentText,
		CreatedAt: createdAt,
	}
	err = dal.CreateComment(l.ctx, l.svcCtx, comment)
	if err != nil {
		return nil, err
	}

	commentDTO, err := util.CopyStruct[types.Comment](comment)
	if err != nil {
		return nil, err
	}

	// 填充user信息
	userInfoLogic := NewUserInfoLogic(l.ctx, l.svcCtx)
	userInfoResponse, err := userInfoLogic.UserInfo(&types.UserInfoRequest{UserID: userID})
	if err != nil {
		return nil, err
	}
	commentDTO.User = userInfoResponse.User
	commentDTO.Create_date = createdAt.Format("2006-01-02 15:04")

	resp = &types.VideoCommentActionResponse{
		BaseResponse: types.BaseResponse{Message: "评论成功"},
		Comment:      commentDTO,
	}

	return
}

func (l *VideoCommentActionLogic) removeTag(content string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	content = re.ReplaceAllString(content, "")
	content = strings.Replace(content, " ", "", -1)
	return content
}

func (l *VideoCommentActionLogic) filterSensitiveWords(content string) string {
	// 敏感词列表，可以从数据库或者配置文件中读取
	sensitiveWords := []string{"骂人的词汇", "政治相关的词汇"}

	// 替换敏感词汇
	for _, word := range sensitiveWords {
		content = strings.ReplaceAll(content, word, "***")
	}

	return content
}

func (l *VideoCommentActionLogic) DeleteComment(req *types.VideoCommentActionRequest) (resp *types.VideoCommentActionResponse, err error) {
	err = dal.DeleteComment(l.ctx, l.svcCtx, req.CommentID)
	if err != nil {
		return nil, err
	}
	resp = &types.VideoCommentActionResponse{
		BaseResponse: types.BaseResponse{
			Message: "删除成功",
		},
	}
	return
}
