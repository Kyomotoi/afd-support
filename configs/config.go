package configs

import (
	_ "embed"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

//go:embed default_config.yml
var defConfig string

var Conf *Config

func genConfig() error {
	sb := strings.Builder{}
	sb.WriteString(defConfig)
	err := os.WriteFile("config.yml", []byte(sb.String()), 0644)
	if err != nil {
		return err
	}
	return nil
}

func Parse() {
	content, err := os.ReadFile("config.yml")
	if err != nil {
		err = genConfig()
		if err != nil {
			panic("无法生成设置文件: config.yml, 请确认是否给足系统权限")
		}
		logrus.Warn("未检测到 config.yml，已自动于同目录生成，请配置并重新启动")
		logrus.Warn("将于 5 秒后退出...")
		os.Exit(-1)
	}

	Conf = &Config{}
	err = yaml.Unmarshal(content, Conf)
	if err != nil {
		logrus.Fatal("解析 config.yml 失败，请检查格式、内容是否输入正确")
	}
}
