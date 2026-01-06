package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/proxy"
)

var router *gin.Engine

var apiUrl = "https://studio-api.prod.suno.com"
var authUrl = "https://auth.suno.com"

// suno
func init() {
	router = gin.Default()
	// 初始化代理中间件，默认管理接口挂载在 /endpoints 下
	p := proxy.New(router)

	// 配置路由规则
	setupProxyRoutes(p)

	// 你的 Gin 依然可以处理普通路由
	// router.GET("/ping", func(c *gin.Context) {
	//     c.JSON(200, gin.H{"message": "我是网关本体"})
	// })
}

func setupProxyRoutes(p *proxy.Proxy) {
	// 规则1：将 /proxy/ 开头的请求分发给集群 A
	// 默认使用轮询策略，健康检查间隔5秒
	err := p.Pass("/suno-api/", []string{
		apiUrl,
	})
	if err != nil {
		panic(err)
	}

	// 规则2：将 /personal/ 开头的请求分发给集群 B
	err = p.Pass("/auth-api/", []string{
		authUrl,
	})
	if err != nil {
		panic(err)
	}
}

func Listen(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
