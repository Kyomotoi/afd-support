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
		dIns, _ := db.DbM.Prepare("INSERT INTO `afdian_users` (user_id, user_name) VALUES (?, ?)")
		defer dIns.Close()
		_, _ = dIns.Exec(a.Data.Order.UserID, "")
		// TODO: 摆烂
	}

	dIns, _ = db.DbM.Prepare("INSERT INTO `afdian_orders` (order_no, time, user_id, consumed) VALUES (?, ?, ?, ?)")
	defer dIns.Close()
	_, err = dIns.Exec(a.Data.Order.OutTradeNo, time.Now().Unix(), a.Data.Order.UserID, 0)
	if err != nil {
		logrus.Error("写入 mysql/MariaDB 失败，将写入本地数据库")
		logrus.Error(err)
		err = db.DbS.Insert("afdian_orders", &db.AfdianOrders{
			OrderNo:  a.Data.Order.OutTradeNo,
			Time:     time.Now().Unix(),
			UserID:   a.Data.Order.UserID,
			Consumed: 0,
		})
		if err != nil {
			logrus.Error("写入本地数据库失败，本次记录将跳过")
			logrus.Error(err)
		}
	}
}
