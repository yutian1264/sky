package timeUtils

import (
	"time"
)
const base_format = "2006-01-02 15:04:05"
//获取当前时间
func GetCurrentTime() string {
	return time.Now().Format(base_format)
}

//时间对比
func TimeChange(t1, t2 string) bool {
	a, _ := time.Parse(base_format, t1)
	b, _ := time.Parse(base_format, t2)
	return a.After(b.Add(300 * time.Second))
}

//时间格式化成字符串
func Time2String(t time.Time) string {
	return t.Format(base_format)
}
//时间格式化成字符串
func Str2Time(t string)  time.Time {
	time, _ := time.ParseInLocation(base_format, t,time.Local)
	return time
}
