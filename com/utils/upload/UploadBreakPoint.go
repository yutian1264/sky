package upload

/**
	This method supports sharding upload and single file upload.
	If single file upload does not require chunks and chunks attributes.
	Shard uploads support breakpoint continuation
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
	"time"
	"github.com/yutian1264/sky/com/utils"
)
type ResultChan struct {
	FileName string
	ThumbPath string
	FilePath string
	Status int
	Mess string
}

func UploadBreakPoint(req *http.Request,key,savePath string,ch chan ResultChan){
	req.ParseForm()
	forms:=req.Form
	var rc=ResultChan{}
	//b,_:=json.Marshal(forms)
	formFile, header, err := getFileMessage(req,key)
	defer func(){
		if formFile!=nil{
			formFile.Close()
		}else{
			recover()
			rc.Mess="连接断开，停止上传！"
			rc.Status=0
			ch<-rc
			return
		}
	}()
	//todo check file message
	if !CheckErr(err) {
		rc.Mess="获取文件失败:"+key
		rc.Status=0
		ch<-rc
		return
	}
	b,err:=utils.PathNotExistsCreate(savePath)
	if !b {
		rc.Mess="创建目录失败"+savePath
		rc.Status=0
		ch<-rc
		return
	}
	thumbPath:=savePath+"thumb/"
	b,err =utils.PathNotExistsCreate(thumbPath)
	if !b {
		rc.Mess="创建缩略图目录失败:"+thumbPath
		rc.Status=0
		ch<-rc
		return
	}
	fileName:=header.Filename
	fileSaveName:=strconv.FormatInt(time.Now().Unix(),10)+"_"+fileName
	rc.FileName=fileSaveName
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
		b,_:=utils.IsPathExists(saveFileName)
		if b{
			rc.Mess="文件已经存在"
			log.Println(saveFileName+":文件已经存在")
			goto Nothing
		}
		b,err:=utils.IsPathExists(tempPath)
		if !b{
			err=os.MkdirAll(tempPath,os.ModePerm)
			if err!=nil{
				rc.Mess="Create directory error:%s"+ err.Error()
				log.Printf("Create directory error: %s\n", err)
				rc.Status=0;
				ch<-rc
				return
			}
		}
		var isOK bool
		//保存文件
		isOK,err=saveFile(saveFileName,formFile)
		if !isOK {
			rc.Mess = "保存失败" + err.Error()
			log.Println(saveFileName + ":保存失败")
			rc.Status = 0
			ch <- rc;
			return
		}
			rc.Status=1
		ch<-rc;
		return

	}else{
		//整体上传
		saveFileName=savePath+fileSaveName
		rc.FilePath=saveFileName
		var isOK bool
		//保存文件
		isOK,err=saveFile(saveFileName,formFile)
		if !isOK {
			log.Println(saveFileName + ":保存失败")
			rc.Status = 0
			ch <- rc;
			return
		}
		//创建缩略图
		thumbName:=strconv.FormatInt(time.Now().Unix(),10)+".jpg";
		rc.ThumbPath=thumbPath+thumbName
		err:=utils.CreateThumb(savePath+fileSaveName,savePath+"thumb/",thumbName,500,0)
		if err!=nil{
			rc.Status=0
			rc.Mess="创建缩略图失败！"
			ch<-rc;
			return
		}
		rc.Status=1
		rc.Mess=""
		ch<-rc;
		return
	}
	goto Ok
	Nothing:{
		log.Println("nothing")
		formFile.Close()
		rc.Status=1
		ch<-rc;
		return
	}
	Ok:{
		log.Println(" not BreakPoint")
	}
	//如果是分片上传判断是否上传完成,然后组合成完整文件
	if ok{
		if chunk==totalChunk-1{
			assembleFile(header.Filename,fileSaveName,savePath,totalChunk)
			thumbName:=strconv.FormatInt(time.Now().Unix(),10)+".jpg";
			thumbPath:=savePath+"thumb/"+thumbName
			rc.ThumbPath=thumbPath
			utils.CreateThumb(savePath+header.Filename,savePath+"thumb/",thumbName,500,0)
			rc.Status=1
			ch<-rc;
			return
		}
	}


}
//组装上传文件
func  assembleFile(fileName,fileSaveName,path string,totalCount int){

	fileFullPath:=path+fileSaveName
	fileNameTemp:=path+fileName
	b,_:=utils.IsPathExists(fileFullPath)
	if b{
		log.Println("文件已经存在:",fileFullPath)
		removeTempFile(fileNameTemp+"_temp",totalCount)
		return
	}

	file, err := os.OpenFile(fileFullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
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
