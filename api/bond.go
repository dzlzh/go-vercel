package api

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dzlzh/httpc"
	"github.com/dzlzh/notify"
	"github.com/tidwall/gjson"
)

func Ф(w http.ResponseWriter, r *http.Request) {
	request := httpc.NewRequest(httpc.NewClient())
	request.SetMethod("GET").SetURL("https://datacenter-web.eastmoney.com/api/data/v1/get")
	request.SetQuery("source", "WEB")
	request.SetQuery("client", "WEB")
	request.SetQuery("reportName", "RPT_BOND_CB_LIST")
	request.SetQuery("columns", "SECURITY_CODE,SECURITY_NAME_ABBR,PUBLIC_START_DATE,BOND_START_DATE,LISTING_DATE")
	request.SetQuery("sortColumns", "PUBLIC_START_DATE")
	request.SetQuery("sortTypes", "-1")
	request.SetQuery("pageSize", "50")
	request.SetQuery("pageNumber", "1")
	request.Send()
	_, res, err := request.End()
	if err != nil || !gjson.ValidBytes(res) {
		w.Write([]byte("获取数据失败"))
	}

	result := gjson.ParseBytes(res)
	datas := result.Get("result.data").Array()
	now := time.Now().Format("2006-01-02 00:00:00")
	var publicStart, bondStart, listing string
	for _, data := range datas {
		m := data.Get("SECURITY_CODE").String() + "|" + data.Get("SECURITY_NAME_ABBR").String() + "\n"
		publicStartDate := data.Get("PUBLIC_START_DATE").String()
		bondStartDate := data.Get("BOND_START_DATE").String()
		listingDate := data.Get("LISTING_DATE").String()
		if now == publicStartDate {
			publicStart += m
		}
		if now == bondStartDate {
			bondStart += m
		}
		if now == listingDate {
			listing += m
		}
	}

	n := notify.New()
	corpid := os.Getenv("WXW_CORPID")
	agentid := os.Getenv("WXW_AGENTID")
	corpsecret := os.Getenv("WXW_CORPSECRET")
	if corpid != "" && agentid != "" && corpsecret != "" {
		n.UseService(notify.NewWeiXinWork(corpid, agentid, corpsecret))
	}

	var message strings.Builder
	message.Reset()
	message.WriteString(time.Now().Format("2006-01-02"))
	message.WriteString("\n")
	message.WriteString("新债申购")
	message.WriteString("\n")
	message.WriteString(publicStart)
	message.WriteString("\n")
	message.WriteString("新债中签")
	message.WriteString("\n")
	message.WriteString(bondStart)
	message.WriteString("\n")
	message.WriteString("新债上市")
	message.WriteString("\n")
	message.WriteString(listing)
	n.Send("新债", message.String())
	w.Write([]byte(message.String()))
}
