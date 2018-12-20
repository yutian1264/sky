package main

import (
	"github.com/yutian1264/sky/com/utils/imgUtils"
	//"github.com/yutian1264/sky/com/utils"

	"fmt"
)

func main() {

	//生成二维码
	//filename:=time.Now().Unix()
	//path:="d:/upload/barcode"//utils.GetCurrentDirectory()
	//imgUtils.CreateQRCode("abcd",path,strconv.FormatInt(filename,10)+".png",200,200)
	//读取exif 信息
	//str:=[]string{"D:/images/测试场景/测试场景/IMG_8360.JPG"}
	//imgs:=imgUtils.GetExifMess(str);
	//fmt.Println(imgs)
	//获取城市
	m:=imgUtils.GetCityByPoint(39.99811944444445, 116.44526666666667)
	fmt.Println(m)
}
