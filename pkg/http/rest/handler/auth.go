package handler

import (
	"github.com/HongXiangZuniga/login-go/pkg/authentication"
	"github.com/HongXiangZuniga/login-go/pkg/authorize"
	http "github.com/HongXiangZuniga/login-go/pkg/http/rest"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authorizeService authorize.Service
	authService      authentication.Service
	logger           *zap.Logger
}

func NewAuthHanler(authorizeService authorize.Service, authenticationService authentication.Service, logger *zap.Logger) GinHandler {
	return &AuthHandler{
		authorizeService: authorizeService,
		authService:      authenticationService,
		logger:           logger}
}

func (impl *AuthHandler) RegisterHandler(router *gin.RouterGroup) {
	router.POST("", impl.GetAuth)
}

func (impl *AuthHandler) GetAuth(ctx *gin.Context) {
	var loginData http.LoginRequest
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result, err := impl.authService.Authorization(loginData.User, loginData.Password)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if *result {
		hash, err := impl.authorizeService.SetHash(loginData.User)
		if err != nil {
			impl.logger.Error(err.Error())
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.SetCookie("session", *hash, 300, "/", "localhost", false, true)
		ctx.JSON(200, nil)
	} else {
		ctx.JSON(401, gin.H{"error": "User Unauthorized"})
	}

}
