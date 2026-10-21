package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rcontrisha/api-url-shortener/internal/config"
	"github.com/rcontrisha/api-url-shortener/internal/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file")
	}

	pdb := config.NewPsql(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	pool, err := pdb.Connect()
	if err != nil {
		log.Println("Cannot Connect to DB\nReason: ", err.Error())
		return
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Println("Database is not ready\nReason: ", err.Error())
		return
	}
	log.Println("Database is ready")

	rc := config.NewRedisClient(os.Getenv("REDIS_USER"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	rdb := rc.Connect()
	defer rdb.Close()
	status, err := rdb.Ping(context.Background()).Result()

	if err != nil {
		fmt.Println("[failed to connect redis]:", err)
	}
	fmt.Println(status)

	fmt.Println("connected to redis")

	r := gin.Default()
	router.MainRouter(r, pool, rdb)
	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}