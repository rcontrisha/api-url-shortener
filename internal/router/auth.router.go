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

func AuthRouter(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := r.Group("auth")

	repo := repository.NewAuthRepo(db, rdb)
	service := service.NewAuthService(repo)
	controller := controller.NewAuthController(service)

	authRouter.POST("login", controller.LoginController)
	authRouter.POST("register", controller.RegisterController)

	authRouter.POST("logout", middleware.AuthMiddleware(rdb), controller.LogoutController)
}
