/*
@Time : 2018/12/20 14:14 
@Author : sky
@Description:
@File : timedemo
@Software: GoLand
*/
package main

import (
	"log"
	"github.com/yutian1264/sky/com/utils/scheduler"
	"time"
	"github.com/yutian1264/sky/com/utils/timeUtils"
)

func init() {
}

func main() {
	//fmt.Println(time.Now())
	scheduler.Cron = scheduler.NewCron()
	go scheduler.Cron.Start()
	//cron.AddFunc(10*int64(time.Second)+time.Now().UnixNano(),run)
	//scheduler.SchedulerInit()
	scheduler.Cron.AddFuncSpaceNumber(1,-1,run)
	//cron.AddTask(&scheduler.Task{
	//	Job: scheduler.FuncJob(func() {
	//		fmt.Println("hello cron")
	//	}),
	//	RunTime: time.Now().UnixNano() + 20*int64(time.Second),
	//})
	timer := time.NewTimer(10 * time.Second)
	for {
		select {
			case <-timer.C:
				break
		}

	}
}
func run(){
	t1:=timeUtils.Str2Time("2018-12-21 18:00:00")
	t2:=time.Now()
	local2, _ := time.LoadLocation("Local")
	t1=t1.In(local2)
	t2=t2.In(local2)
	t3:=t1.Sub(t2)
	log.Println("下班倒计时:",int(t3.Seconds()),"秒")
	//run2()
	//scheduler.Cron.AddFunc(10*int64(time.Second)+time.Now().UnixNano(),run2)
}

func run2(){
	log.Println("run222222")
}
