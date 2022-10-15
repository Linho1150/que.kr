//go:build server

package main

import (
	"fmt"
	"net/http"
	"quekr/server/service"
	"quekr/server/view"

	"github.com/gin-gonic/gin"
)
type ChangeUrlRequestBody struct {
	Url string
}

func setupRouter() *gin.Engine {
	svc, err := service.NewService()
	 if err != nil {
		panic(err)
	}	

	server := gin.Default()
	server.Static("/static", "./static")
	server.LoadHTMLGlob("templates/*")


	server.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html",nil)
	})
	server.POST("/",func (c *gin.Context)  {
		originalURL := c.PostForm("originalURL");
		responseHTML, err := view.CreatUrl(c,svc,originalURL);
		if(err != nil){
			if(err.Error()=="ValidationURL"){
				view.ResponseErrorHtml(c,http.StatusBadRequest,"It's not Url (ValidationURL)")
				return
			}
			if(err.Error()=="CraetMapping"){
				view.ResponseErrorHtml(c,http.StatusServiceUnavailable,"Server Error (CraetMapping)")
				return
			}
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(responseHTML))
	})
	server.POST("/urls", func(c *gin.Context) {
		innerUrl := c.PostForm("innerUrl");
		adminUrl := c.PostForm("adminUrl");
		c.HTML(http.StatusOK, "result.html", gin.H{
			"innerUrl": innerUrl,
			"adminUrl": adminUrl,
		},
		)
	})
	server.GET("/:shortkey", func(c *gin.Context) {
		shortkey := c.Param("shortkey")
		item, err := view.RedirectUrl(c,svc,shortkey)
		if(err != nil){
			if(err.Error()=="QueryMapping"){
				view.ResponseErrorHtml(c,http.StatusNotFound,"Not Found (QueryMapping)")
				return
			}
			if(err.Error()=="CraetMapping"){
				view.ResponseErrorHtml(c,http.StatusServiceUnavailable,"Server Error (TouchStatistics)")
				return
			}
		}
    c.Redirect(http.StatusTemporaryRedirect, item.OriginalUrl)
	})

	server.GET("/:shortkey/:secrettoken", func(c *gin.Context) {
		shortkey := c.Param("shortkey")
		secrettoken := c.Param("secrettoken")

		info, err := svc.QueryMapping(shortkey)
		if err != nil {
			view.ResponseErrorHtml(c,http.StatusBadRequest,"Bad request(QueryMapping)")
			return
		}

		if info.SecretToken != secrettoken{
			view.ResponseErrorHtml(c,http.StatusUnauthorized,"Secret token isn't correct(QueryMapping)")
			return
		}

		dataReferer,dataDeviceType,dataTimerPerDate,dataTimePerMinute,err := view.GetStaticData(c,svc,shortkey)
		if(err != nil){
			view.ResponseErrorHtml(c,http.StatusServiceUnavailable,"Server Error("+err.Error()+")")
			return
		}

		accessDeviceString, refererString, accessDayString, accessMinString, err := view.RefineData(
			svc,dataDeviceType,dataReferer,dataTimerPerDate,dataTimePerMinute)
		if(err != nil){
			view.ResponseErrorHtml(c,http.StatusServiceUnavailable,"Server Error("+err.Error()+")")
			return
		}

		fmt.Println(
			string(accessMinString)+
			string(accessDayString)+
			string(accessDeviceString)+
			string(refererString)+
			info.OriginalUrl)

		c.HTML(http.StatusOK, "statement.html", gin.H{
			"accessMin": string(accessMinString),
			"accessDay": string(accessDayString),
			"accessDevice": string(accessDeviceString),
			"referer": string(refererString),
			"origianlURL": info.OriginalUrl,
		},
		)
	})
	server.PUT("/:shortkey/:secrettoken", func(c *gin.Context) {
		shortkey := c.Param("shortkey")
		secrettoken := c.Param("secrettoken")
		var requestBody ChangeUrlRequestBody
		if err := c.BindJSON(&requestBody); err != nil {
			view.ResponseErrorHtml(c,http.StatusBadRequest,"Bad Request(BindJSON)")
			return
		}
		if(view.ValidationURL(requestBody.Url)){
			view.ResponseErrorHtml(c,http.StatusBadRequest,"It's not URL(ValidationURL)")
			return
		}
		err = svc.UpdateMapping(shortkey, secrettoken, requestBody.Url)
		if err != nil {
			view.ResponseErrorHtml(c,http.StatusUnauthorized,"Unauthorized(UpdateMapping)")
			return
		}
		c.Status(http.StatusOK);
	})
	server.DELETE("/:shortkey/:secrettoken", func(c *gin.Context) {	
		shortkey := c.Param("shortkey")
		secrettoken := c.Param("secrettoken")
		err = svc.RemoveMapping(shortkey, secrettoken)
		if err != nil {
			view.ResponseErrorHtml(c,http.StatusBadRequest,"Bad request(RemoveMapping)")
			return
		}
		c.Status(http.StatusOK);
	})
	return server
}

func main() {
	server := setupRouter()
	server.Run(":7979")
}