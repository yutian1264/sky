package utils

import (
	"os"
	"log"
	"image/jpeg"
	"github.com/nfnt/resize"
	"path/filepath"
	"github.com/astaxie/beego"
	"strings"
)

/*
	校验文件是否存在
 */


func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func PathNotExistsCreate(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		err= os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false, nil
		}
		return true, nil
	}
	return false, err
}

func CreateThumb(path,savePath,saveName string,W,H uint)error {
	file, err := os.Open(path)
	if err != nil {
		recover()
		log.Println("create thumb open file:",err)
		file.Close()
		return err
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		recover()
		log.Println("create thumb decode:",err)
		file.Close()
		return err
	}
	file.Close()
	m := resize.Resize(W, H, img, resize.Lanczos3)
	out, err := os.Create(savePath+"/"+saveName)
	if err != nil {
		log.Println("create thumb:",err)
		return err
	}
	defer out.Close()
	err=jpeg.Encode(out, m, nil)
	if err!=nil{
		recover()
		return err
	}
	return nil
}

/*
获取程序运行路径
*/
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		beego.Debug(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}