package logic

import (
	"context"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"
	"douyin-zero/internal/dal/model"
	"github.com/duke-git/lancet/v2/slice"
	"strconv"

	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListRequest) (resp *types.FriendListResponse, err error) {
	ctxUserID := l.ctx.Value(consts.UserIDKey).(int64)
	// 获取好友列表
	friendIDs, err := dal.GetFriendIDs(l.ctx, l.svcCtx, ctxUserID)
	if err != nil {
		return nil, err
	}
	// 根据粉丝好友id查询聊天最后一条消息的列表
	ctxUserIDStr := strconv.FormatInt(ctxUserID, 10)

	// 获取所有好友的最后一条消息
	allLastMessage := l.svcCtx.RedisClient.HGetAll(l.ctx, consts.MessagePrefix+ctxUserIDStr).Val()

	// 遍历好友id，从redis中拿到每一个好友的最后一条消息
	lastMessages := make([]*model.Message, 0, len(friendIDs))
	for _, friendID := range friendIDs {
		lastMessage := &model.Message{}

		friendIDStr := strconv.FormatInt(friendID, 10)
		err = sonic.UnmarshalString(allLastMessage[friendIDStr], lastMessage)
		if err != nil {
			return
		}

		lastMessages = append(lastMessages, lastMessage)
	}

	// 将lastMessages按照时间排序
	slice.SortBy(lastMessages, func(i, j *model.Message) bool {
		return i.CreateTime > j.CreateTime
	})

	// 初始化返回好友列表
	friendList := make([]*types.FriendUser, 0, len(friendIDs))

	// 获取已排序好的lastMessages中的所有用户信息
	userInfoLogic := NewUserInfoLogic(l.ctx, l.svcCtx)

	for _, lastMessage := range lastMessages {
		friend := &types.FriendUser{}
		// 填充消息内容
		friend.Message = lastMessage.Content
		// 填充User详情和消息类型
		if lastMessage.FromUserID == ctxUserID {
			// 我是发送者
			friend.MsgType = consts.MessageSend
			// 获取User信息详情
			userInfoResp, err1 := userInfoLogic.UserInfo(&types.UserInfoRequest{UserID: lastMessage.ToUserID})
			if err1 != nil {
				err = err1
				return
			}

			friend.User = *userInfoResp.User
		} else {
			// 我是接收者
			friend.MsgType = consts.MessageReceive
			// 获取User信息详情
			userInfoResp, err1 := userInfoLogic.UserInfo(&types.UserInfoRequest{UserID: lastMessage.FromUserID})
			if err1 != nil {
				err = err1
				return
			}

			friend.User = *userInfoResp.User
		}
		friendList = append(friendList, friend)
	}

	resp = &types.FriendListResponse{
		FriendList: friendList,
	}

	return
}
