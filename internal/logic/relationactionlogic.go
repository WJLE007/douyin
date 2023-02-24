package logic

import (
	"context"
	"douyin-zero/common/errorx"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"

	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelationActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationActionLogic {
	return &RelationActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelationActionLogic) RelationAction(req *types.RelationActionRequest) (resp *types.RelationActionResponse, err error) {
	// 查询该用户
	_, err = dal.GetUserByID(l.ctx, l.svcCtx, req.ToUserID)
	if err != nil {
		return
	}

	ctxUserID := l.ctx.Value(consts.UserIDKey).(int64)

	switch req.ActionType {
	case 1:
		// 关注
		err = dal.CreateFollow(l.ctx, l.svcCtx, ctxUserID, req.ToUserID)
		// 检验对方是否关注了自己
		isFollow, _ := dal.IsFollow(l.ctx, l.svcCtx, req.ToUserID, ctxUserID)
		if isFollow {
			// 如果互相关注了，那么我自动给对方发条消息
			messageActionLogic := NewMessageActionLogic(l.ctx, l.svcCtx)
			_, _ = messageActionLogic.MessageAction(&types.MessageActionRequest{
				ToUserID:   req.ToUserID,
				ActionType: consts.MessageActionSend,
				Content:    "嗨~ 我们已成为好友，来聊聊天吧",
			})
		}
	case 2:
		// 取消关注
		err = dal.DeleteFollow(l.ctx, l.svcCtx, ctxUserID, req.ToUserID)
	default:
		err = errorx.NewDefaultError("参数错误")
	}

	if err != nil {
		return
	}

	resp = &types.RelationActionResponse{
		BaseResponse: types.BaseResponse{
			Code:    0,
			Message: "操作成功",
		},
	}

	return
}
