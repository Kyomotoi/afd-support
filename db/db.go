package db

import (
	"afd-support/configs"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var Database *sql.DB

func InitDb() {
	dbConf := configs.Conf.Db
	dbPw := ""
	if dbConf.Password != "" {
		dbPw = fmt.Sprintf(":%s", dbPw)
	}
	dsn := fmt.Sprintf("%s%s@tcp(%s:%s)/%s", dbConf.Username, dbPw, dbConf.Host, strconv.Itoa(dbConf.Port), dbConf.DbName)
	logrus.Debug(fmt.Sprintf("DSN INFO: %s", dsn))

	var err error
	Database, err = sql.Open("mysql", dsn)
	if err != nil {
		logrus.Error("连接数据库时出错")
		logrus.Fatal(err)
	}
	Database.SetConnMaxLifetime(time.Minute * 3)

	var version string
	Database.QueryRow("SELECT VERSION()").Scan(&version)
	logrus.Info("数据库版本: " + version)
}

func DoSearch(t, k1, k2, k3 string) (string, error) {
	do, err := Database.Prepare(fmt.Sprintf("SELECT %s FROM `%s` WHERE %s = ?", k1, t, k2))
	if err != nil {
		return "", err
	}
	defer do.Close()

	var att string
	err = do.QueryRow(k3).Scan(&att)
	if err != nil {
		return "", err
	}
	return att, nil
}
