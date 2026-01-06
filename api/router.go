package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

// suno
func init() {
	router = gin.Default()
	Suno(router)
}

func Listen(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
