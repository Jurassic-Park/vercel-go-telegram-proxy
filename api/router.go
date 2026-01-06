package api

import (
	"fmt"
	"net/http"
	"strings"

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

const apiUrl = "https://studio-api.prod.suno.com"
const authUrl = "https://auth.suno.com"

func Suno(router *gin.Engine) {

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
		req, err := http.NewRequestWithContext(context, context.Request.Method, url, context.Request.Body)
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

		reader := resp.Body

		// body, err := io.ReadAll(resp.Body)
		// fmt.Println("请求结果：", string(body), resp.StatusCode)

		contentLength := resp.ContentLength
		contentType := resp.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			// "Content-Disposition": `attachment; filename="gopher.png"`,
		}
		for key, values := range resp.Header {
			if key == "Content-Encoding" || key == "Content-Length" ||
				key == "Transfer-Encoding" || key == "Connection" {
				continue
			}
			for _, value := range values {
				extraHeaders[key] = value
			}
		}

		context.DataFromReader(resp.StatusCode, contentLength, contentType, reader, extraHeaders)
	})
}
