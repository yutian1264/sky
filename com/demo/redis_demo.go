package main

import (
	"sky/com/utils/redisUtils"
	"fmt"
)

func main() {

	utils.RedisInit("127.0.0.1","6379","auth_name")
	utils.RedisSetItem("test","123345")

	fmt.Println(utils.RedisGetItem("test"))

}

