//go:build practice

package main

import (
	"fmt"
	"time"
	"quekr/server/service"
)

func main() {
	svc, err := service.NewService()
	// define service
	if err != nil {
		panic(err)
	}

	info, err := svc.CreateMapping("https://naver.com", "127.0.0.1")
	// make s	hortUrl
	if err != nil {
		panic(err)
	}

	fmt.Println("created")
	fmt.Println(info)
	
	err = svc.UpdateMapping(info.ShortKey, "abc", "https://daum.net")
	// this is bad example
	if err == nil {
		panic("update why succeed?")
	}

	err = svc.UpdateMapping("notexist", "abc", "https://daum.net")
	// this is bad example
	if err == nil {
		panic("update why succeed?")
	}
	
	err = svc.UpdateMapping(info.ShortKey, info.SecretToken, "https://daum.net")
	
	if err != nil {
		panic(err)
	}

	info, err = svc.QueryMapping(info.ShortKey)
	// get original url
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s => %s\n", info.ShortKey, info.OriginalUrl)

	// queuing raw access record for accumlating statistics info
	err = svc.TouchStatistics(info.ShortKey, svc.NowLocalTime(), "127.0.0.1", "https://referer", service.DeviceTypePC)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 40)

	starows, err := svc.QueryStatistics(info.ShortKey, service.StatisticLegendTypeReferer, false)
	// getStatistics information (false = asc or desc)
	if err != nil {
		panic(err)
	}

	for _, starow := range starows {
		fmt.Printf("%s => %d \n", starow.Legend.(string), starow.Counter)
	}

	err = svc.RemoveMapping(info.ShortKey, info.SecretToken)
	// remove mapping
	if err != nil {
		panic(err)
	}
}
