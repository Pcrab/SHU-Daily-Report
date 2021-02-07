package main

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

func NormalLogin(student *Student) bool {

	fmt.Println(student.Username,"is in Normal Login mode")

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// init client
	client := http.Client{Jar: jar}

	// want to log into this URL
	wantURL := "https://selfreport.shu.edu.cn"

	// init request to get the login url (begin with newsso.shu.edu.cn), and put it into loginURL
	req, err := http.NewRequest("GET", wantURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	parseReq(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	loginURL := resp.Request.URL.String()

	cookie := jar.Cookies(resp.Request.URL)[0].String()

	// init form used in the request POST
	form := url.Values{}

	form.Add("username", student.Username)
	form.Add("password", student.Password)
	form.Add("login_submit", "")

	// init request which is used to check if we can login correctly
	req, err = http.NewRequest("POST", loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	parseReq(req)

	// 33 means length of "username=...&password=...&login_submit="
	length := strconv.Itoa(len(student.Username) + len(student.Password) + 33)

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Length", length)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "newsso.shu.edu.cn")
	req.Header.Set("Origin", "https://newsso.shu.edu.cn")
	req.Header.Set("Referer", loginURL)
	req.Header.Set("User-Agent", userAgent)

	// fmt.Println(req)

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	cookie = jar.Cookies(resp.Request.URL)[0].String()

	v := viper.New()
	v.SetConfigName(student.Config)
	v.SetConfigType(configExtension)
	v.AddConfigPath(configPath)
	v.Set("Cookie", cookie)
	v.Set("Username", student.Username)
	v.Set("Password", student.Password)
	err = v.WriteConfig()
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(student.Username, "ugziping...")
	body = UGZipBytes(body)

	if strings.Contains(string(body), "健康之路") {
		fmt.Println("login succeeded")
		student.Cookie = cookie
		return true
	}

	return false
}

