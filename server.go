package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	server := gin.Default()
	server.Static("/static", "./static")
	server.LoadHTMLGlob("templates/*")
	server.GET("/", func(response *gin.Context) {
		response.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home Page",
		},
		)
	})
	server.GET("/urls", func(response *gin.Context) {
		response.HTML(http.StatusOK, "result.html", gin.H{
			"innerUrl": "www.naver.com",
			"outerUrl": "www.google.com",
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
