package main

import (
	"github.com/yutian1264/sky/com/utils/dbUtils"
	"os"
	"log"
)

	func init() {
		err:=utils.InitMySQLDB("root","wang","localhost:3306","qgmanager")
		if err!=nil{
			log.Println("mysql connected error")
			os.Exit(3)
		}
	}
	func main() {
		result:=utils.Common("user","ONEUSER_SM",[]string{})
		log.Println(result)

	}

