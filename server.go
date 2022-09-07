package main

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
	server := gin.Default()
	server.LoadHTMLGlob("templates/*")
	server.GET("/", func(reposne *gin.Context) {
		reposne.String(200, "Hello World")
	})
	return server
}

func main() {
	server := setupRouter()
	server.Run(":7979")
}
