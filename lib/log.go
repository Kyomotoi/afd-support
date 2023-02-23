package lib

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

const logDir = "data/logs/"

func init() {
	exi := IsDir(logDir)
	if !exi {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			panic("创建日志缓存文件夹失败，请尝试手动创建：data/logs")
		}
	}
}

func InitLogger() {
	logrus.SetFormatter(&easy.Formatter{
		TimestampFormat: "01-02 15:04:05",
		LogFormat:       "%time% | %lvl%: %msg%\n",
	})
	now := time.Now().Format("20060102-15")
	fileName := now + ".log"
	file, _ := os.OpenFile(logDir+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	mw := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(mw)
}
