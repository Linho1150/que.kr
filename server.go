//go:build server

package main

import (
	"encoding/json"
	"net/http"
	"quekr/server/service"
	"strconv"
	"strings"
	"time"

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
	server.GET("/:shortkey", func(response *gin.Context) {
		ipAddress := ReadUserIP(response.Request)
		referer := ReadUserReferer(response.Request)
		shortkey := response.Param("shortkey")

		deviceType := service.DeviceTypePC
		mobileDevice := isMobile(response.Request)
		if(mobileDevice){
			deviceType = service.DeviceTypeMobile
		}

		item, err := svc.QueryMapping(shortkey)
		if err != nil {
			panic(err)
		}
	
		if (referer==""){
			referer="Direct access"
		}
		err = svc.TouchStatistics(shortkey, svc.NowLocalTime(), ipAddress, referer, deviceType)
		if err != nil {
			panic(err)
		}
    response.Redirect(http.StatusTemporaryRedirect, "http://"+item.OriginalUrl)
	})

	server.GET("/:shortkey/:secrettoken", func(response *gin.Context) {
		shortkey := response.Param("shortkey")
		secrettoken := response.Param("secrettoken")

		info, err := svc.QueryMapping(shortkey)
		if err != nil {
			panic(err)
		}

		if info.SecretToken != secrettoken{
			panic("Authentication failure");
		}

		dataReferer, err := svc.QueryStatistics(shortkey, service.StatisticLegendTypeReferer, false)
		if err != nil {
			panic(err)
		}
		dataDeviceType, err := svc.QueryStatistics(shortkey, service.StatisticLegendTypeDevicetype, false)
		if err != nil {
			panic(err)
		}

		dataTimerPerDate, err := svc.QueryStatistics(shortkey, service.StatisticLegendTypeTimePerDate, false)
		if err != nil {
			panic(err)
		}

		dataTimePerMinute, err := svc.QueryStatistics(shortkey, service.StatisticLegendTypeTimePerMinute, false)
		if err != nil {
			panic(err)
		}

		accessMin := make([]map[string]interface{}, 0, 0)
		accessDay := make([]map[string]interface{}, 0, 0)
		accessDevice := make([]map[string]interface{}, 0, 0)
		referer := make([]map[string]interface{}, 0, 0)

		for cnt:=0;cnt<3;cnt++{
			if cnt == len(dataDeviceType){
				break
			}
			deviceTypeJson := make(map[string]interface{})
			deviceTypeJson["devicetype"] = dataDeviceType[cnt].Legend.(string)
			deviceTypeJson["devicecount"] = strconv.Itoa(dataDeviceType[cnt].Counter)
			accessDevice = append(accessDevice, deviceTypeJson)
		}
		accessDeviceString,_ := json.Marshal(accessDevice)

		for cnt:=0;cnt<3;cnt++{
			if cnt == len(dataReferer){
				break
			}
			refererJson := make(map[string]interface{})
			refererJson["refererurl"] = dataReferer[cnt].Legend.(string)
			refererJson["referercount"] = strconv.Itoa(dataReferer[cnt].Counter)
			referer = append(referer,refererJson)
		}
		refererString,_ := json.Marshal(referer)


		for cnt:=0;cnt<5;cnt++{
			if cnt == len(dataTimerPerDate){
				break
			}
			accessDayJson := make(map[string]interface{})
			accessDayJson["accessday"]=dataTimerPerDate[cnt].Legend.(time.Time).Format("2006-01-02")
			accessDayJson["accessdaycount"]=strconv.Itoa(dataTimerPerDate[cnt].Counter)
			accessDay = append(accessDay, accessDayJson)
		}
		accessDayString,_ := json.Marshal(accessDay)
		
		for cnt:=0;cnt<3;cnt++{
			if cnt == len(dataTimePerMinute){
				break
			}
			accessMinJson := make(map[string]interface{})
			accessMinJson["accessmin"]=dataTimePerMinute[cnt].Legend.(time.Time).Format("2006-01-02 15:04:05")
			accessMinJson["accessmincount"]=strconv.Itoa(dataTimePerMinute[cnt].Counter)
			accessMin = append(accessMin, accessMinJson)
		}
		accessMinString,_ := json.Marshal(accessMin)
		
		response.HTML(http.StatusOK, "statement.html", gin.H{
			"accessMin": string(accessMinString),
			"accessDay": string(accessDayString),
			"accessDevice": string(accessDeviceString),
			"referer": string(refererString),
		},
		)

	})
	server.PUT("/:shortkey/:secrettoken", func(response *gin.Context) {
		shortkey := response.Param("shortkey")
		secrettoken := response.Param("secrettoken")
		err = svc.UpdateMapping(shortkey, secrettoken, "https://daum.net")
		//todo: 변경할 URL 어떻게 가져올지 고민하기
	
		if err != nil {
			panic(err)
		}
		response.String(200, "update")
	})
	server.DELETE("/:shortkey/:secrettoken", func(response *gin.Context) {	
		shortkey := response.Param("shortkey")
		secrettoken := response.Param("secrettoken")
		err = svc.RemoveMapping(shortkey, secrettoken)
		if err != nil {
			panic(err)
		}
		response.String(200, "remove")
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

func ReadUserReferer(r *http.Request) string {
	return r.Header.Get("referer")
}

func isMobile(r *http.Request) bool  {
	userAgent:=r.Header.Get("User-Agent")
	if(strings.Contains(userAgent,"Mobi")){
		return true
	}
	return false
}

func getMinuteItemsByIterator(iterator int,starow []*service.QueryStatisticsResultRow) string {
	retList := ""
	for cnt:=0;cnt<iterator;cnt++{
		//fromTime := starow[cnt].Legend.(time.Time).Format()
		//retList += 
		retList += strconv.Itoa(starow[cnt].Counter)
	}
	return retList
}