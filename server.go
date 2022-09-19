//go:build server

package main

import (
	"net/http"

	"quekr/server/service"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	svc, err := service.NewService()
	 if err != nil {
		panic(err)
	}

	server := gin.Default()
	server.Static("/static", "./static")
	server.LoadHTMLGlob("templates/*")
	server.GET("/", func(response *gin.Context) {
		response.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home Page",
		},
		)
	})
	server.POST("/",func (response *gin.Context)  {
		originalURL := response.PostForm("originalURL");
		ipAddress := ReadUserIP(response.Request)
		info, err := svc.CreateMapping(originalURL, ipAddress)
		if err != nil {
			panic(err)
		}
		responseHTML :=
		`<html>
			<form action="/urls" method="POST">
				<input type="text" name="innerUrl" value="https://que.kr/` +info.ShortKey+ `"/>
				<input type="text" name="adminUrl" value="https://que.kr/`+info.ShortKey+ "/" +info.SecretToken+`"/>
			</form>
			<script>document.forms[0].submit();</script>
		</html>`
		response.Data(http.StatusOK, "text/html; charset=utf-8", []byte(responseHTML))
	})
	server.POST("/urls", func(response *gin.Context) {
		innerUrl := response.PostForm("innerUrl");
		adminUrl := response.PostForm("adminUrl");
		response.HTML(http.StatusOK, "result.html", gin.H{
			"innerUrl": innerUrl,
			"adminUrl": adminUrl,
		},
		)
	})
	server.GET("/statements", func(response *gin.Context) {
		response.HTML(http.StatusOK, "statement.html", gin.H{
			"accessMin": `5분전/4분전/3분전/2분전/1분전,23/24/25/26/27`,
			"accessDay": `23일/24일/25일/26일/27일,23/24/25/26/27`,
			"accessDevice": "20/50/30", //Moblie, Web, Etc
			"referer": `naver.com/32,daum.net/21,google.com/12,tistory.com/3,kakao.com/6,linho.kr/1`,
		},
		)
	})
	server.GET("/:shortkey", func(response *gin.Context) {
		shortkey := response.Param("shortkey")
		item, err := svc.QueryMapping(shortkey)
		if err != nil {
			panic(err)
		}
    response.Redirect(http.StatusTemporaryRedirect, "http://"+item.OriginalUrl)
	})

	server.PUT("/urls", func(response *gin.Context) {
		response.String(200, "Hi")
	})
	server.DELETE("/urls", func(response *gin.Context) {	
		response.String(200, "Hi")
	})
	return server
}

func main() {
	server := setupRouter()
	server.Run(":7979")
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
			IPAddress = r.RemoteAddr
	}
	return IPAddress
}