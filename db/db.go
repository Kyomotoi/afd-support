package db

import (
	"afd-support/configs"
	"afd-support/lib"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	sqlite "github.com/FloatTech/sqlite"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var (
	DbM *sql.DB
	DbS *sqlite.Sqlite
)

const dbDir = "data/db/"

func init() {
	exi := lib.IsDir(dbDir)
	if !exi {
		err := os.MkdirAll(dbDir, 0755)
		if err != nil {
			panic("创建数据库文件夹失败，请尝试手动创建: data/db")
		}
	}
}

func InitDb() {
	dbConf := configs.Conf.Db
	dbPw := ""
	if dbConf.Password != "" {
		dbPw = fmt.Sprintf(":%s", dbPw)
	}
	dsn := fmt.Sprintf("%s%s@tcp(%s:%s)/%s", dbConf.Username, dbPw, dbConf.Host, strconv.Itoa(dbConf.Port), dbConf.DbName)
	logrus.Debug(fmt.Sprintf("DSN INFO: %s", dsn))

	var err error
	DbM, err = sql.Open("mysql", dsn)
	if err != nil {
		logrus.Error("连接 mysql/MariaDB 数据库时出错")
		logrus.Fatal(err)
	}
	DbM.SetConnMaxLifetime(time.Minute * 3)

	var version string
	DbM.QueryRow("SELECT VERSION()").Scan(&version)
	logrus.Info("数据库版本: " + version)

	// TODO
	dbPath := dbDir + "sqlite.db"
	DbS = &sqlite.Sqlite{DBPath: dbPath}
	err = DbS.Open(time.Minute * 15)
	if err != nil {
		logrus.Error("连接 sqlite 数据库时出错")
		logrus.Fatal(err)
	}
	if !lib.IsExists(dbPath) {
		err = DbS.Create("afdian_orders", &AfdianOrders{})
		if err != nil {
			logrus.Error("初始化 sqlite 时创建表失败")
			logrus.Fatal(err)
		}
	}
}

func DoSearch(t, k1, k2, k3 string) (string, error) {
	do, err := DbM.Prepare(fmt.Sprintf("SELECT %s FROM `%s` WHERE %s = ?", k1, t, k2))
	if err != nil {
		return "", err
	}

	var att string
	err = do.QueryRow(k3).Scan(&att)
	if err != nil {
		return "", err
	}
	return att, nil
}
