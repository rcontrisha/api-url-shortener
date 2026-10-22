package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-url-shortener/internal/controller"
	"github.com/rcontrisha/api-url-shortener/internal/middleware"
	"github.com/rcontrisha/api-url-shortener/internal/repository"
	"github.com/rcontrisha/api-url-shortener/internal/service"
	"github.com/redis/go-redis/v9"
)

func LinkRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	linkRouter := r.Group("link")

	repo := repository.NewLinkRepo(db)
	service := service.NewLinkService(repo)
	controller := controller.NewLinkController(service)

	linkRouter.GET("links", middleware.AuthMiddleware(rdb), controller.GetUserLinks)
	linkRouter.POST("create", middleware.AuthMiddleware(rdb), controller.CreateShortLink)
	linkRouter.DELETE("/:id/delete", middleware.AuthMiddleware(rdb), controller.DeleteShortenedLink)

	r.GET("/:slug", controller.GetOriginalUrl)
}
