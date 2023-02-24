package middleware

import (
	"context"
	"douyin-zero/common/errorx"
	"douyin-zero/consts"
	"douyin-zero/util"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"strconv"
)

type AuthMiddleware struct {
	RedisClient *redis.Client
}

func NewAuthMiddleware(redisClient *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		RedisClient: redisClient,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAllowed := util.MatchURI(r.RequestURI,
			"/douyin/feed/**",
			"/douyin/publish/list/**",
			"/douyin/comment/list/**",
		)
		// 处理允许通过的路由
		if isAllowed {
			r = m.HandleAllow(r)
			next(w, r)
			return
		}

		token := r.FormValue("token")
		if token == "" {
			token = r.PostFormValue("token")
			if token == "" {
				httpx.ErrorCtx(r.Context(), w, errorx.NewDefaultError("token is empty"))
				return
			}
		}

		// token不为空，从redis中获取用户信息
		userIDString := m.RedisClient.Get(context.TODO(), consts.TokenPrefix+token).Val()
		if userIDString == "" {
			httpx.ErrorCtx(r.Context(), w, errorx.NewDefaultError("token错误"))
			return
		}
		// 存入上下文
		userID, _ := strconv.Atoi(userIDString)
		ctx := context.WithValue(r.Context(), consts.UserIDKey, int64(userID))

		r = r.WithContext(ctx)

		next(w, r)
	}
}

func (m *AuthMiddleware) HandleAllow(r *http.Request) *http.Request {
	token := r.FormValue("token")
	// token为空，设置上下文id为0
	if token == "" {
		token = r.PostFormValue("token")
		if token == "" {
			var userID int64 = 0
			ctx := context.WithValue(r.Context(), consts.UserIDKey, int64(userID))
			r = r.WithContext(ctx)
			return r
		}
	}

	// token不为空，从redis中获取用户信息
	userIDString := m.RedisClient.Get(context.TODO(), consts.TokenPrefix+token).Val()
	if userIDString == "" {
		userIDString = "0"
	}
	// 存入上下文
	userID, _ := strconv.Atoi(userIDString)
	ctx := context.WithValue(r.Context(), consts.UserIDKey, int64(userID))
	r = r.WithContext(ctx)
	return r
}
