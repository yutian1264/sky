package main

import (
	"sky/com/utils/dbUtils"
	"fmt"
)

	func init() {
		err:=utils.InitMySQLDB("root","wang","localhost:3306","qgmanager")
		if err!=nil{
			fmt.Println("mysql connected error")
			return
		}
	}
	func main() {


	}

