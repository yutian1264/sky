package main

import (
	"sky/com/utils"
)

func main() {
	//file, err := os.Open("F:/切图/bg.jpg")
	utils.CreateThumb("F:/切图/bg.jpg","F:/thumb","abc.jpg",100,0)
}