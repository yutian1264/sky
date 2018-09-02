/**
	MySQL

 */
package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"encoding/json"
)

type DataResult struct {
	error  string
	errMsg string
	Data   []map[string]string
}

type DBParam struct {
	TokenId string
	ItemId  string
	Params  []string
}

func (d *DataResult) GetDataResult() {

}

type ObjectResult struct {
	Msg   string
	Status int
	Data  map[string]interface{}
}

//获取全局变量db
var db *sql.DB

//初始化数据库
func InitMySQLDB(uname, pwd, address, dbName string) error {
	var err error
	db, err = sql.Open("mysql", uname+":"+pwd+"@tcp("+address+")/"+dbName+"?charset=utf8")

	if err != nil {

		return err
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
	return nil
}
/*
	获取完整sql 语句
*/
func GetResultPool(m map[string]string) string {
	mp := getSql(m)

	param := strings.Replace(string(m["param"]), "[", "", -1)
	param = strings.Replace(param, "]", "", -1)
	//解析参数为数组
	paramArr := strings.Split(param, ",")

	_sql:= marryParam(mp.Sqltxt, paramArr)
	if m == nil {
		//未找到当前sql
		return "{\"Msg\": \"cannot find item sql" + string(m["key"]) + "\",\"Error\": 0,\"Data\":[]}"
	}
	var result string
	switch mp.Optype {
	case "1": //查询
		resultJson,_:= Query(_sql)
		b,_:=json.Marshal(resultJson)
		fmt.Println(b)
	case "2": //删除 、修改
		result = Sql_update(_sql, true)
	case "3": //添加
		result = Add(_sql)
	}
	return result
}

func Common(tokenId, itemId string, param []string, isResultObject bool) interface{} {

	p := &DBParam{}
	p.TokenId = tokenId
	p.ItemId = itemId
	p.Params = param
	result := DBCommon(p, isResultObject)
	return result
}

/**
param 参数 {tokenid,itemid,[]param}
resultType true []map[string]interface{}
resultType false  string
return interface
*/
func DBCommon(param *DBParam, resultType bool) interface{} {

	m := make(map[string]string)
	m["key"] = param.TokenId
	m["key1"] = param.ItemId
	mp := getSql(m)
	_sql:= marryParam(mp.Sqltxt, param.Params)
	var result interface{}
	switch mp.Optype {
	case "1": //查询
		resultJson,_ := Query(_sql)
		if resultType {
			result = resultJson
		} else {
			result = ""
		}
		fmt.Println("query")
	case "2": //删除 、修改
		result = Sql_update(_sql, true)
	case "3": //添加
		result = Add(_sql)
	}
	return result

}

//删除数据返回删除记录数
func Sql_update(_sql string, isReturnJson bool) string {

	stmt, err := db.Prepare(_sql)
	checkErr(err)
	res, err := stmt.Exec()
	checkErr(err)
	num, err := res.LastInsertId()
	checkErr(err)

	if isReturnJson {
		return "{\"Msg\": \"\",\"Error\": 1,\"Data\":" + strconv.FormatInt(num, 10) + "}"
	}
	return strconv.Itoa(int(num))
}

//添加信息返回添加信息Id
func Add(_sql string) string {
	stmt, err := db.Prepare(_sql)
	checkErr(err)
	res, err := stmt.Exec()
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	return "{\"Msg\": \"\",\"Error\": 1,\"Data\":" + strconv.FormatInt(id, 10) + "}"
}

//查询数据返回json对象
func Query(sql string)([]map[string]interface{},error ){

	rows, err := db.Query(sql)
	defer rows.Close()
	checkErr(err)

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	//var num int=0
	result := []map[string]interface{}{}
	for rows.Next() {
		record := make(map[string]interface{})
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		result = append(result, record)
	}

	return result,err
}

func QueryByPage(tname,condition string,pageSize,pageCount int )map[string]interface{}{

	result:=make(map[string]interface{})
	sql:="select count(*) count from "+tname
	if len(condition)>0{
		sql+=" where "+condition
	}
	total,_:=Query(sql)
	sql1:= "select * from "+tname +" where 1=1 "
	if len(condition)>0{
		sql1+=" and "+condition
	}
	sql1+=" limit "+strconv.Itoa(pageSize*pageCount) +","+strconv.Itoa(pageCount)
	data,_:=Query(sql1)
	result["totalCount"]=total
	result["data"]=data
	return result
}

//匹配sql中参数 返回完整sql
func marryParam(sqlStr string, param []string) string {

	if len(param) > 0 {
		for i := 0; i < len(param); i++ {

			sqlStr = strings.Replace(sqlStr, "[#"+strconv.Itoa(i+1)+"#]", param[i], -1)
		}
	}


	return sqlStr
}

/*
	从sql字典中获取sql语句，并匹配参数
*/
func getSql(m map[string]string) *SqlMessage {
	mp := &SqlMessage{}
	if m != nil {
		mp = GetItemSql(m["key"], m["key1"])
	}
	return mp
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

/**
	获取完整sql
 */
func GetFullSql(key,key1 string,param []string)string{
	mp := &SqlMessage{}
	if (key != "" && key1!="") {
		mp = GetItemSql(key, key1)
	}
	return marryParam(mp.Sqltxt,param)
}

