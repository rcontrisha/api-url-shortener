package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rcontrisha/api-url-shortener/internal/dto"
	"github.com/rcontrisha/api-url-shortener/internal/service"
	"github.com/rcontrisha/api-url-shortener/pkg"
)

type LinkController struct {
	service *service.LinkService
}

func NewLinkController(service *service.LinkService) *LinkController {
	return &LinkController{
		service: service,
	}
}

func (l *LinkController) CreateShortLink(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist {
		ctx.JSON(401, dto.Response{
			Success: false,
			Message: "Token Data Not Found in Context.",
		})
		return
	}
	claims := token.(pkg.JWTClaims)
	userId := claims.Id

	var payload dto.ShortenLink

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := l.service.ShortenLinkService(ctx, userId, payload); err != nil {
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, dto.Response{
		Success: true,
		Message: "successfully shorten link.",
	})
}

func (l *LinkController) GetOriginalUrl(ctx *gin.Context) {
	slug := ctx.Param("slug")

	if slug == "" {
		ctx.JSON(400, dto.Response{
			Success: false,
			Message: "slug is required",
		})
		return
	}

	originalUrl, err := l.service.GetOriginalUrl(ctx, slug)
	if err != nil {
		ctx.JSON(404, dto.Response{
			Success: false,
			Message: "URL Not Found",
		})
		return
	}
	log.Println("[Link Controller]-Original URL: ", originalUrl)

	ctx.Redirect(301, originalUrl)
}
