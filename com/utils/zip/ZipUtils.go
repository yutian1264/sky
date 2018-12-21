/*
@Time : 2018/12/21 14:28 
@Author : sky
@Description:
@File : ZipUtils
@Software: GoLand
*/
package zip

import (
	"os"
	"archive/zip"
	"io"
	"github.com/astaxie/beego/logs"
	"net/url"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/simplifiedchinese"
	"bytes"
	"io/ioutil"
)

func ToZip(files []*os.File, dest string) error {

	//f, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	d, _ := os.Create(dest + ".zip")
	defer d.Close()
	w := zip.NewWriter(d)
	for _, file := range files {
		err := compress(file, dest, w)
		if err != nil {
			return err
		}
	}
	return nil
}
func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		err = os.MkdirAll(info.Name(), os.ModePerm)
		if err != nil {
			logs.Error("create dir : " + err.Error())
			return err
		}

		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

/**
@zipFile：压缩文件
@dest：解压之后文件保存路径
*/func UnZip(srcFile *os.File, dest string) error {
	zipFile, err := zip.OpenReader(srcFile.Name())
	if err != nil {
		logs.Error("Unzip File Error：", err.Error())
		return err
	}
	defer zipFile.Close()
	for _, innerFile := range zipFile.File {
		info := innerFile.FileInfo()
		if info.IsDir() {
			data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(innerFile.Name)),simplifiedchinese.GB18030.NewEncoder()))
			err = os.MkdirAll(dest+"/"+ string(data), os.ModePerm)//url.QueryEscape(
			if err != nil {
				logs.Error("Unzip File Error : " + err.Error())
				return err
			}
			continue
		}
		srcFile, err := innerFile.Open()
		if err != nil {
			logs.Error("Unzip File Error : " + err.Error())
			continue
		}
		defer srcFile.Close()
		newFile, err := os.Create(dest+"/"+url.QueryEscape(innerFile.Name))
		if err != nil {
			logs.Error("Unzip File Error : " + err.Error())
			continue
		}
		io.Copy(newFile, srcFile)
		newFile.Close()
	}
	return nil
}
