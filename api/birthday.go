package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dzlzh/httpc"
	"github.com/dzlzh/notify"
	"github.com/golang-module/carbon/v2"
)

type birthday struct {
	Name string
	Desc string
	Month int
	Day int
	Lunar bool
}

func RunBirthday() {
	res := getConfig()
	var b []birthday
	err := json.Unmarshal(res, &b)
	if err != nil {
		fmt.Println("解析生日错误")
	}
	currentTime := carbon.Now()
	threeDaysLater := currentTime.AddDays(3)

	currentTimeLunar := currentTime.Lunar()
	threeDaysLaterLunar := threeDaysLater.Lunar()

	var message strings.Builder
	message.Reset()
	message.WriteString("生日提醒\n\n")

	for _, v := range b {
		if v.Lunar {
			if v.Month == currentTimeLunar.Month() && v.Day == 30 && !is30Days(currentTimeLunar.Year(), v.Month) {
				v.Day = 29
				fmt.Println(v.Name)
			}
			if v.Month == currentTimeLunar.Month() && v.Day == currentTimeLunar.Day() {
				message.WriteString(v.Desc)
				message.WriteString(" - 今天生日")
				message.WriteString("\n")
			} else if v.Month == threeDaysLaterLunar.Month() && v.Day == threeDaysLaterLunar.Day() {
				message.WriteString(v.Desc)
				message.WriteString(" - 3天后生日")
				message.WriteString("\n")
			}
		} else {
			if v.Month == currentTime.Month() && v.Day == currentTime.Day() {
				message.WriteString(v.Desc)
				message.WriteString(" - 今天生日")
				message.WriteString("\n")
			} else if v.Month == threeDaysLater.Month() && v.Day == threeDaysLater.Day() {
				message.WriteString(v.Desc)
				message.WriteString(" - 3天后生日")
				message.WriteString("\n")
			}
		}
	}
	fmt.Println(message.String())

	sendNotify("生日提醒", message.String())
}

func getConfig() []byte {
	url := os.Getenv("EDGE_CONFIG_URL")
	token := os.Getenv("EDGE_CONFIG_TOKEN")
	request := httpc.NewRequest(httpc.NewClient())
	request.SetMethod("GET").SetURL(url)
	request.SetQuery("token", token)
	request.Send()
	_, res, err := request.End()
	fmt.Println(string(res))
	if err != nil {
		fmt.Println("获取数据失败")
	}
	return res
}

func sendNotify(subject, message string) {
	n := notify.New()
	corpid := os.Getenv("WXW_WC_CORPID")
	agentid := os.Getenv("WXW_WC_AGENTID")
	corpsecret := os.Getenv("WXW_WC_CORPSECRET")
	if corpid != "" && agentid != "" && corpsecret != "" {
		n.UseService(notify.NewWeiXinWork(corpid, agentid, corpsecret))
	}
	n.Send(subject, message)
}

func is30Days(year, month int) bool {
	return carbon.CreateFromLunar(year, month, 29, 0, 0, 0, false).
		AddDay().Lunar().Day() == 30
}
