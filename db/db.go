package db

import (
	"afd-support/configs"
	"afd-support/lib"
	"fmt"
	"os"
	"strconv"
	"time"

	sqlite "github.com/FloatTech/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DbM *gorm.DB
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

// InitDB 初始化数据库
func InitDB() {
	dbConf := configs.Conf.Db
	dbPw := ""
	if dbConf.Password != "" {
		dbPw = fmt.Sprintf(":%s", dbPw)
	}
	dsn := fmt.Sprintf("%s%s@tcp(%s:%s)/%s", dbConf.Username, dbPw, dbConf.Host, strconv.Itoa(dbConf.Port), dbConf.DbName)
	logrus.Debug(fmt.Sprintf("DSN INFO: %s", dsn))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DbM = db

	sqlitePath := dbDir + "sqlite.db"
	DbS = &sqlite.Sqlite{DBPath: sqlitePath}
	err = DbS.Open(time.Minute * 15)
	if err != nil {
		panic(err)
	}

	if !lib.IsExists(sqlitePath) {
		DbS.Create("afdian_users", &AfdianUsers{})
		if db.Error != nil {
			panic(db.Error)
		}
		DbS.Create("afdian_orders", &AfdianOrders{})
		if db.Error != nil {
			panic(db.Error)
		}
	}
}
