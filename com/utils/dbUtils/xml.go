package utils

import (
	"os"
	"io/ioutil"
	"encoding/xml"
)

/*
	解析本地xml
*/
func AnalysisXml(path string ,obj interface{}){
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(data, obj)
	if err != nil {
		panic(err)
	}
}