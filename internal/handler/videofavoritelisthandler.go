package handler

import (
	"douyin-zero/common/response"
	"douyin-zero/internal/logic"
	"douyin-zero/internal/svc"
	"douyin-zero/internal/types"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func VideoFavoriteListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoFavoriteListRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Error(w, errors.New("参数错误"))
			return
		}
		if err := validator.New().StructCtx(r.Context(), req); err != nil {
			response.Error(w, errors.New("参数错误"))
			return
		}

		l := logic.NewVideoFavoriteListLogic(r.Context(), svcCtx)
		resp, err := l.VideoFavoriteList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
