package handler

import (
	"net/http"

	"github.com/dzlzh/go-vercel/job"
	"github.com/golang-module/carbon/v2"
)

func init() {
	carbon.SetDefault(carbon.Default{
		Layout: carbon.DateTimeLayout,
		Timezone: carbon.PRC,
		WeekStartsAt: carbon.Monday,
		Locale: "zh-CN",
	})
}

func Cron(w http.ResponseWriter, r *http.Request) {
	bond := new(job.Bond)
	bond.Run()

	reminder := new(job.Reminder)
	reminder.Run()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.Host))
	w.Write([]byte(r.Method))
	w.Write([]byte("ok"))
}
