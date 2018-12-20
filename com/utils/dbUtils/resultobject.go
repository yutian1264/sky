/*
@Time : 2018/10/29 15:21 
@Author : sky
@Description:
@File : resultobject
@Software: GoLand
*/
package utils

import "encoding/json"

const (
	Status0  ="未知错误"
	Status1  ="正常数据"
	Status2  ="用户错误"
	Status3  ="暂无权限"

)

type ResultObject struct {
	Status int   //0:未知错误 1:正常数据 2:用户错误 3:暂无权限
	Msg string
	Data interface{}

}
func (t ResultObject)GetResuleString(status int,msg string,data interface{})string{

	r:=ResultObject{status,msg,data}
	b,e:=json.Marshal(r)
	if e!=nil{
		return string(`{"Status":0,"Msg":`+e.Error()+`,"Data":null}`)
	}
	return string(b)
}

