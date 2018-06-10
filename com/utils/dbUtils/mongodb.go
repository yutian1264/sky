package utils

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"time"
	"fmt"
)

var (
	mgoSession *mgo.Session
	dialInfo  *mgo.DialInfo
)

func MongodbInit(dbName,addr,port,userName,pwd,source string,poolLimit int){
	dialInfo = &mgo.DialInfo{
		Addrs:     []string{addr+":"+port},
		Direct:    false,
		Timeout:   time.Second * 1,
		Database:  dbName,
		Source:    source,
		Username:  userName,
		Password:  pwd,
		PoolLimit: poolLimit, // Session.SetPoolLimit
	}
}

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getSession() *mgo.Session {

	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			panic(err) //直接终止程序运行
			er:=recover()
			fmt.Println(er)
		}
	}
	//最大连接池默认为4096

	return mgoSession.Clone()
}
//公共方法，获取collection对象
func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dialInfo.Database).C(collection)
	return s(c)
}


/**
 * 获取一条记录通过objectid
 */
func GetElementById(collectionName ,id string) (result interface{},err error) {

	objid := bson.ObjectIdHex(id)

	query := func(c *mgo.Collection) error {
		return c.FindId(objid).One(&result)
	}
	err=witchCollection(collectionName, query)

	return
}

func GetItemByAny(collationName string ,query bson.M)(result interface{},err error){

	exop := func(c *mgo.Collection) error {
		return c.Find(query).One(&result)
	}
	err = witchCollection(collationName, exop)
	return
}


func UpdateItem(query bson.M, change bson.M) (string,error){
	exop := func(c *mgo.Collection) error {
		return c.Update(query, change)
	}
	err := witchCollection("person", exop)
	if err != nil {
		return "true",nil
	}
	return "false",err
}

/**
 * 执行查询，此方法可拆分做为公共方法
 * [SearchPerson description]
 * @param {[type]} collectionName string [description]
 * @param {[type]} query          bson.M [description]
 * @param {[type]} sort           bson.M [description]
 * @param {[type]} fields         bson.M [description]
 * @param {[type]} skip           int    [description]
 * @param {[type]} limit          int)   (results      []interface{}, err error [description]
 */
func GetListByPage(collectionName string, query bson.M, sort string, fields bson.M, skip int, limit int) (results []interface{}, err error) {
	exop := func(c *mgo.Collection) error {
		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	err = witchCollection(collectionName, exop)
	return
}

func FindAll(collectionName string)(results []interface{}, err error){
	exop := func(c *mgo.Collection) error {
		return c.Find(nil).All(&results)
	}
	err = witchCollection(collectionName, exop)
	return
}


