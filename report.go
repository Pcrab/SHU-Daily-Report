package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func report(student Student) bool {
	defer wg.Done()
	tryTimes := 0
	for {
		tryTimes++

		if tryTimes >= 5 {
			return false
		}

		// Check if student has cookie or not. If student doesn't have cookie, set it.
		if !strings.Contains(student.Cookie, ".ncov") {
			NormalLogin(&student)
		} else if !CookieLogin(&student) {
			// Check if student's cookie is up to date. If not, correct it.
			NormalLogin(&student)
		}

		fmt.Println(student.Username, "reporting...")

		viewstate, viewstategenerator := getViewState(student)
		fstate, addr := getFState(student)

		reportURL := "https://selfreport.shu.edu.cn/DayReport.aspx"
		req, err := http.NewRequest("GET", reportURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		parseReq(req)
		req.Header.Set("Cookie", student.Cookie)

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		form := url.Values{}

		form.Add("__EVENTTARGET", "p1$ctl00$btnSubmit")
		form.Add("__EVENTARGUMENT", "")
		form.Add("__VIEWSTATE", viewstate)
		form.Add("__VIEWSTATEGENERATOR", viewstategenerator)
		form.Add("p1$ChengNuo", "p1_ChengNuo")
		form.Add("p1$BaoSRQ", time.Now().Format("2006-01-02"))
		form.Add("p1$DangQSTZK", "良好")
		form.Add("p1$TiWen", "")
		form.Add("p1$JiuYe_ShouJHM", "")
		form.Add("p1$JiuYe_Email", "")
		form.Add("p1$JiuYe_Wechat", "")
		form.Add("p1$QiuZZT", "")
		form.Add("p1$JiuYKN", "")
		form.Add("p1$JiuYSJ", "")
		form.Add("p1$GuoNei", "国内")
		form.Add("p1$ddlGuoJia$Value", "-1")
		form.Add("p1$ddlGuoJia", "选择国家")
		form.Add("p1$ShiFSH", "是")
		form.Add("p1$ShiFZX", "否")
		form.Add("p1$ddlSheng$Value", addr[0])
		form.Add("p1$ddlSheng", addr[0])
		form.Add("p1$ddlShi$Value", addr[1])
		form.Add("p1$ddlShi", addr[1])
		form.Add("p1$ddlXian$Value", addr[2])
		form.Add("p1$ddlXian", addr[2])
		form.Add("p1$XiangXDZ", addr[3])
		form.Add("p1$ShiFZJ", "是")
		form.Add("p1$FengXDQDL", "否")
		form.Add("p1$TongZWDLH", "否")
		form.Add("p1$CengFWH", "否")
		form.Add("p1$CengFWH_RiQi", "")
		form.Add("p1$CengFWH_BeiZhu", "")
		form.Add("p1$JieChu", "否")
		form.Add("p1$JieChu_RiQi", "")
		form.Add("p1$JieChu_BeiZhu", "")
		form.Add("p1$TuJWH", "否")
		form.Add("p1$TuJWH_RiQi", "")
		form.Add("p1$TuJWH_BeiZhu", "")
		form.Add("p1$QueZHZJC$Value", "否")
		form.Add("p1$QueZHZJC", "否")
		form.Add("p1$DangRGL", "否")
		form.Add("p1$GeLDZ", "")
		form.Add("p1$FanXRQ", "")
		form.Add("p1$WeiFHYY", "")
		form.Add("p1$ShangHJZD", "")
		form.Add("p1$DaoXQLYGJ", "")
		form.Add("p1$DaoXQLYCS", "")
		form.Add("p1$JiaRen_BeiZhu", "")
		form.Add("p1$SuiSM", "绿色")
		form.Add("p1$LvMa14Days", "是")
		form.Add("p1$Address2", "")
		form.Add("F_TARGET", "p1_ctl00_btnSubmit")
		form.Add("p1_ContentPanel1_Collapsed", "true")
		form.Add("p1_GeLSM_Collapsed", "false")
		form.Add("p1_Collapsed", "false")
		form.Add("F_STATE", fstate)

		// formString := "__EVENTTARGET="+url.QueryEscape("p1$ctl00$btnSubmit")+"&__EVENTARGUMENT="+"&__VIEWSTATE="+url.QueryEscape("qCY0WrBVh9d3LR8z3GDhdYc27rp6oyNYvEcQm5fPnypqJdsj8VGvrlCpjWLuqnOIssegmhqnBYbMcMvFoh5Hx+u4GwmFcPN/1oVknLi/r3j8XLk0Iq2rvvDO8aFGYyGo8ezNV+H8RDSoD5gCbBnQCxKQ8KB2r+HRIa3Is0OBKsPVy37d+sTd8a0WZ6uxgKnk+cWhrra2SgAwhKa+LrP2+C5x+KrOErbavOESBRHXQ/juwGkUmJZjuuN0y0YWX32tG1lqdXTIqWGIptPnLdLwotaP4lNytSZnbnUbVpd7Q0VdfO42eB2dje/ehvk2PaGVdCnafi6zwbT6bBuEZYn0R8d2Vra1H/HmHcVW/JzRQKe8IUG9z8RVKIf7MEQs5YnAux91e82mwLWehuWlVdTMq9Qy4NvQpTBlwUXf2k01hsi/kYpzSKFrjfQgg79uAXx6gpwqUgp8I4ISjviEOR6Lz8jetolDNq9ztXvmLRkXAPRYwc8lbfm0FssZh4h2iOlft/TejD+MrHlcnb1yVLv+F7Em6WnLjjRV3R1/jRZJvBjEcZIYq9qbc2+6CNi5CJtYXIUVwWL4XAczrzQOws43J9MI1NF6q35c4SOiDykNnDkwn9OuzgBeW9T4IC2g+bnKGV2i+FRqdtqXBQhyimaRiH/dwSl/GciSx6oUY99QysoXIBSrpSGSc4dBvp4YylDu73gZby9ZcwWHBoYSRrR5dKo2BQ6pcqX+q3f/ErRzEupGzwMy/tN6+qId1STxBZBO8+pgwytEIXm7Ry/BW1J52isA8H14ggHknVmL2gGadRVyh5rSL7FE02zXKC2TrC+JbiVUFo96VD56QqxzGS4O4P2rPXGciLPs7Qnj/J36CrHXpFQR4URwAS/s+eqxG6v9JJOFw1dgAXLwbguYzTQTA+NpFeTwA9+rjLutSNZKBFu42vRFv8fQKra+z34uQx1AadEqYrcn16lohEn4Hv1e8aYMF2/KT/GoVjw4DAvC8snKnHmvAlVYVxwuquF7kDpdJt528PqjWxjwgtOs8Cd+q2grj16Bvcu0oOU5+vg7Ns549qQRbv51d3AwVFs5u2bt")+"&__VIEWSTATEGENERATOR="+url.QueryEscape("7AD7E509")+"&p1$ChengNuo="+url.QueryEscape("p1_ChengNuo")+"&p1$BaoSRQ="+url.QueryEscape("2021-02-07")+"&p1$DangQSTZK="+url.QueryEscape("良好")+"&p1$TiWen="+"&p1$JiuYe_ShouJHM="+"&p1$JiuYe_Email="+"&p1$JiuYe_Wechat="+"&p1$QiuZZT="+"&p1$JiuYKN="+"&p1$JiuYSJ="+"&p1$GuoNei="+url.QueryEscape("国内")+"&p1$ddlGuoJia$Value="+url.QueryEscape("-1")+"&p1$ddlGuoJia="+url.QueryEscape("选择国家")+"&p1$ShiFSH="+url.QueryEscape("是")+"&p1$ShiFZX="+url.QueryEscape("否")+"&p1$ddlSheng$Value="+url.QueryEscape("上海")+"&p1$ddlSheng="+url.QueryEscape("上海")+"&p1$ddlShi$Value="+url.QueryEscape("上海市")+"&p1$ddlShi="+url.QueryEscape("上海市")+"&p1$ddlXian$Value="+url.QueryEscape("静安区")+"&p1$ddlXian="+url.QueryEscape("静安区")+"&p1$XiangXDZ="+url.QueryEscape("武定路1145弄2号2604室")+"&p1$ShiFZJ="+url.QueryEscape("是")+"&p1$FengXDQDL="+url.QueryEscape("否")+"&p1$TongZWDLH="+url.QueryEscape("否")+"&p1$CengFWH="+url.QueryEscape("否")+"&p1$CengFWH_RiQi="+"&p1$CengFWH_BeiZhu="+"&p1$JieChu="+url.QueryEscape("否")+"&p1$JieChu_RiQi="+"&p1$JieChu_BeiZhu="+"&p1$TuJWH="+url.QueryEscape("否")+"&p1$TuJWH_RiQi="+"&p1$TuJWH_BeiZhu="+"&p1$QueZHZJC$Value="+url.QueryEscape("否")+"&p1$QueZHZJC="+url.QueryEscape("否")+"&p1$DangRGL="+url.QueryEscape("否")+"&p1$GeLDZ="+"&p1$FanXRQ="+"&p1$WeiFHYY="+"&p1$ShangHJZD="+"&p1$DaoXQLYGJ="+"&p1$DaoXQLYCS="+"&p1$JiaRen_BeiZhu="+"&p1$SuiSM="+url.QueryEscape("绿色")+"&p1$LvMa14Days="+url.QueryEscape("是")+"&p1$Address2="+"&F_TARGET="+url.QueryEscape("p1_ctl00_btnSubmit")+"&p1_ContentPanel1_Collapsed="+url.QueryEscape("true")+"&p1_GeLSM_Collapsed="+url.QueryEscape("false")+"&p1_Collapsed="+url.QueryEscape("false")+"&F_STATE="+url.QueryEscape(fstate)

		req, err = http.NewRequest("POST", reportURL, strings.NewReader(form.Encode()))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Accept", "text/plain, */*; q=0.01")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

		// req.Header.Set("Cache-Control", "max-age=0")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Cookie", student.Cookie)
		req.Header.Set("Host", "selfreport.shu.edu.cn")

		req.Header.Set("Origin", "https://selfreport.shu.edu.cn")
		req.Header.Set("Referer", "https://selfreport.shu.edu.cn/DayReport.aspx")

		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("X-FineUI-Ajax", "true")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")

		resp, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if strings.Contains(string(body), "运行时") {
			fmt.Println("fail")
		} else {
			bodyString := string(UGZipBytes(body))
			if strings.Contains(bodyString, "成功") {
				fmt.Println("!!!!!!!!!")
			}
		}

		// fmt.Println(string(body))


		break

	}
	return true
}
