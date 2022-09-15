//go:build practice

package main

import (
	"fmt"
	"quekr/server/service"
	"time"
)

func main() {
	svc, err := service.NewService()

	if err != nil {
		panic(err)
	}

	info, err := svc.CreateMapping("https://naver.com", "127.0.0.1")

	if err != nil {
		panic(err)
	}

	fmt.Println("created")
	fmt.Println(info)

	err = svc.UpdateMapping(info.ShortKey, "abc", "https://daum.net")

	if err == nil {
		panic("update why succeed?")
	}

	err = svc.UpdateMapping("notexist", "abc", "https://daum.net")

	if err == nil {
		panic("update why succeed?")
	}

	err = svc.UpdateMapping(info.ShortKey, info.SecretToken, "https://daum.net")

	if err != nil {
		panic(err)
	}

	info, err = svc.QueryMapping(info.ShortKey)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s => %s\n", info.ShortKey, info.OriginalUrl)

	err = svc.RemoveMapping(info.ShortKey, info.SecretToken)

	if err != nil {
		panic(err)
	}

	// queuing raw access record for accumlating statistics info
	err = svc.TouchStatistics("abc", time.Now(), "127.0.0.1", "https://daum.net", service.DeviceTypePC)

	if err != nil {
		panic(err)
	}

	fmt.Println("go!")

	starows, err := svc.QueryStatistics("aaa", service.StatisticLegendTypeReferer, false)

	if err != nil {
		panic(err)
	}

	for _, starow := range starows {
		fmt.Printf("%s => %d \n", starow.Legend.(string), starow.Counter)
	}
}
