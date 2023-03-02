package lib

import (
	"time"

	"github.com/robfig/cron"
)

var Scheduler = cron.NewWithLocation(time.Local)

// InitSchedule 初始化定时任务
func InitSchedule() {
	Scheduler.Start()
}
