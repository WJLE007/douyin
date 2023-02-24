package logic

import (
	"context"
	"douyin-zero/internal/dal"
	"douyin-zero/util"

	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCommentListLogic {
	return &VideoCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoCommentListLogic) VideoCommentList(req *types.VideoCommentListRequest) (resp *types.VideoCommentListResponse, err error) {
	// 判断视频是否存在
	_, err = dal.GetVideoByID(l.ctx, l.svcCtx, req.VideoID)
	if err != nil {
		return nil, err
	}
	commentList, err := dal.GetCommentList(l.ctx, l.svcCtx, req.VideoID)

	commentListDTO := make([]*types.Comment, 0, len(commentList))

	userInfoLogic := NewUserInfoLogic(l.ctx, l.svcCtx)

	for i := range commentList {

		commentDTO, err1 := util.CopyStruct[types.Comment](commentList[i])
		if err1 != nil {
			err = err1
			return
		}

		// 填充user信息
		userInfoResponse, err1 := userInfoLogic.UserInfo(&types.UserInfoRequest{UserID: commentList[i].UserID})
		if err1 != nil {
			err = err1
			return
		}

		commentDTO.User = userInfoResponse.User

		commentDTO.Create_date = commentList[i].CreatedAt.Format("2006-01-02 15:04")

		commentListDTO = append(commentListDTO, commentDTO)

	}
	resp = &types.VideoCommentListResponse{
		CommentList: commentListDTO,
	}

	return
}
