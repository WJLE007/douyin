package logic

import (
	"context"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"
	"douyin-zero/util"

	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoFavoriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoFavoriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoFavoriteListLogic {
	return &VideoFavoriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoFavoriteListLogic) VideoFavoriteList(req *types.VideoFavoriteListRequest) (resp *types.VideoFavoriteListResponse, err error) {
	ctxUserID := l.ctx.Value(consts.UserIDKey).(int64)

	// 获取用户喜欢的视频ids
	favoriteVideoIDs, err := dal.GetFavoriteVideoIDs(l.ctx, l.svcCtx, req.UserID)
	if err != nil {
		return
	}

	// 获取视频列表
	videoList, err := dal.GetVideoListByIDs(l.ctx, l.svcCtx, favoriteVideoIDs)
	if err != nil {
		return
	}

	// 获取用户信息
	userInfoLogic := NewUserInfoLogic(l.ctx, l.svcCtx)

	// 获取用户喜欢的视频列表
	videoListDTO := make([]*types.Video, 0, len(videoList))
	for i := range videoList {
		video, err1 := util.CopyStruct[types.Video](videoList[i])
		if err1 != nil {
			err = err1
			return
		}

		userInfoResponse, err1 := userInfoLogic.UserInfo(&types.UserInfoRequest{UserID: videoList[i].AuthorID})
		if err != nil {
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

	resp = &types.VideoFavoriteListResponse{
		VideoList: videoListDTO,
	}
	return
}
