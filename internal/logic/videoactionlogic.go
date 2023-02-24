package logic

import (
	"context"
	"douyin-zero/common/errorx"
	"douyin-zero/consts"
	"douyin-zero/internal/dal"
	"douyin-zero/internal/dal/model"
	"douyin-zero/util"
	"github.com/zeromicro/go-zero/core/threading"
	"mime/multipart"
	"path"
	"time"

	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoActionLogic {
	return &VideoActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoActionLogic) VideoAction(req *types.VideoActionRequest, file *multipart.FileHeader) (resp *types.VideoActionResponse, err error) {
	l.Infof("VideoFavoriteAction req:%v", req)
	if file.Size > consts.MaxVideoSize {
		return nil, errorx.NewDefaultError("视频过大，无法上传")
	}

	ext := path.Ext(file.Filename)
	allowExt := map[string]bool{
		".mp4": true,
		".jpg": true,
		".png": true,
	}
	if _, ok := allowExt[ext]; !ok {
		return nil, errorx.NewDefaultError("视频格式错误")
	}

	playName := util.UUID() + ext
	coverName := util.UUID() + ".jpg"
	playURL := consts.PlayURLPrefix + playName
	CoverURL := consts.CoverURLPrefix + coverName

	routineGroup := threading.NewRoutineGroup()

	// 保存视频并截取封面
	routineGroup.Run(func() {
		saveStart := time.Now().UnixMilli()
		err = util.SaveUploadedFile(file, consts.PlayFilePath+"/"+playName)
		if err != nil {
			return
		}
		saveEnd := time.Now().UnixMilli()
		l.Infof("保存视频用时为: %v", saveEnd-saveStart)

		coverStart := time.Now().UnixMilli()
		err = l.CaptureCover(playName, coverName)
		coverEnd := time.Now().UnixMilli()
		l.Infof("截取视频用时为: %v", coverEnd-coverStart)
	})

	userID := l.ctx.Value(consts.UserIDKey).(int64)

	video := &model.Video{
		AuthorID:   userID,
		Title:      req.Title,
		PlayURL:    playURL,
		CoverURL:   CoverURL,
		CreateTime: time.Now().UnixMilli(),
	}

	// 存到数据库和redis
	routineGroup.Run(func() {
		err = dal.CreateVideo(l.ctx, l.svcCtx, video)
	})

	routineGroup.Wait()
	if err != nil {
		return
	}

	resp = &types.VideoActionResponse{
		BaseResponse: types.BaseResponse{
			Message: "上传成功",
		},
	}
	return
}

func (l *VideoActionLogic) CaptureCover(playName, coverName string) error {
	return util.CmdWithDirNoOut(consts.PlayFilePath,
		"ffmpeg",
		"-i", consts.PlayFilePath+"/"+playName,
		"-y",
		"-f", "image2",
		"-frames", "1",
		"-ss", "1",
		consts.CoverFilePath+"/"+coverName,
	)
}
