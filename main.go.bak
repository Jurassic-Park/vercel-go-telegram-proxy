package main

import (
	"fmt"
	"vercel-go-telegram-proxy/api"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	api.Suno(router)

	fmt.Println("网关已启动，监听 :8080")
	router.Run(":8080")
}
