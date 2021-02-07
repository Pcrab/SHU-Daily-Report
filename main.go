package main

import (
	"fmt"
	"net/http"
	"sync"
)

type Student struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Cookie   string `mapstructure:"cookie"`
	Config   string
}

type Host struct {
	Address string
	Port int
}

type Config struct {
	AppId string
	Secret string
	Host Host
}


var (
	wg                       sync.WaitGroup
	students                 []Student
	userAgent                = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15"
	tooManyRequestsSleepTime = 100
	configPath               = "./students/"
	configExtension          = "toml"
)

func main() {


	MakeStudents()

	for _, student := range students {
		wg.Add(1)
	 	go report(student)

	}

	wg.Wait()
	fmt.Println("done!")

	return

}

func parseReq(req *http.Request) {
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", userAgent)
}
