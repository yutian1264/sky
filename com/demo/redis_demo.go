package main

import (
	"sky/com/utils/redisUtils"
	"fmt"
)

func main() {

	utils.RedisInit("127.0.0.1","6379")
	utils.RedisSetItem("test","123345")

	fmt.Println(string(utils.RedisGetItem("test").([]byte)))

}

