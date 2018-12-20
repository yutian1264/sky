/*
@Time : 2018/10/30 10:10 
@Author : sky
@Description:
@File : pool
@Software: GoLand
*/
package main

import (
	"bytes"
	"net/http"
	"fmt"
	"io/ioutil"
	"time"
	"github.com/yutian1264/sky/com/threadpool"

	//"runtime"
)

func SamplePost(url string,param string,index int)string {

	reader := bytes.NewReader([]byte(param))

	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	fmt.Println(index,"===="+string(respBytes))
	////byte数组直接转成string，优化内存
	//str := (*string)(unsafe.Pointer(&respBytes))
	//fmt.Println(*str)

	return string(respBytes)

}

func main() {

	time.Sleep(1 * time.Second)
	go addQueue()
	time.Sleep(10000 * time.Second)
}

type Abc struct {}

func(a Abc)DoTask(v interface{}){

	fmt.Println("interface-------",v)
	SamplePost("http://localhost:8500/f/d/common",
		`{"tokenid":"user","itemid":"login","param":["wyt","111"]}`,1)
}

func addQueue() {

	for i := 0; i <10; i++ {
		var v interface{}=i
		threadpool.JobQueue <- threadpool.Job{Abc{},&v}
		//fmt.Println(runtime.NumGoroutine())

	}

}