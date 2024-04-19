package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/monaco-io/request"
	"github.com/qm012/dun"
	"go.uber.org/zap"
	"net/http"
	"vland.live/app/global"
)

const (
	BearerSchema         = "Bearer "
	ContextMemberIdKey   = "memberId"
	ContextMemberTypeKey = "memberType" // 用户类型
)

var (
	statusCodeNoPermission        = dun.NewStatusCode(http.StatusUnauthorized, "权限不足，无法进一步操作")
	statusCodeNeedAdminPermission = dun.NewStatusCode(http.StatusUnauthorized, "需要管理员权限才可以操作")
	statusCodeChooseSigninType    = dun.NewStatusCode(http.StatusUnauthorized, "需要选择一种登录方式")
	statusCodeNeedReSignin        = dun.NewStatusCode(http.StatusUnauthorized, "账号已失效，请重新登录")
)

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.StandardClaims
}

type BaseClaims struct {
	ID          int64
	Username    string
	NickName    string
	AuthorityId string
}

func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		if authorization == "" {
			ctx.Abort()
			dun.Failed401(ctx, statusCodeNoPermission)
			return
		}
		tokenString := authorization[len(BearerSchema):]
		// 验证是否是管理后台地址
		_, ok := adminPath[ctx.Request.Method+ctx.FullPath()]
		// 2. 从后台管理中去验证
		if ok {
			client := &request.Client{
				//URL:    global.Config.App.AdminAuthVerifyUrl,
				Method: http.MethodPost,
				Header: map[string]string{
					"x-token": tokenString,
					"app":     "verseland-go",
				},
			}
			type Response struct {
				Code   int          `json:"code"`
				Msg    string       `json:"msg"`
				Claims CustomClaims `json:"data"`
			}

			r := new(Response)
			resp := client.Send().Scan(r)
			if !resp.OK() || r.Code != 200 {
				ctx.Abort()
				global.Logger.Error("验证后台管理的管理员权限失败", zap.Error(resp.Error()), zap.Any("code", r.Code))
				dun.Failed401(ctx, dun.StatusCodeInternalService)
				return
			}
			ctx.Next()
			global.Logger.Info("管理员操作数据获取到的 Claims", zap.Any("Claims", r.Claims))
			ctx.Set(ContextMemberIdKey, r.Claims.ID)
			return
		}
	}
}
