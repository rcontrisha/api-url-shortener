package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rcontrisha/api-url-shortener/internal/dto"
	"github.com/rcontrisha/api-url-shortener/internal/service"
	"github.com/rcontrisha/api-url-shortener/pkg"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

func (a *AuthController) LoginController(ctx *gin.Context) {
	var payload dto.LoginRequest
	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		log.Println("error", err.Error())
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: "Internal Server Error",
		})
		return
	}

	token, err := a.service.LoginService(ctx.Request.Context(), payload)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "Login Failed. Invalid email or password.",
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "Login Success.",
		Results: gin.H{
			"token": token,
		},
	})
}

func (a *AuthController) RegisterController(ctx *gin.Context) {
	var payload dto.LoginRequest

	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		log.Println("error", err.Error())
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: "Internal Server Error",
		})
		return
	}

	err := a.service.RegisterService(ctx.Request.Context(), payload)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "Register failed.",
		})
		return
	}

	ctx.JSON(201, dto.Response{
		Success: true,
		Message: "Register Success.",
	})
}

func (a *AuthController) LogoutController(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "Token Data Not Found in Context.",
		})
	}

	claims := token.(pkg.JWTClaims)
	jti := claims.ID
	expiresAt := claims.ExpiresAt
	if err := a.service.Logout(ctx, jti, expiresAt.Time); err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "Successfully logged out",
	})
}