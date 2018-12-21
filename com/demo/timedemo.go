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
)

func init() {
}

func main() {
	//fmt.Println(time.Now())
	scheduler.Cron = scheduler.NewCron()
	go scheduler.Cron.Start()
	//cron.AddFunc(10*int64(time.Second)+time.Now().UnixNano(),run)
	//scheduler.SchedulerInit()
	scheduler.Cron.AddFuncSpaceNumber(2,-1,run)
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
	log.Println("00000")
	//run2()
	//scheduler.Cron.AddFunc(10*int64(time.Second)+time.Now().UnixNano(),run2)
}

func run2(){
	log.Println("run222222")
}
