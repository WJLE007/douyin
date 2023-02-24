package dal

import (
	"douyin-zero/common/errorx"
	"douyin-zero/consts"
	"douyin-zero/internal/dal/model"
	"douyin-zero/internal/dal/query"
	"douyin-zero/internal/svc"
	"douyin-zero/util"
	"github.com/bytedance/sonic"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"strconv"
)

func GetFriendIDs(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) ([]int64, error) {
	// 获取好友列表
	f1 := query.Follow.As("f1")
	f2 := query.Follow.As("f2")
	friendIDs := make([]int64, 0)
	err := f1.Select(f1.FollowID.As("friend_id")).Join(f2, f1.UserID.EqCol(f2.FollowID), f1.FollowID.EqCol(f2.UserID)).Where(f1.UserID.Eq(userID)).Scan(&friendIDs)
	if err != nil {
		return nil, err
	}
	return friendIDs, nil
}

func GetMessageList(ctx context.Context, svcCtx *svc.ServiceContext, ctxUserID int64, toUserID int64, preMsgTime int64) ([]*model.Message, error) {
	// 获取聊天记录
	//m := query.Message
	//return m.Where(m.Where(m.ToUserID.Eq(toUserID), m.FromUserID.Eq(ctxUserID)).Or(m.Where(m.ToUserID.Eq(ctxUserID), m.FromUserID.Eq(toUserID))), m.CreateTime.Gt(preMsgTime)).
	//	Order(m.CreateTime).Find()

	sql := `SELECT * FROM message 
         WHERE ((to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)) 
        	AND create_time > ?
            ORDER BY create_time`

	//logSql := "SELECT * FROM message WHERE ((to_user_id = %v AND from_user_id = %v) OR (to_user_id = %v AND from_user_id = %v)) AND create_time > %v ORDER BY create_time"
	//
	//logx.Infof(logSql, ctxUserID, toUserID, toUserID, ctxUserID, preMsgTime)

	messageList := []*model.Message{}
	err := svcCtx.DB.Raw(sql, ctxUserID, toUserID, toUserID, ctxUserID, preMsgTime).Scan(&messageList).Error
	return messageList, err
}

func CreateMessage(ctx context.Context, svcCtx *svc.ServiceContext, message *model.Message) (err error) {
	// 添加消息
	m := query.Message
	err = m.Create(message)
	if err != nil {
		return err
	}
	// 保存消息到redis
	fromUserID := strconv.Itoa(int(message.FromUserID))
	toUserID := strconv.Itoa(int(message.ToUserID))

	messageValue, err := sonic.MarshalString(message)
	if err != nil {
		return err
	}

	svcCtx.RedisClient.HSet(ctx, consts.MessagePrefix+fromUserID, toUserID, messageValue)
	svcCtx.RedisClient.HSet(ctx, consts.MessagePrefix+toUserID, fromUserID, messageValue)
	return
}

func GetCommentList(ctx context.Context, svcCtx *svc.ServiceContext, videoID int64) ([]*model.Comment, error) {
	c := query.Comment
	return c.Where(c.VideoID.Eq(videoID)).Order(c.CreatedAt.Desc()).Find()
}

