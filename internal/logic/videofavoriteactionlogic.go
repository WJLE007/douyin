package logic

import (
	"context"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"
	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoFavoriteActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoFavoriteActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoFavoriteActionLogic {
	return &VideoFavoriteActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoFavoriteActionLogic) VideoFavoriteAction(req *types.VideoFavoriteActionRequest) (resp *types.VideoFavoriteActionResponse, err error) {
	userID := l.ctx.Value(consts.UserIDKey).(int64)

	isFavorite := 0

	switch req.ActionType {
	case consts.FavoriteActionType:
		isFavorite = 1
	case consts.UnFavoriteActionType:
		isFavorite = 0
	}

	err = dal.VideoFavoriteAction(l.ctx, l.svcCtx, req.VideoID, userID, isFavorite)
	if err != nil {
		return
	}

	resp = &types.VideoFavoriteActionResponse{}

	return
}
