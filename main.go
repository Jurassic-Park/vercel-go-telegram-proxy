package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var apiUrl = "https://studio-api.prod.suno.com"

var authUrl = "https://auth.suno.com"

// var authUrl = "https://dev.factory.51zhiyong.com"

func main() {
	router := gin.Default()
	router.Any("/*path", func(context *gin.Context) {
		uri := context.Param("path")
		reqUrl := apiUrl

		// 获取query参数
		queryParams := context.Request.URL.RawQuery
		if queryParams != "" {
			uri = uri + "?" + queryParams
		}

		if after, ok := strings.CutPrefix(uri, "/suno.com-api"); ok {
			uri = after
		} else if after0, ok0 := strings.CutPrefix(uri, "/suno.com-auth-api"); ok0 {
			uri = after0
			reqUrl = authUrl
		} else {
			context.String(http.StatusNotFound, "404 Not found")
			return
		}
		url := reqUrl + uri
		fmt.Println("请求", url, context.Request.Method)
		req, err := http.NewRequest(context.Request.Method, url, context.Request.Body)
		if err != nil {
			fmt.Println(err)
			context.String(http.StatusBadRequest, err.Error())
			return
		}
		cookies := context.Request.Cookies()
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		// 获取除cookie外的header
		for key, values := range context.Request.Header {
			if key == "Cookie" || key == "Accept-Encoding" ||
				key == "Accept" {
				continue
			}
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			context.String(http.StatusBadRequest, err.Error())
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		fmt.Println("请求结果：", string(body), resp.StatusCode)

		// context.DataFromReader(resp.StatusCode, resp.ContentLength, "application/json", resp.Body, resp.Cookies())
		context.JSON(resp.StatusCode, resp.Body)
	})

	fmt.Println("网关已启动，监听 :8080")
	router.Run(":8080")
}
