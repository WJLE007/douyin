package logic

import (
	"context"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"
	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"
	"douyin-zero/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	// 查询该用户
	user, err := dal.GetUserByID(l.ctx, l.svcCtx, req.UserID)
	if err != nil {
		return
	}

	tUser, err := util.CopyStruct[types.User](user)
	if err != nil {
		return
	}

	// 查询ctxUser是否关注user
	ctxUser := l.ctx.Value(consts.UserIDKey)
	var ctxUserID int64
	if ctxUser != nil {
		ctxUserID = ctxUser.(int64)
	}

	isFollow, err := dal.IsFollow(l.ctx, l.svcCtx, ctxUserID, user.ID)
	if err != nil {
		return nil, err
	}

	tUser.IsFollow = isFollow

	resp = &types.UserInfoResponse{
		User: tUser,
	}
	return
}
