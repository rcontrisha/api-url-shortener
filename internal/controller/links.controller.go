package controller

import (
	"log"
	"strconv"

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

	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	ctx.Redirect(301, originalUrl)
}

func (l *LinkController) DeleteShortenedLink(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, dto.Response{
			Success: false,
			Message: "invalid link id",
		})
		return
	}

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

	err = l.service.DeleteShortenedLink(ctx.Request.Context(), id, userId)
	if err != nil {
		ctx.JSON(404, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "short link deleted successfully",
	})
}

func (l *LinkController) GetUserLinks(ctx *gin.Context) {
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

	search := ctx.DefaultQuery("search", "")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	res, err := l.service.GetUserLinks(ctx.Request.Context(), userId, search, page, limit)
	if err != nil {
		ctx.JSON(500, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.Response{
		Success: true,
		Message: "successfully retrieve user links",
		Results: res,
	})
}
