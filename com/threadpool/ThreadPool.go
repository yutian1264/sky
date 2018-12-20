/*
@Time : 2018/10/30 9:15
@Author : sky
@Description:
@File : ThreadPoll
@Software: GoLand
*/
package threadpool

import (
	"fmt"
	//"math/rand"
	"runtime"
)

var (
	MaxQueue = 101
	MaxWorker =10
)

type Task interface {
	DoTask(v interface{})
}
type Job struct {
	Task Task
	V interface{}

}


func init(){
	runtime.GOMAXPROCS(runtime.NumGoroutine())
	JobQueue = make(chan Job, MaxQueue)
	dispatcher := NewDispatcher(MaxWorker)
	dispatcher.Run()
}
//var count=0;
//func (j Job) DoTask(){
//	count++
//	//time.Sleep(100 * time.Millisecond)
//	fmt.Println(count,"===================",rand.Float64())
//}

//任务channal
var JobQueue chan Job

//执行任务单元
type Worker struct {
	WorkerPool chan chan Job
	JobChannal chan Job
	quit       chan bool
	no         int
}

//创建一个新的工作单元

func NewWorker(workerPool chan chan Job, no int) Worker {

	//返回一个新的worker 实例不需要地址符
	return Worker{
		workerPool,
		make(chan Job),
		make(chan bool),
		no,
	}
}

//循环监听任务和结束信号
func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannal
			//fmt.Println("w.WorkerPool <- w.JobChannel", w)
			select {
			case job := <-w.JobChannal:
				job.Task.DoTask(job.V)
			case <-w.quit:
				fmt.Println(w.no, "job exit")
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

//任务调度

type Dispatcher struct {
	//工作池
	WorkerPool chan chan Job
	MaxWorker  int
}

func NewDispatcher(maxPoolCount int) *Dispatcher {
	pool := make(chan chan Job, maxPoolCount)
	return &Dispatcher{pool, maxPoolCount}
}

//工作者池的初始化
func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 1; i < d.MaxWorker+1; i++ {
		worker := NewWorker(d.WorkerPool, i)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			//fmt.Println("job := <-JobQueue:")
			go func(job Job) {
				//fmt.Println("等待空闲worker (任务多的时候会阻塞这里")
				//等待空闲worker (任务多的时候会阻塞这里)
				jobChannel := <-d.WorkerPool
				//fmt.Println("jobChannel := <-d.WorkerPool", reflect.TypeOf(jobChannel))
				// 将任务放到上述woker的私有任务channal中
				jobChannel <- job
				//fmt.Println("jobChannel <- job")
			}(job)
		}
	}
}


