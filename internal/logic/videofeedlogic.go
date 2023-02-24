package logic

import (
	"context"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"
	"douyin-zero/util"
	"time"

	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoFeedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoFeedLogic {
	return &VideoFeedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoFeedLogic) VideoFeed(req *types.VideoFeedRequest) (resp *types.VideoFeedResponse, err error) {
	if req.LatestTime == 0 {
		req.LatestTime = time.Now().UnixMilli()
	}

	// 获取视频列表
	videoList, err := dal.GetVideoListFeed(l.ctx, l.svcCtx, req.LatestTime)
	if err != nil {
		return
	}

	userID := l.ctx.Value(consts.UserIDKey).(int64)

	// 转换为dto
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
		isFavorite, err1 := dal.IsFavoriteVideo(l.ctx, l.svcCtx, video.ID, userID)

		video.IsFavorite = isFavorite

		videoListDTO = append(videoListDTO, video)
	}

	resp = &types.VideoFeedResponse{
		VideoList: videoListDTO,
	}

	return
}
