package handler

import (
	"douyin-zero/common/response"
	"douyin-zero/internal/logic"
	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func VideoActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoActionRequest
		logx.Infof("r.Header : %v", r.Header)
		logx.Infof("r.Method : %v", r.Method)
		logx.Infof("r.Body : %v", r.PostForm)
		if err := httpx.Parse(r, &req); err != nil {
			response.Error(w, errors.New("参数错误"))
			return
		}
		if err := validator.New().StructCtx(r.Context(), req); err != nil {
			response.Error(w, errors.New("参数错误"))
			return
		}
		logx.Infof("req : %v", req)

		_, fileHeader, err := r.FormFile("data")
		if err != nil {
			logx.Infof("文件上传失败,err : %v", err)
			response.Error(w, errors.New("文件上传失败"))
			return
		}

		l := logic.NewVideoActionLogic(r.Context(), svcCtx)
		resp, err := l.VideoAction(&req, fileHeader)

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}

}
