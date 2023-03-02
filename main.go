package main

import (
	"afd-support/afdian"
	"afd-support/api"
	"afd-support/configs"
	"afd-support/db"
	"afd-support/lib"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var conf *configs.Config

func init() {
	timelocal, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = timelocal

	lib.InitLogger()

	configs.Parse()
	logrus.Info("config.yml 加载成功")
	conf = configs.Conf
	if configs.Conf.Self.Debug {
		logrus.Info("已启用 debug")
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func handleRequest() {
	aFapi := api.New(conf.API.APIToken, conf.Afdian.APIToken, conf.Afdian.UserID)

	http.HandleFunc(conf.API.OrderEndpoint, aFapi.GetOrders)
	http.HandleFunc(conf.API.SponsorsEndpoint, aFapi.GetSponsorsWithPage)
	http.HandleFunc(conf.API.GetuseridEndpoint, aFapi.GetUserIDByURL)

	if conf.Webhook.Enabled {
		http.HandleFunc(conf.Webhook.Endpoint, afdian.Webhook)
	}

	logrus.Fatal(http.ListenAndServe(
		fmt.Sprintf("%s:%s", conf.Self.Host, strconv.Itoa(conf.Self.Port)), nil,
	))
}

func main() {
	logrus.Info("将运行于端口: ", conf.Self.Port)

	db.InitDB()
	defer db.DbS.Close()

	sche := lib.Scheduler
	_ = sche.AddFunc("0 0 0/3 * * ? ", func() {
		db.DbSync()
	})

	lib.InitSchedule()

	handleRequest()
}
