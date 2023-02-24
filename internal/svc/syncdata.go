package svc

import (
	"context"
	"douyin-zero/consts"
	"douyin-zero/internal/dal/model"
	"douyin-zero/internal/dal/query"
	"fmt"
	"github.com/robfig/cron"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"strconv"
	"strings"
	"time"
)

func SyncTask(svcCtx *ServiceContext) {
	ctx := context.Background()

	routineGroup := threading.NewRoutineGroup()

	spec := fmt.Sprintf("*/%d * * * * *", consts.SyncTime)
	c := cron.New()
	c.AddFunc(spec, func() {
		start := time.Now().UnixMilli()
		logx.Infof("sync task start %v", time.Now().Format("2006-01-02 15:04:05"))

		routineGroup.Run(func() { SyncVideoFavoriteRecord(ctx, svcCtx) })
		routineGroup.Run(func() { SyncUserFavoriteCount(ctx, svcCtx) })
		routineGroup.Run(func() { SyncUserFavoritedCount(ctx, svcCtx) })
		routineGroup.Run(func() { SyncVideoFavoriteCount(ctx, svcCtx) })

		routineGroup.Wait()

		end := time.Now().UnixMilli()
		logx.Infof("sync task end %v", time.Now().Format("2006-01-02 15:04:05"))
		logx.Infof("sync task cost %vms", end-start)
	})

	go c.Start()
	defer c.Stop()
}

func SyncVideoFavoriteRecord(ctx context.Context, svcCtx *ServiceContext) {
	resVal := svcCtx.RedisClient.HGetAll(ctx, consts.VideoFavorite).Val()
	for k, v := range resVal {
		// 对k分离出视频id和用户id
		videoAndUserID := strings.Split(k, "-")
		videoIDStr := videoAndUserID[0]
		userIDStr := videoAndUserID[1]

		videoID, _ := strconv.Atoi(videoIDStr)
		userID, _ := strconv.Atoi(userIDStr)

		// 判断用户是否点赞
		if v == "1" {
			// 点赞，插入点赞记录
			_ = CreateFavorite(ctx, svcCtx, int64(videoID), int64(userID))
		} else {
			// 取消点赞，删除点赞记录
			_ = DeleteFavorite(ctx, svcCtx, int64(videoID), int64(userID))
		}
		// 删除redis中的点赞记录
		svcCtx.RedisClient.HDel(ctx, consts.VideoFavorite, k)
	}
}

func SyncVideoFavoriteCount(ctx context.Context, svcCtx *ServiceContext) {
	resVal := svcCtx.RedisClient.HGetAll(ctx, consts.VideoFavoriteCount).Val()
	for k, v := range resVal {
		// 对k分离出视频id和用户id
		videoID, _ := strconv.Atoi(k)
		favoriteCount, _ := strconv.Atoi(v)
		// 更新视频点赞数
		_ = UpdateVideoFavoriteCount(ctx, svcCtx, int64(videoID), int64(favoriteCount))
	}
}

func SyncUserFavoriteCount(ctx context.Context, svcCtx *ServiceContext) {
	resVal := svcCtx.RedisClient.HGetAll(ctx, consts.UserFavoriteCount).Val()
	for k, v := range resVal {
		// 对k分离出视频id和用户id
		userID, _ := strconv.Atoi(k)
		favoriteCount, _ := strconv.Atoi(v)
		// 更新用户点赞数
		_ = UpdateUserFavoriteCount(ctx, svcCtx, int64(userID), int64(favoriteCount))
	}
}

func SyncUserFavoritedCount(ctx context.Context, svcCtx *ServiceContext) {
	resVal := svcCtx.RedisClient.HGetAll(ctx, consts.UserFavoritedCount).Val()
	for k, v := range resVal {
		// 对k分离出视频id和用户id
		userID, _ := strconv.Atoi(k)
		favoritedCount, _ := strconv.Atoi(v)
		// 更新作者获赞数
		_ = UpdateUserFavoritedCount(ctx, svcCtx, int64(userID), int64(favoritedCount))
	}
}

func DeleteFavorite(ctx context.Context, svcCtx *ServiceContext, videoID int64, userID int64) error {
	f := query.Favorite
	f.Where(f.UserID.Eq(userID), f.VideoID.Eq(videoID)).Delete()
	return nil
}

func CreateFavorite(ctx context.Context, svcCtx *ServiceContext, videoID int64, userID int64) error {
	f := &model.Favorite{
		VideoID:    videoID,
		UserID:     userID,
		CreateTime: time.Now().UnixMilli(),
	}

	return query.Favorite.Create(f)
}

func UpdateUserFavoriteCount(ctx context.Context, svcCtx *ServiceContext, userID int64, favoriteCount int64) error {
	u := query.User
	query.User.Where(query.User.ID.Eq(userID)).Update(u.FavoriteCount, favoriteCount)
	return nil
}

func UpdateUserFavoritedCount(ctx context.Context, svcCtx *ServiceContext, userID int64, favoritedCount int64) error {
	u := query.User
	query.User.Where(query.User.ID.Eq(userID)).Update(u.TotalFavorited, favoritedCount)
	return nil
}

func UpdateVideoFavoriteCount(ctx context.Context, svcCtx *ServiceContext, videoID int64, count int64) error {
	_, err := GetVideoByID(ctx, svcCtx, videoID)
	if err != nil {
		return err
	}

	v := query.Video
	query.Video.Where(query.Video.ID.Eq(videoID)).Update(v.FavoriteCount, count)
	return nil
}

func GetVideoByID(ctx context.Context, svcCtx *ServiceContext, videoID int64) (*model.Video, error) {
	//// 先查缓存
	//video := &model.Video{}
	//err := svcCtx.RedisClient.Get(ctx, consts.VideoInfo+strconv.Itoa(int(videoID))).Scan(video).Err()
	//if err == nil {
	//	return video, nil
	//}
	//
	//// 缓存中没有查数据库
	//v := query.Video
	//video, err = v.Where(v.ID.Eq(videoID)).First()
	//if err != nil {
	//	return nil, err
	//}
	//
	//// 缓存到redis
	//err = svcCtx.RedisClient.Set(ctx, consts.VideoInfo+strconv.Itoa(int(videoID)), video, 0).Err()
	//if err != nil {
	//	logx.Error(err)
	//}

	// 直接查数据库
	v := query.Video
	return v.Where(v.ID.Eq(videoID)).First()
}
