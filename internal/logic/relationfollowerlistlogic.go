package logic

import (
	"context"
	"douyin-zero/internal/dal"

	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationFollowerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationFollowerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationFollowerListLogic {
	return &RelationFollowerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationFollowerListLogic) RelationFollowerList(req *types.RelationFollowerListRequest) (resp *types.RelationFollowerListResponse, err error) {
	// 查询该用户
	_, err = dal.GetUserByID(l.ctx, l.svcCtx, req.UserID)
	if err != nil {
		return
	}

	// 查询该用户的关注列表
	followerIDs, err := dal.GetFollowerIDs(l.ctx, l.svcCtx, req.UserID)
	if err != nil {
		return
	}

	// 查询用户信息
	userList := []*types.User{}
	userInfoLogic := NewUserInfoLogic(l.ctx, l.svcCtx)
	for _, followerID := range followerIDs {
		userInfoResp, err := userInfoLogic.UserInfo(&types.UserInfoRequest{UserID: followerID})
		if err != nil {
			return nil, err
		}

		userList = append(userList, userInfoResp.User)
	}

	resp = &types.RelationFollowerListResponse{
		UserList: userList,
	}

	return

}
