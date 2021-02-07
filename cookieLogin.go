package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func CookieLogin(student *Student) bool {

	fmt.Println(student.Username,"is in Cookie Login mode")

	client := http.Client{}
	// urla := "https://newsso.shu.edu.cn/login/eyJ0aW1lc3RhbXAiOjE2MTI1OTEzMDI3NTc4ODIyNzEsInJlc3BvbnNlVHlwZSI6ImNvZGUiLCJjbGllbnRJZCI6IldVSFdmcm50bldZSFpmelE1UXZYVUNWeSIsInNjb3BlIjoiMSIsInJlZGlyZWN0VXJpIjoiaHR0cHM6Ly9zZWxmcmVwb3J0LnNodS5lZHUuY24vTG9naW5TU08uYXNweD9SZXR1cm5Vcmw9JTJmIiwic3RhdGUiOiIifQ=="
	urla := "https://selfreport.shu.edu.cn"

	req, err := http.NewRequest("GET", urla, nil)
	if err != nil {
		log.Fatal(err)
	}
	parseReq(req)
	req.Header.Set("Cookie", student.Cookie)
	req.Header.Set("Host", "selfreport.shu.edu.cn")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()


	body, err := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(body), "上海大学统一身份认证") {
		fmt.Println(student.Username, "Cookie login failed")
		return false
	}

	body = UGZipBytes(body)

	if strings.Contains(string(body), "健康之路") {
		fmt.Println(student.Username, "login succeeded")
		return true
	}

	return false
}
