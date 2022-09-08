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
			"accessMin":    "Home Page",
			"accessDay":    "Home Page",
			"accessDevice": "Home Page",
			"referer":      "Home Page",
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
