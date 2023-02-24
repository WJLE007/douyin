package logic

import (
	"context"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"
	"douyin-zero/internal/dal/model"
	"time"

	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageActionLogic {
	return &MessageActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageActionLogic) MessageAction(req *types.MessageActionRequest) (resp *types.MessageActionResponse, err error) {
	// 判断对方是否存在
	_, err = dal.GetUserByID(l.ctx, l.svcCtx, req.ToUserID)
	if err != nil {
		return nil, err
	}

	switch req.ActionType {
	case consts.MessageActionSend:
		// 发送消息
		return l.SendMessage(req)
	}

	return
}

func (l *MessageActionLogic) SendMessage(req *types.MessageActionRequest) (resp *types.MessageActionResponse, err error) {
	ctxUserID := l.ctx.Value("userID").(int64)
	message := &model.Message{
		FromUserID: ctxUserID,
		ToUserID:   req.ToUserID,
		Content:    req.Content,
		CreateTime: time.Now().UnixMilli(),
	}
	err = dal.CreateMessage(l.ctx, l.svcCtx, message)
	if err != nil {
		return nil, err
	}

	resp = &types.MessageActionResponse{}
	return
}
