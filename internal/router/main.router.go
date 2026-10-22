package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rcontrisha/api-url-shortener/internal/middleware"
	"github.com/redis/go-redis/v9"
)

func MainRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	router.Use(middleware.Cors)

	AuthRouter(router, db, rdb)
	LinkRouter(router, db, rdb)
}
