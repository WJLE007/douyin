package logic

import (
	"context"
	"douyin-zero/common/errorx"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"
	"douyin-zero/util"

	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoListLogic {
	return &VideoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoListLogic) VideoList(req *types.VideoListRequest) (resp *types.VideoListResponse, err error) {
	_, err = dal.GetUserByID(l.ctx, l.svcCtx, req.UserID)
	if err != nil {
		return nil, errorx.NewDefaultError("用户不存在")
	}

	videoList, err := dal.GetUserVideoList(l.ctx, l.svcCtx, req.UserID)
	if err != nil {
		return
	}

	ctxUserID := l.ctx.Value(consts.UserIDKey).(int64)

	userInfoLogic := NewUserInfoLogic(l.ctx, l.svcCtx)

	videoListDTO := make([]*types.Video, 0, len(videoList))
	for i := range videoList {
		video, err1 := util.CopyStruct[types.Video](videoList[i])
		if err1 != nil {
			err = err1
			return
		}

		// 填充author信息
		userInfoResponse, err1 := userInfoLogic.UserInfo(&types.UserInfoRequest{UserID: videoList[i].AuthorID})
		if err1 != nil {
			err = err1
			return
		}

		video.Author = userInfoResponse.User

		// 查询是否点赞
		isFavorite, err1 := dal.IsFavoriteVideo(l.ctx, l.svcCtx, video.ID, ctxUserID)
		if err1 != nil {
			err = err1
			return
		}

		video.IsFavorite = isFavorite

		videoListDTO = append(videoListDTO, video)
	}

	resp = &types.VideoListResponse{
		VideoList: videoListDTO,
	}

	return
}
