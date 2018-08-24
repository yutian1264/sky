package upload

/**
	This method supports sharding upload and single file upload.
	If single file upload does not require chunks and chunks attributes.
	Shard uploads support breakpoint continuation
	post上传
	该方法支持分片上传和单文件整体上传,如果单文件上传不需要chunk和chunks属性.
	分片上传支持断点续传

 */

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"os"
	"io"
	"io/ioutil"
	"sky/com/utils/file"
)

func UploadBreakPoint(req *http.Request,key,savePath string,ch chan int){
	req.ParseForm()
	forms:=req.Form
	//b,_:=json.Marshal(forms)
	formFile, header, err := getFileMessage(req,key)
	defer formFile.Close()
	//todo check file message
	if !CheckErr(err){
		return
	}
	fileName:=header.Filename
	var chunk ,totalChunk int
	// is sub_slice
	_, ok := forms["chunk"]
	var saveFileName string
	if ok{
		chunk,_=strconv.Atoi(forms["chunk"][0])
		totalChunk,_=strconv.Atoi(forms["chunks"][0])
		fmt.Println("current point:,total:",chunk,totalChunk)

		tempPath:=savePath+fileName+"_temp/"
		//如果是分包上传需要保存临时文件,上传完成后组合成一个完整文件
		saveFileName=tempPath+header.Filename+strconv.Itoa(chunk)+".g"
		b,_:=file.IsPathExists(saveFileName)
		if b{
			log.Println(saveFileName+":文件已经存在")
			goto Nothing
		}
		b,err:=file.IsPathExists(tempPath)
		if !b{
			err=os.MkdirAll(tempPath,os.ModePerm)
			if err!=nil{
				log.Printf("Create directory error: %s\n", err)
				return
			}
		}
		var isOK bool
		//保存文件
		isOK,err=saveFile(saveFileName,formFile)
		if !isOK{
			log.Println(saveFileName+":保存失败")
			ch<-0;
		}
		ch<-1;

	}else{
		//整体上传
		saveFileName=savePath+header.Filename
		var isOK bool
		//保存文件
		isOK,err=saveFile(saveFileName,formFile)
		if !isOK{
			log.Println(saveFileName+":保存失败")
			ch<-0;
		}
		ch<-1;
	}

	Nothing:{
		log.Println("nothing")
		formFile.Close()
		ch<-1;
	}

	//如果是分片上传判断是否上传完成,然后组合成完整文件
	if ok{
		if chunk==totalChunk-1{
			assembleFile(header.Filename,savePath,totalChunk)
			ch<-1;
		}
	}


}
//组装上传文件
func  assembleFile(fileName,path string,totalCount int){

	fileNameTemp:=path+fileName
	b,_:=file.IsPathExists(fileName)
	if b{
		log.Println("文件已经存在:",fileName)
		removeTempFile(fileNameTemp+"_temp",totalCount)
		return
	}

	file, err := os.OpenFile(fileNameTemp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < totalCount; i++ {
		f, err := os.OpenFile(fileNameTemp+"_temp/"+fileName+strconv.Itoa(i)+".g", os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		file.Write(b)
		f.Close()
	}
	 file.Close()
	removeTempFile(fileNameTemp+"_temp",totalCount)


}
func removeTempFile(filePath string,count int){
	err:=os.RemoveAll(filePath)
	if err!=nil{
		log.Println("删除临时文件异常:",err)
	}

}
//保存临时文件
func saveFile(fileName string,formFile multipart.File)(bool,error){

	destFile, err :=os.OpenFile(fileName,os.O_APPEND|os.O_RDWR|os.O_CREATE,0777)
	if err != nil {
		log.Printf("Create failed: %s\n", err)
		return false,err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, formFile)
	if err != nil {
		log.Printf("Write file failed: %s\n", err)
		return false,err
	}
	return true,nil
}
func CheckErr(err error)bool {
	if err != nil {
		defer func(err error)bool {
			log.Fatal(err)
			recover()
			return false
		}(err)
		panic(err)
	}
	return true
}
func getFileMessage(req *http.Request,key string)(multipart.File, *multipart.FileHeader, error){
	return req.FormFile(key)
}
