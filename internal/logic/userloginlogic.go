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

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.UserLoginRequest) (resp *types.UserLoginResponse, err error) {
	// 查询用户是否存在
	user, err := dal.GetUserByName(l.ctx, l.svcCtx, req.Username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if !util.VerifyPassword(req.Password, user.Password) {
		return nil, errorx.NewDefaultError("密码错误")
	}

	// 生成token，保存到redis
	token := util.UUID()
	l.svcCtx.RedisClient.Set(l.ctx, consts.TokenPrefix+token, strconv.Itoa(int(user.ID)), consts.TokenTTL)

	resp = &types.UserLoginResponse{
		BaseResponse: types.BaseResponse{
			Message: "登录成功",
		},
		UserID: user.ID,
		Token:  token,
	}

	return
}
