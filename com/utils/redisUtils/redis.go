package utils

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var redisConn redis.Conn

func RedisInit(serverPath,port string)bool {

	b:=true
	var err error
	redisConn, err = redis.Dial("tcp", serverPath+":"+port)
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return false
	}
	//defer redisConn.Close()

	return b
}

func CheckRedis(){

	fmt.Print(redisConn)
}


func RedisSetItem(key,data string)bool{

	b:=true
	_, err:= redisConn.Do("SET", key, data)
	if err != nil {
		b=false
		fmt.Println("redis set failed:", err)
	}
	return b
}

func RedisCheckKeyExists(key string)bool{
	is_key_exit, err := redis.Bool(redisConn.Do("EXISTS", key))
	if err != nil {
		is_key_exit=false;
		fmt.Println("error:", err)
	}
	return is_key_exit;
}

func RedisDeleteItemByKey(key string){

	if RedisCheckKeyExists(key){
		_, err := redisConn.Do("DEL", key)
		if err != nil {
			fmt.Println("redis delelte failed:", err)
		}
	}
}
func RedisAddHash(){

}
func RedisAddList(){

}
func RedisAddSet(){

}
func RedisAddZset(){

}

