package main

import (
	"os"

	//"encoding/json"
	//	"database/sql"
	//"fmt"
	"cache"
	"db"
	_ "github.com/Go-SQL-Driver/MySQL"
	//	"io"
	"log"
	//"logfile"
	"net/http"
	"webserver"
)

var fileout *os.File

func main() {
	//创建日志文件

	//	fileout = logfile.InitLogfile("maraapplog" + logfile.GetNowtime() + ".txt")
	//	io.WriteString(fileout, logfile.GetNowtime()+"   starting the server..........\n")

	// 创建数据库
	//	io.WriteString(fileout, logfile.GetNowtime()+"   创建数据库..........\n")
	db.G_db = db.InitDbOperation(`/test?charset=utf8`)
	db.G_db.Open()
	defer db.G_db.Close()
	db.G_db.CreateTable()

	// 初始化位置管理topic
	cache.G_PositionTopic.InitPositionTopic()
	
	// 创建缓存
	//	io.WriteString(fileout, logfile.GetNowtime()+"   begining: create the cache of the server..........\n")

	cache.G_CacheData = cache.InitAllCacheData()

	// 推送服务
	go cache.PushToServer(cache.G_PushInfo)
	go cache.G_CacheData.UpdateDb()
	
	// 写入位置信息到数据库
	go cache.G_CacheData.Buffer4Db.WriteDb()	
	go cache.G_CacheData.FriendSet.Write2DbFriend()

	go cache.G_CacheData.UserSet.Write2DbUser() //用户写入数据库
	//	io.WriteString(fileout, logfile.GetNowtime()+"   ending: create the cache of the server successfully..........\n")
	//	io.WriteString(fileout, logfile.GetNowtime()+"   begining: load all datas to the cache..........\n")

	//go cache.G_CacheData.FriendSet.Response4AddFriend() //好友加入数据库

	//加载所有数据
	cache.LoadAllData(cache.G_CacheData, db.G_db)
	//	io.WriteString(fileout, logfile.GetNowtime()+"   ending: load all data successfully..........\n")
	cache.G_CacheData.GoCache()

	//	io.WriteString(fileout, logfile.GetNowtime()+"   begining: set handler..........\n")

	http.HandleFunc("/mobile", webserver.HandleRequest)

	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		//		io.WriteString(fileout, logfile.GetNowtime()+"   ListenAndServe:"+err.Error())
		log.Fatal("ListenAndServe: ", err)
	}
	//	io.WriteString(fileout, logfile.GetNowtime()+"   success!..........\n")

}
