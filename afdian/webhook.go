package afdian

import (
	"afd-support/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

	var do *gorm.DB
	var result map[string]interface{} = make(map[string]interface{})
	db.DbM.Model(&db.AfdianUsers{}).First(&result, "user_id = ?", a.Data.Order.UserID)
	if len(result) == 0 {
		logrus.Warn(a.Data.Order.UserID + " 为新用户")
		if do = db.DbM.Create(&db.AfdianUsers{
			UserID:   a.Data.Order.UserID,
			UserName: "",
		}); do.Error != nil {
			logrus.Error("远程数据库写入失败，将写入本地 sqlite 数据库")
			err = db.DbS.Insert("afdian_users", &db.AfdianUsers{
				UserID:   a.Data.Order.UserID,
				UserName: "",
			})
			if err != nil {
				logrus.Error("本地数据库写入失败，将跳过本次记录")
			}
		}
	}

	db.DbM.Model(&db.AfdianOrders{}).First(&result, "order_no = ?", a.Data.Order.OutTradeNo)
	if len(result) == 0 {
		if do = db.DbM.Create(&db.AfdianOrders{
			OrderNo:  a.Data.Order.OutTradeNo,
			Time:     time.Now().Unix(),
			UserID:   a.Data.Order.UserID,
			Consumed: 0,
		}); do.Error != nil {
			logrus.Error("远程数据库写入失败，将写入本地 sqlite 数据库")
			err = db.DbS.Insert("afdian_users", &db.AfdianOrders{
				OrderNo:  a.Data.Order.OutTradeNo,
				Time:     time.Now().Unix(),
				UserID:   a.Data.Order.UserID,
				Consumed: 0,
			})
			if err != nil {
				logrus.Error("本地数据库写入失败，将跳过本次记录")
			}
		}
	} else {
		logrus.Info("该订单已有记录，跳过")
	}
}
