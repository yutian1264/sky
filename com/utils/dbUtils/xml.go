package utils

import (
	"os"
	"io/ioutil"
	"encoding/xml"
)

/*
	解析本地xml
*/
func AnalysisXml(path string ,obj interface{})error{
	f, err := os.Open(path)
	if err != nil {
		return err;
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err;
	}

	err = xml.Unmarshal(data, obj)
	if err != nil {
		return err;
	}
	return nil;
}