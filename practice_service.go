package main

import (
	"fmt"
	"quekr/server/service"
)

func main() {
	service, err := service.NewService()

	if err != nil {
		panic(err)
	}

	info, err := service.CreateMapping("https://naver.com", "127.0.0.1")

	if err != nil {
		panic(err)
	}

	fmt.Println("created")
	fmt.Println(info)

	err = service.UpdateMapping(info.ShortKey, "abc", "https://daum.net")

	if err == nil {
		panic("update why succeed?")
	}

	err = service.UpdateMapping("notexist", "abc", "https://daum.net")

	if err == nil {
		panic("update why succeed?")
	}

	err = service.UpdateMapping(info.ShortKey, info.SecretToken, "https://daum.net")

	if err != nil {
		panic(err)
	}

	info, err := service.QueryMapping(info.ShortKey)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s => %s\n", info.ShortKey, info.OriginalUrl)

	err = service.RemoveMapping(info.ShortKey, info.SecretToken)

	if err != nil {
		panic(err)
	}
}
