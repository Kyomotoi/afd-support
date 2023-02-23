package afdian

import (
	"afd-support/db"
	"database/sql"
	"encoding/json"
	"fmt"
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
	logrus.Info(fmt.Sprintf("webhook: 新订单请求，订单号 %s from %s", a.Data.Order.OutTradeNo, a.Data.Order.UserID))
	rw.Write([]byte(`{"ec":200}`))

	var dIns *sql.Stmt

	_, err = db.DoSearch("afdian_users", "user_name", "user_id", a.Data.Order.UserID)
	if err != nil {
		logrus.Warn(fmt.Sprintf("%s 为新用户", a.Data.Order.UserID))
		dIns, _ := db.Database.Prepare("INSERT INTO `afdian_users` (user_id, user_name) VALUES (?, ?)")
		defer dIns.Close()
		_, _ = dIns.Exec(a.Data.Order.UserID, "")
		// TODO: 摆烂
	}

	dIns, _ = db.Database.Prepare("INSERT INTO `afdian_orders` (order_no, time, user_id, consumed) VALUES (?, ?, ?, ?)")
	defer dIns.Close()
	_, err = dIns.Exec(a.Data.Order.OutTradeNo, time.Now().Unix(), a.Data.Order.UserID, 0)
	if err != nil {
		// TODO: 此处为服务器数据库出现问题，由此写入本地 sqlite
		logrus.Error(err)
	}
}
