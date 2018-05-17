/**
	MySQL

 */
package utils

import (
"database/sql"
"fmt"
"github.com/astaxie/beego"
_ "github.com/go-sql-driver/mysql"

"strconv"
"strings"

	"qgtemp/tutils"
	"qgtemp/conf/contant"
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

type SqlResult struct {
	Msg   string
	Error int
	Data  []map[string]string
}

//获取全局变量db
var db *sql.DB

//初始化数据库
func InitMySQLDB(uname, pwd, address, dbName string) {
	var err error
	db, err = sql.Open("mysql", uname+":"+pwd+"@tcp("+address+")/"+dbName+"?charset=utf8")

	if err != nil {
		beego.Error("openDB error:", err)
		return
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
}
func GetResultPool(m map[string]string) string {
	mp := getSql(m)

	param := strings.Replace(string(m["param"]), "[", "", -1)
	param = strings.Replace(param, "]", "", -1)
	//解析参数为数组
	paramArr := strings.Split(param, ",")

	_sql := marryParam(mp.Sqltxt, paramArr)
	if m == nil {
		//未找到当前sql
		return "{\"Msg\": \"cannot find item sql" + string(m["key"]) + "\",\"Error\": 0,\"Data\":[]}"
	}
	var result string
	switch mp.Optype {
	case "1": //查询
		resultJson := Query(_sql)
		result = tutils.ResultJson(resultJson)
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
	_sql := marryParam(mp.Sqltxt, param.Params)
	var result interface{}
	switch mp.Optype {
	case "1": //查询
		resultJson := Query(_sql)
		if resultType {
			result = resultJson
		} else {
			result = tutils.ResultJson2(resultJson)
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
func Query(sql string) []map[string]interface{} {

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
func getSql(m map[string]string) *contant.SqlMessage {
	mp := &contant.SqlMessage{}
	if m != nil {
		mp = contant.GetItemSql(m["key"], m["key1"])
	}
	return mp
}

//

// func pool(w http.ResponseWriter, r *http.Request) {
//	rows, err := db.Query("SELECT * FROM ns_user limit 1")
//	defer rows.Close()
//	checkErr(err)
//
//	columns, _ := rows.Columns()
//	scanArgs := make([]interface{}, len(columns))
//	values := make([]interface{}, len(columns))
//	for j := range values {
//		scanArgs[j] = &values[j]
//	}
//
//	record := make(map[string]string)
//	for rows.Next() {
//		//将行数据保存到record字典
//		err = rows.Scan(scanArgs...)
//		for i, col := range values {
//			if col != nil {
//				record[columns[i]] = string(col.([]byte))
//			}
//		}
//	}
//
//	fmt.Println(record)
//	fmt.Fprintln(w, "finishtttt")
//}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

