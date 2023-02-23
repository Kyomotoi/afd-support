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

	_ "github.com/go-sql-driver/mysql"
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
	aFapi := api.New(conf.Self.AuthToken, conf.HTTP.AfdAPIToken, conf.HTTP.AfdUserID)

	http.HandleFunc("/order", aFapi.GetOrders)
	http.HandleFunc("/sponsors", aFapi.GetSponsorsWithPage)
	http.HandleFunc("/getuserid", aFapi.GetUserIDByURL)

	if conf.Webhook.Enabled {
		http.HandleFunc(conf.Webhook.Point, afdian.Webhook)
	}

	logrus.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", conf.Self.Host, strconv.Itoa(conf.Self.Port)), nil))
}

func main() {
	logrus.Info("将运行于端口: ", conf.Self.Port)

	handleRequest()

	db.InitDb()
	defer db.Database.Close()
}
