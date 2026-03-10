package api

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	Proxy(router)
}

func Listen(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}

func Proxy(router *gin.Engine) {

	router.Any("/*path", func(context *gin.Context) {
		uri := context.Param("path")
		var reqUrl string

		// 获取query参数
		queryParams := context.Request.URL.RawQuery
		if queryParams != "" {
			uri = uri + "?" + queryParams
		}

		if after, ok := strings.CutPrefix(uri, "/suno.com-api"); ok {
			reqUrl = "https://studio-api.prod.suno.com"
			uri = after
		} else if after0, ok0 := strings.CutPrefix(uri, "/suno.com-auth-api"); ok0 {
			reqUrl = "https://auth.suno.com"
			uri = after0
		} else if after1, ok1 := strings.CutPrefix(uri, "/aiplatform.googleapis.com"); ok1 {
			// 默认 us-central1
			reqUrl = "https://us-central1-aiplatform.googleapis.com"
			uri = after1
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
		cookies := context.Request.Cookies()
		for _, cookie := range cookies {
			fmt.Printf("添加cookie: %s=%s\n", cookie.Name, cookie.Value)
			req.AddCookie(cookie)
		}
		// 获取除cookie外的header
		for key, values := range context.Request.Header {
			if strings.Contains(key, "Vercel") || key == "Cookie" || key == "Accept-Encoding" ||
				key == "Accept" || key == "Connection" || key == "Content-Length" || key == "Host" ||
				key == "Origin" || key == "Referer" || key == "User-Agent" ||
				strings.Contains(key, "Forwarded") || key == "X-Real-Ip" {
				continue
			}
			for _, value := range values {
				fmt.Printf("添加header: %s=%s\n", key, value)
				req.Header.Add(key, value)
			}
		}
		fmt.Println("请求地址:", url, context.Request.Method, "body:", context.Request.Body)

		client := &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   1000 * time.Millisecond,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ResponseHeaderTimeout: 150 * time.Second,
			},
		}

		resp, err := client.Do(req)
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
