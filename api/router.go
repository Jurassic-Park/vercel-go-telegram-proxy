package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

var apiUrl = "https://studio-api.prod.suno.com"
var authUrl = "https://auth.suno.com"

// suno
func init() {
	router = gin.Default()
	router.Any("/*path", func(context *gin.Context) {
		uri := context.Param("path")
		reqUrl := apiUrl
		fmt.Println("请求的链接", uri)
		if after, ok := strings.CutPrefix(uri, "/suno-api"); ok {
			uri = after
		} else if after0, ok0 := strings.CutPrefix(uri, "/auth-api"); ok0 {
			uri = after0
			reqUrl = authUrl
		} else {
			context.String(http.StatusNotFound, "404 Not found")
			return
		}
		url := reqUrl + uri
		req, err := http.NewRequestWithContext(context, context.Request.Method, url, context.Request.Body)
		if err != nil {
			fmt.Println(err)
			context.String(http.StatusBadRequest, err.Error())
			return
		}
		req.Header = context.Request.Header
		req.PostForm = context.Request.PostForm
		req.Form = context.Request.Form
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			context.String(http.StatusBadRequest, err.Error())
			return
		}
		context.DataFromReader(resp.StatusCode, resp.ContentLength, "application/json", resp.Body, nil)
	})
}

func Listen(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
