/*
	服务器启动 读取本地配置sql，缓存到内存中
*/
package utils

import (
	"os"
	"log"
)

var sd SqlDictionary

func Init(xmlPath string) {
	s, err := os.Getwd()
	if err != nil {
		log.Fatalln("init configs error")
	}
	spath := s + "/"+xmlPath
	var v Location
		AnalysisXml(spath, &v)

	for _, countryRegion := range v.CountryRegion {
		AnalysisXml(string(s+countryRegion.Path), &sd)
	}

}

type Location struct {
	CountryRegion []CountryRegion `xml:"datafile"`
}
type CountryRegion struct {
	Name string `xml:"label,attr"`
	Path string `xml:"src,attr"`
}

//获取sql字典数据
type SqlDictionary struct {
	SqlChildList []SqlChildList `xml:"items"`
}

func (s SqlDictionary) GetAllDictionary() SqlDictionary {

	return sd
}

type SqlChildList struct {
	Id   string    `xml:"id,attr"`
	Name string    `xml:"name,attr"`
	Sql  []SqlList `xml:"item"`
}
type SqlList struct {
	Id     string `xml:"id,attr"`
	Name   string `xml:"name,attr"`
	Optype string `xml:"optype,attr"`
	Sqltxt string `xml:"sql"`
}

type SqlMessage struct {
	TokenId string
	ItemId  string
	Optype  string
	Sqltxt  string
}

func (s *SqlList) GetSqlMess() *SqlList {

	return s
}

//获取单条sql
/*
	@parameter perentKey 父节点ID
	@parameter childKey 子节点ID
	return map[string]string
	mp["pid"]=""
	mp["cid"]= ""
	mp["sql"]=""
	mp["op"]=""
*/
func GetItemSql(perentKey, childKey string) *SqlMessage {

	//mp := make(map[string]string)

	mp := &SqlMessage{}
	for _, countryRegion := range sd.SqlChildList {
		if perentKey == countryRegion.Id {
			for _, sc := range countryRegion.Sql {
				if sc.Id == childKey {
					mp.TokenId = countryRegion.Id
					mp.ItemId = sc.Id
					mp.Sqltxt = sc.Sqltxt
					mp.Optype = sc.Optype
					//	fmt.Println(countryRegion.Id, sc.Id, sc.Sqltxt, sc.Optype)
				}
			}
		}
	}
	return mp
}
