package logic

import (
	"context"
	"douyin-zero/internal/dal"

	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationFollowListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationFollowListLogic {
	return &RelationFollowListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationFollowListLogic) RelationFollowList(req *types.RelationFollowListRequest) (resp *types.RelationFollowListResponse, err error) {
	// 查询该用户
	_, err = dal.GetUserByID(l.ctx, l.svcCtx, req.UserID)
	if err != nil {
		return
	}

	// 查询该用户的关注列表
	followIDs, err := dal.GetFollowIDs(l.ctx, l.svcCtx, req.UserID)
	if err != nil {
		return
	}

	// 查询用户信息
	userList := []*types.User{}
	userInfoLogic := NewUserInfoLogic(l.ctx, l.svcCtx)
	for _, followID := range followIDs {
		userInfoResp, err := userInfoLogic.UserInfo(&types.UserInfoRequest{UserID: followID})
		if err != nil {
			return nil, err
		}
		userList = append(userList, userInfoResp.User)
	}

	resp = &types.RelationFollowListResponse{
		UserList: userList,
	}

	return
}
