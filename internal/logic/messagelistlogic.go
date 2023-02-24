package logic

import (
	"context"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"
	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"
	"douyin-zero/util"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageListLogic {
	return &MessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageListLogic) MessageList(req *types.MessageListRequest) (resp *types.MessageListResponse, err error) {
	// 判断对方是否存在
	_, err = dal.GetUserByID(l.ctx, l.svcCtx, req.ToUserID)
	if err != nil {
		return nil, err
	}

	ctxUserID := l.ctx.Value(consts.UserIDKey).(int64)

	// 获取消息列表
	messageList, err := dal.GetMessageList(l.ctx, l.svcCtx, ctxUserID, req.ToUserID, req.PreMsgTime)
	if err != nil {
		return nil, err
	}

	messageListDTO := make([]*types.Message, 0, len(messageList))
	for i := range messageList {
		messageDTO, err1 := util.CopyStruct[types.Message](messageList[i])
		if err1 != nil {
			err = err1
			return
		}

		messageListDTO = append(messageListDTO, messageDTO)
	}

	resp = &types.MessageListResponse{
		MessageList: messageListDTO,
	}

	// 这里睡眠一下
	time.Sleep(300 * time.Millisecond)

	return
}
