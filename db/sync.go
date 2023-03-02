package db

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

// DbSync 将本地所记录的数据同步至远程数据库
func DbSync() error {
	logrus.Warn("开始同步...")

	sqlDB, err := DbM.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	} else {
		logrus.Debug("远程数据库存活")
	}

	auNum, aoNum := 0, 0
	auNum, err = DbS.Count("afdian_users")
	if err != nil {
		logrus.Error("获取本地数据库用户数量失败")
		return err
	}
	aoNum, err = DbS.Count("afdian_orders")
	if err != nil {
		logrus.Error("获取本地数据库订单数量失败")
		return err
	}

	logrus.Info(fmt.Sprintf("本地数据库已记录：%s 位用户, %s 个订单", strconv.Itoa(auNum), strconv.Itoa(aoNum)))
	if auNum != 0 {
		for i := 0; i < auNum; i++ {
			var au *AfdianUsers
			err = DbS.Query("SELECT * FROM afdian_users", &au)
			if err != nil {
				logrus.Error("从本地获取用户数据失败，本次同步失败")
				return err
			}
			result := map[string]interface{}{}
			DbM.Model(&AfdianUsers{}).First(&result, "user_id = ?", au.UserID)
			if result != nil {
				continue
			} else {
				if do := DbM.Model(&AfdianUsers{}).Updates(&AfdianUsers{
					UserID:   au.UserID,
					UserName: au.UserName,
				}); do.Error != nil {
					logrus.Warn(fmt.Sprintf("用户 %s 同步失败", au.UserID))
				} else {
					DbS.Del("afdian_users", fmt.Sprintf("WHERE user_id = %s", au.UserID))
				}
			}
		}
	}
	if aoNum != 0 {
		for i := 0; i < aoNum; i++ {
			var ao *AfdianOrders
			err = DbS.Query("SELECT * FROM afdian_orders", &ao)
			if err != nil {
				logrus.Error("从本地获取用户数据失败，本次同步失败")
				return err
			}
			result := map[string]interface{}{}
			DbM.Model(&AfdianOrders{}).First(&result, "order_no", ao.OrderNo)
			if result != nil {
				continue
			} else {
				if do := DbM.Model(&AfdianUsers{}).Updates(&AfdianOrders{
					OrderNo:  ao.OrderNo,
					Time:     ao.Time,
					UserID:   ao.UserID,
					Consumed: ao.Consumed,
				}); do.Error != nil {
					logrus.Warn(fmt.Sprintf("订单 %s 同步失败", ao.OrderNo))
				} else {
					DbS.Del("afdian_orders", fmt.Sprintf("WHERE order_no = %s", ao.OrderNo))
				}
			}
		}
	}
	logrus.Info("同步完成")
	return nil
}