func CreateComment(ctx context.Context, svcCtx *svc.ServiceContext, comment *model.Comment) (err error) {
	err = query.Q.Transaction(func(tx *query.Query) error {
		// 添加评论
		c := tx.Comment
		err = c.Create(comment)
		if err != nil {

		}
		// 视频评论数+1
		v := tx.Video
		_, err = v.Where(v.ID.Eq(comment.VideoID)).UpdateSimple(v.CommentCount.Add(1))
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func DeleteComment(ctx context.Context, svcCtx *svc.ServiceContext, commentID int64) (err error) {
	err = query.Q.Transaction(func(tx *query.Query) error {
		c := tx.Comment
		// 获取评论
		comment, err := c.Where(c.ID.Eq(commentID)).First()
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		// 添加评论
		_, err = c.Where(c.ID.Eq(commentID)).Delete()
		if err != nil {

		}
		// 视频评论数-1
		v := tx.Video
		_, err = v.Where(v.ID.Eq(comment.VideoID)).UpdateSimple(v.CommentCount.Sub(1))
		if err != nil {
			return err
		}
		return nil
	})
	return

}

func GetVideoByID(ctx context.Context, svcCtx *svc.ServiceContext, videoID int64) (video *model.Video, err error) {
	v := query.Video
	video, err = v.Where(v.ID.Eq(videoID)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.NewDefaultError("视频不存在")
		}
		return nil, err
	}

	return
}

func GetVideoListByIDs(ctx context.Context, svcCtx *svc.ServiceContext, videoIDs []int64) (videoList []*model.Video, err error) {
	v := query.Video
	videoList, err = v.Where(v.ID.In(videoIDs...)).Find()
	slice.SortBy(videoList, func(i, j *model.Video) bool {
		return i.CreateTime > j.CreateTime
	})
	return
}

func GetFavoriteVideoIDs(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) (favoriteVideoIDs []int64, err error) {
	f := query.Favorite
	favoriteList, err := f.Where(f.UserID.Eq(userID)).Find()
	if err != nil {
		return
	}

	favoriteVideoIDs = make([]int64, 0, len(favoriteList))
	for i := range favoriteList {
		favoriteVideoIDs = append(favoriteVideoIDs, favoriteList[i].VideoID)
	}

	return
}

func IsFavoriteVideo(ctx context.Context, svcCtx *svc.ServiceContext, videoID int64, userID int64) (bool, error) {
	// 先查缓存
	key := strconv.Itoa(int(videoID)) + "-" + strconv.Itoa(int(userID))
	isFavorite := svcCtx.RedisClient.HGet(ctx, consts.VideoFavorite, key).Val()
	logx.Infof("videoID:%v,isFavorite:%v", videoID, isFavorite)
	if isFavorite != "" {
		return isFavorite == "1", nil
	}

	// 缓存中没有查数据库
	f := query.Favorite
	count, err := f.Where(f.UserID.Eq(userID), f.VideoID.Eq(videoID)).Count()
	return count != 0, err
}

func VideoFavoriteAction(ctx context.Context, svcCtx *svc.ServiceContext, videoID int64, userID int64, isFavorite int) error {
	// 先查视频是否存在
	v := query.Video
	video, err := v.Where(v.ID.Eq(videoID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.NewDefaultError("视频不存在")
		}
	}

	// 缓存到redis
	// videoID-userID -> isFavorite
	key := strconv.Itoa(int(videoID)) + "-" + strconv.Itoa(int(userID))
	value := isFavorite

	// 先看是否已经点赞
	oldValue := svcCtx.RedisClient.HGet(ctx, consts.VideoFavorite, key).Val()
	if (oldValue == "" || oldValue == "0") && isFavorite == 0 {
		// 未点赞
		return errorx.NewDefaultError("你还未点赞")
	} else if (oldValue == "1") && isFavorite == 1 {
		// 已点赞
		return errorx.NewDefaultError("你已经点过赞了")
	}

	// 更新作者获赞数
	err = svcCtx.RedisClient.HSet(ctx, consts.VideoFavorite, key, value).Err()
	if err != nil {
		return err
	}

	// 更新redis的视频点赞数
	var diff int64
	if isFavorite == 1 {
		diff = 1
	} else {
		diff = -1
	}

	err = svcCtx.RedisClient.HIncrBy(ctx, consts.UserFavoritedCount, strconv.Itoa(int(video.AuthorID)), diff).Err()
	if err != nil {
		return err
	}

	// 更新视频点赞数
	err = svcCtx.RedisClient.HIncrBy(ctx, consts.VideoFavoriteCount, strconv.Itoa(int(videoID)), diff).Err()
	if err != nil {
		return err
	}
	// 更新用户点赞数
	err = svcCtx.RedisClient.HIncrBy(ctx, consts.UserFavoriteCount, strconv.Itoa(int(userID)), diff).Err()
	if err != nil {
		return err
	}

	return nil

}

func GetUserVideoList(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) ([]*model.Video, error) {
	v := query.Video
	videoList, err := v.Where(v.AuthorID.Eq(userID)).Order(v.CreateTime.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

func GetVideoListFeed(ctx context.Context, svcCtx *svc.ServiceContext, latestTime int64) ([]*model.Video, error) {
	v := query.Video
	videoList, err := v.Where(v.CreateTime.Lt(latestTime)).Order(v.CreateTime.Desc()).Limit(consts.VideoFeedCount).Find()
	if err != nil {
		return nil, err
	}

	// 从redis中获取视频列表
	//videoList := make([]*model.Video, 0, consts.VideoFeedCount)
	//
	//resVal := svcCtx.RedisClient.ZRevRangeByScoreWithScores(ctx, consts.FeedCommonPrefix, &redis.ZRangeBy{
	//	Min:    "0",
	//	Max:    strconv.Itoa(int(latestTime)),
	//	Offset: 0,
	//	Count:  consts.VideoFeedCount,
	//}).Val()
	//if len(resVal) == 0 {
	//	return []*model.Video{}, nil
	//}
	//
	//for _, val := range resVal {
	//	video := &model.Video{}
	//	err := sonic.UnmarshalString(val.Member.(string), video)
	//	if err != nil {
	//		return nil, errorx.NewDefaultError("获取视频列表失败")
	//	}
	//
	//	videoList = append(videoList, video)
	//}

	return videoList, nil
}

func CreateVideo(ctx context.Context, svcCtx *svc.ServiceContext, video *model.Video) error {
	err := query.Video.Create(video)
	if err != nil {
		return errors.Wrap(err, "创建视频失败")
	}

	u := query.User
	// 给作者增加视频数
	_, err = u.Where(u.ID.Eq(video.AuthorID)).Update(u.WorkCount, gorm.Expr("work_count + ?", 1))
	if err != nil {
		return errors.Wrap(err, "创建视频失败")
	}

	// 将点赞数缓存到redis
	err = svcCtx.RedisClient.HSet(ctx, consts.VideoFavoriteCount, strconv.Itoa(int(video.ID)), 0).Err()
	if err != nil {
		return errors.Wrap(err, "创建视频失败")
	}
	// 将评论数缓存到redis
	err = svcCtx.RedisClient.HSet(ctx, consts.VideoCommentCount, strconv.Itoa(int(video.ID)), 0).Err()
	if err != nil {
		return errors.Wrap(err, "创建视频失败")
	}

	// 缓存到redis
	//videoJson, err := sonic.MarshalString(video)
	//
	//if err != nil {
	//	return errors.Wrap(err, "创建视频失败")
	//}
	//err = svcCtx.RedisClient.ZAdd(ctx, consts.FeedCommonPrefix, &redis.Z{
	//	Score:  float64(video.CreatedAt),
	//	Member: videoJson,
	//}).Err()

	if err != nil {
		return errors.Wrap(err, "创建视频失败")
	}
	return nil
}

func GetFollowerIDs(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) ([]int64, error) {
	// 从redis中获取粉丝id
	idStrs, err := svcCtx.RedisClient.SMembers(ctx, consts.FollowerPrefix+strconv.Itoa(int(userID))).Result()
	if err != nil {
		return nil, err
	}
	// 转为int64数组
	followerIDs := make([]int64, 0, len(idStrs))
	for _, idStr := range idStrs {
		idInt, _ := strconv.Atoi(idStr)
		followerIDs = append(followerIDs, int64(idInt))
	}

	//f := query.Follow
	//
	//followList, err := f.Where(f.FollowID.Eq(userID)).Find()
	//
	//if err != nil {
	//	return nil, errors.Wrap(err, "查询关注列表失败")
	//}
	//
	//var followerIDs []int64
	//for _, follow := range followList {
	//	followerIDs = append(followerIDs, follow.UserID)
	//}

	return followerIDs, nil
}

func GetFollowIDs(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) ([]int64, error) {
	f := query.Follow

	followList, err := f.Where(f.UserID.Eq(userID)).Find()

	if err != nil {
		return nil, errors.Wrap(err, "查询关注列表失败")
	}

	var followIDs []int64
	for _, follow := range followList {
		followIDs = append(followIDs, follow.FollowID)
	}

	return followIDs, nil
}

func CreateFollow(ctx context.Context, svcCtx *svc.ServiceContext, userID int64, followUserID int64) error {
	if userID == followUserID {
		return errorx.NewDefaultError("不能关注自己哦")
	}

	isFollow, err := IsFollow(ctx, svcCtx, userID, followUserID)
	if err != nil {
		return err
	}

	if isFollow {
		return errorx.NewDefaultError("您已经关注了该用户")
	}

	follow := &model.Follow{
		UserID:   userID,
		FollowID: followUserID,
	}

	// 事务
	err = query.Q.Transaction(func(tx *query.Query) error {
		// 创建关注
		err = tx.Follow.Create(follow)
		if err != nil {
			return err
		}

		// redis添加关注和被关注
		err = svcCtx.RedisClient.SAdd(ctx, consts.FollowPrefix+strconv.Itoa(int(userID)), followUserID).Err()
		if err != nil {
			return err
		}
		err = svcCtx.RedisClient.SAdd(ctx, consts.FollowerPrefix+strconv.Itoa(int(followUserID)), userID).Err()
		if err != nil {
			return err
		}

		u := tx.User

		// 更新用户关注数
		_, err = u.Where(u.ID.Eq(userID)).UpdateColumn(u.FollowCount, gorm.Expr("follow_count + ?", 1))
		if err != nil {
			return err
		}

		// 更新用户粉丝数
		_, err = u.Where(u.ID.Eq(followUserID)).UpdateColumn(u.FollowerCount, gorm.Expr("follower_count + ?", 1))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "关注失败")
	}

	return nil
}

func DeleteFollow(ctx context.Context, svcCtx *svc.ServiceContext, userID int64, followUserID int64) error {
	if userID == followUserID {
		return errorx.NewDefaultError("不能关注自己哦")
	}

	isFollow, err := IsFollow(ctx, svcCtx, userID, followUserID)
	if err != nil {
		return err
	}

	if !isFollow {
		return errorx.NewDefaultError("您还没有关注该用户")
	}

	// 事务
	err = query.Q.Transaction(func(tx *query.Query) error {
		// 删除关注
		f := tx.Follow
		_, err = f.Where(f.UserID.Eq(userID), f.FollowID.Eq(followUserID)).Delete()
		if err != nil {
			return err
		}

		// redis删除关注和被关注
		err = svcCtx.RedisClient.SRem(ctx, consts.FollowPrefix+strconv.Itoa(int(userID)), followUserID).Err()
		if err != nil {
			return err
		}
		err = svcCtx.RedisClient.SRem(ctx, consts.FollowerPrefix+strconv.Itoa(int(followUserID)), userID).Err()
		if err != nil {
			return err
		}

		u := tx.User

		// 更新用户关注数
		_, err = u.Where(u.ID.Eq(userID)).UpdateColumn(u.FollowCount, gorm.Expr("follow_count - ?", 1))
		if err != nil {
			return err
		}

		// 更新用户粉丝数
		_, err = u.Where(u.ID.Eq(followUserID)).UpdateColumn(u.FollowerCount, gorm.Expr("follower_count - ?", 1))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "取消关注失败")
	}

	return nil
}

func IsFollow(ctx context.Context, svcCtx *svc.ServiceContext, userID int64, followUserID int64) (bool, error) {
	if userID == followUserID {
		return true, nil
	}

	return svcCtx.RedisClient.SIsMember(ctx, consts.FollowPrefix+strconv.Itoa(int(userID)), followUserID).Result()

	//q := query.Follow
	//count, err := q.Where(q.UserID.Eq(userID), q.FollowID.Eq(followUserID)).Count()
	//if err != nil {
	//	return false, err
	//}

	//return count != 0, nil
}

func GetUserByName(ctx context.Context, svcCtx *svc.ServiceContext, username string) (*model.User, error) {
	q := query.User
	user, err := q.Where(q.Name.Eq(username)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewDefaultError("用户不存在")
		}
		return nil, err
	}

	return user, nil
}

func GetUserByID(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*model.User, error) {
	q := query.User
	user, err := q.Where(q.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewDefaultError("用户不存在")
		}
		return nil, err
	}

	return user, nil
}

func CreateUser(ctx context.Context, svcCtx *svc.ServiceContext, username, password string) (int64, error) {
	q := query.User
	user := &model.User{
		Name:            username,
		Password:        util.EncryptPassword(password),
		Avatar:          consts.DefaultAvatarURL,
		BackgroundImage: consts.DefaultBackgroundImageURL,
	}

	err := q.Create(user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
