package afdian

import (
	"afd-support/db"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Webhook(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("webhook 收到新请求")
	a := &AfdianWebhookResponse{}
	body, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(body, a)
	if err != nil {
		logrus.Error("webhook：解析返回数据失败")
		return
	}
	logrus.Info("webhook: 新订单请求，订单号 ", a.Data.Order.OutTradeNo)
	rw.Write([]byte(`{"ec":200}`))

	dIns, _ := db.Database.Prepare("INSERT INTO `afdian_orders` (order_no, time, user_id, consumed)VALUES(?, ?, ?, ?)")
	defer dIns.Close()
	_, err = dIns.Exec(a.Data.Order.OutTradeNo, time.Now().Unix(), a.Data.Order.UserID, 0)
	if err != nil {
		// TODO: 此处为服务器数据库出现问题，由此写入本地 sqlite
		logrus.Error(err)
	}
}
