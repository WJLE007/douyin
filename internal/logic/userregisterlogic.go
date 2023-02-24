package logic

import (
	"context"
	"douyin-zero/common/errorx"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"
	"douyin-zero/util"
	"strconv"

	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisteRequest) (resp *types.UserRegisterResponse, err error) {
	// 查询用户是否存在
	_, err = dal.GetUserByName(l.ctx, l.svcCtx, req.Username)
	if err == nil {
		return nil, errorx.NewDefaultError("用户名已存在")
	}

	// 创建用户
	userID, err := dal.CreateUser(l.ctx, l.svcCtx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	token := util.UUID()
	// 保存token到redis
	l.svcCtx.RedisClient.Set(l.ctx, consts.TokenPrefix+token, strconv.Itoa(int(userID)), consts.TokenTTL)

	resp = &types.UserRegisterResponse{
		BaseResponse: types.BaseResponse{
			Message: "注册成功",
		},
		UserID: userID,
		Token:  token,
	}

	return
}
