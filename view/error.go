package view

import (
	"github.com/gin-gonic/gin"
)

func ResponseErrorHtml(response *gin.Context,statusCode int, errMessage string){
	response.HTML(statusCode,"error.html",gin.H{"error":errMessage})
}
