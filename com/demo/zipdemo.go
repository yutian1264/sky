/*
@Time : 2018/12/21 14:36 
@Author : sky
@Description:
@File : zipdemo
@Software: GoLand
*/
package main

import (
	"os"
	"fmt"
	"github.com/yutian1264/sky/com/utils/zip"
)

func main() {
	f1, err := os.Open("D:/images/测试场景/测试场景/IMG_7209.JPG")
	if err != nil {
		fmt.Println(err)
	}
	defer f1.Close()
	f2, err := os.Open("D:/images/a/")
	if err != nil {
		fmt.Println(err)
	}
	defer f2.Close()
	f3, err := os.Open("D:/images/test_img.zip")
	if err != nil {
		fmt.Println(err)
	}
	defer f3.Close()
	//var files = []*os.File{f1, f2, f3}
	dest := "D:/images/ttt"
	//err = zip.ToZip(files, dest)
	//if err != nil {
	//	fmt.Println(err)
	//}

	zip.UnZip(f3,dest)
}
