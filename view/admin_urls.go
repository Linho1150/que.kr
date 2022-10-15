package view

import (
	"encoding/json"
	"errors"
	"net/http"
	"quekr/server/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetStaticData(c *gin.Context, svc *service.Service, shortKey string) ([]*service.QueryStatisticsResultRow,[]*service.QueryStatisticsResultRow,[]*service.QueryStatisticsResultRow,[]*service.QueryStatisticsResultRow,error) {
	dataReferer, err := svc.QueryStatistics(shortKey, service.StatisticLegendTypeReferer, false)
	if err != nil {
		
		c.HTML(http.StatusServiceUnavailable, "error.html", gin.H{"error": "Server Error"})
		err = errors.New("dataReferer")
		return nil,nil,nil,nil,err
	}

	dataDeviceType, err := svc.QueryStatistics(shortKey, service.StatisticLegendTypeDevicetype, false)
	if err != nil {
		c.HTML(http.StatusServiceUnavailable, "error.html", gin.H{"error": "Server Error"})
		return nil,nil,nil,nil,err
	}

	dataTimerPerDate, err := svc.QueryStatistics(shortKey, service.StatisticLegendTypeTimePerDate, false)
	if err != nil {
		c.HTML(http.StatusServiceUnavailable, "error.html", gin.H{"error": "Server Error"})
		return nil,nil,nil,nil,err
	}

	dataTimePerMinute, err := svc.QueryStatistics(shortKey, service.StatisticLegendTypeTimePerMinute, false)
	if err != nil {
		c.HTML(http.StatusServiceUnavailable, "error.html", gin.H{"error": "Server Error"})
		return nil,nil,nil,nil,err
	}
	return dataReferer, dataDeviceType, dataTimerPerDate, dataTimePerMinute,nil
}

func RefineData(svc *service.Service,dataDeviceType []*service.QueryStatisticsResultRow,dataReferer []*service.QueryStatisticsResultRow,dataTimerPerDate []*service.QueryStatisticsResultRow,dataTimePerMinute []*service.QueryStatisticsResultRow)([]byte,[]byte,[]byte,[]byte,error){
	accessDeviceString,err:=RefineDeviceData(dataDeviceType)
	if(err != nil){
		return nil,nil,nil,nil,err
	}	

	refererString,err := RefineRefererData(dataReferer)
	if(err != nil){
		return nil,nil,nil,nil,err
	}
	accessDayString,err := RefineAccessDayData(svc,dataTimerPerDate)
	if(err != nil){
		return nil,nil,nil,nil,err
	}
	accessMinString,err := RefineAccessMinuteData(svc,dataTimePerMinute)
	if(err != nil){
		return nil,nil,nil,nil,err
	}

	return accessDeviceString, refererString, accessDayString, accessMinString, nil
}

func RefineDeviceData(dataDeviceType []*service.QueryStatisticsResultRow)([]byte,error){
	accessDevice := make([]map[string]interface{}, 0, 0)
	for cnt:=0;cnt<3;cnt++{
		if cnt == len(dataDeviceType){
			break
		}
		deviceTypeJson := make(map[string]interface{})
		deviceTypeJson["devicetype"] = dataDeviceType[cnt].Legend.(string)
		deviceTypeJson["devicecount"] = strconv.Itoa(dataDeviceType[cnt].Counter)
		accessDevice = append(accessDevice, deviceTypeJson)
	}
	accessDeviceString,err := json.Marshal(accessDevice)
	if(err != nil){
		err := errors.New("It's not JSON type(Marshal)")
		return nil,err
	}
	return accessDeviceString,nil
}

func RefineRefererData(dataReferer []*service.QueryStatisticsResultRow)([]byte,error){
	referer := make([]map[string]interface{}, 0, 0)
	for cnt:=0;cnt<3;cnt++{
		if cnt == len(dataReferer){
			break
		}
		refererJson := make(map[string]interface{})
		refererJson["refererurl"] = dataReferer[cnt].Legend.(string)
		refererJson["referercount"] = strconv.Itoa(dataReferer[cnt].Counter)
		referer = append(referer,refererJson)
	}
	refererString,err := json.Marshal(referer)
	if(err != nil){
		err := errors.New("It's not JSON type(Marshal)")
		return nil,err
	}
	return refererString, nil
}

func RefineAccessDayData(svc *service.Service,dataTimerPerDate []*service.QueryStatisticsResultRow)([]byte,error){
	accessDay := make([]map[string]interface{}, 0, 0)
	for cnt:=0;cnt<5;cnt++{
		if cnt == len(dataTimerPerDate){
			break
		}
		accessDayJson := make(map[string]interface{})
		accessDayJson["accessday"]=dataTimerPerDate[cnt].Legend.(time.Time).In(svc.LocalTimezone).Format("2006-01-02")
		accessDayJson["accessdaycount"]=strconv.Itoa(dataTimerPerDate[cnt].Counter)
		accessDay = append(accessDay, accessDayJson)
	}
	accessDayString,err := json.Marshal(accessDay)
	if(err != nil){
		err := errors.New("It's not JSON type(Marshal)")
		return nil,err
	}
	return accessDayString,err
}

func RefineAccessMinuteData(svc *service.Service,dataTimePerMinute []*service.QueryStatisticsResultRow)([]byte,error){
	accessMin := make([]map[string]interface{}, 0, 0)
	for cnt:=0;cnt<3;cnt++{
		if cnt == len(dataTimePerMinute){
			break
		}
		accessMinJson := make(map[string]interface{})
		accessMinJson["accessmin"]=dataTimePerMinute[cnt].Legend.(time.Time).In(svc.LocalTimezone).Format("2006-01-02 15:04:05")
		accessMinJson["accessmincount"]=strconv.Itoa(dataTimePerMinute[cnt].Counter)
		accessMin = append(accessMin, accessMinJson)
	}
	accessMinString,err := json.Marshal(accessMin)
	if(err != nil){
		err := errors.New("It's not JSON type(Marshal)")
		return nil,err
	}
	return accessMinString,err
}