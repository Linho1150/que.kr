//go:build statistics_handler

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

	err = svc.TouchStatistics("abc", time.Now(), "127.0.0.1", "https://daum.net", service.DeviceTypePC)

	if err != nil {
		panic(err)
	}

	fmt.Println("go!")
}
