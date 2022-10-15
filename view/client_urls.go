package view

import (
	"errors"
	"quekr/server/service"

	"github.com/gin-gonic/gin"
)

func CreatUrl(c *gin.Context,svc *service.Service,originalURL string) (string, error){
	if ValidationURL(originalURL) {
		err := errors.New("ValidationURL")
		return "", err
	}
	ipAddress := ReadUserIP(c.Request)
	info, err := svc.CreateMapping(originalURL, ipAddress)
	if err != nil {
		err := errors.New("CraetMapping")
		return "", err
	}
	responseHTML :=
		`<html>
		<form action="/urls" method="POST">
			<input type="text" name="innerUrl" value="https://que.kr/` + info.ShortKey + `"/>
			<input type="text" name="adminUrl" value="https://que.kr/` + info.ShortKey + "/" + info.SecretToken + `"/>
		</form>
		<script>document.forms[0].submit();</script>
	</html>`

	return responseHTML,nil
}

func RedirectUrl(c *gin.Context, svc *service.Service,shortKey string) (*service.MappingInfo,error){
	ipAddress := ReadUserIP(c.Request)
	referer := ReadUserReferer(c.Request)

	deviceType := service.DeviceTypePC
	mobileDevice := ReadDeivceMobile(c.Request)
	if(mobileDevice){
		deviceType = service.DeviceTypeMobile
	}

	item, err := svc.QueryMapping(shortKey)
	if err != nil {
		err = errors.New("QueryMapping")
		return nil, err
	}

	if (referer==""){
		referer="Direct access"
	}
	err = svc.TouchStatistics(shortKey, svc.NowLocalTime(), ipAddress, referer, deviceType)
	if err != nil {
		err = errors.New("TouchStatistics")
		return nil, err
	}
	return item,nil;
}