package cache

import (
	"db"
	"fmt"
	"strconv"
)

// 标识数据库缓冲区的大小
const BUFFERSIZE4DB = 1000

// 写入数据库的缓冲区
type Buffer4Db struct {
	/*	idx   int            // 标识写哪个缓存
		buf1  []UserPosition // 标识第一个缓冲区
		buf2  []UserPosition // 标识第二个缓冲区
		chIdx chan int       // 标识写入数据库的哪一个缓冲区*/
	chBuf chan UserPosition
}

// 初始化
func initBuffer4Db() *Buffer4Db {
	pRet := new(Buffer4Db)

	//pRet.idx = 1
	pRet.chBuf = make(chan UserPosition, BUFFERSIZE4DB)
	//pRet.buf2 = make([]UserPosition, 0)

	//pRet.chIdx = make(chan int, 1)

	return pRet
}

// 写入数据库
func (pBuf Buffer4Db) WriteDb() {
	fmt.Println("---------------位置开始写入数据库：")
	//pBuf := G_CacheData.Buffer4Db
	for userposition := range pBuf.chBuf {
		sql := "insert into userlocation(id,uid,ulongitude,ulatitude,createtime) values " +
			"('" + strconv.FormatInt(userposition.Id, 10) + "','" + strconv.Itoa(userposition.UserId) + "'," + strconv.FormatFloat(userposition.TmSpace.Pos.Longitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(userposition.TmSpace.Pos.Latitude, 'f', -1, 64) + ",'" + userposition.TmSpace.CreateTime + "')"
		if !db.G_db.Insert2Table(sql) {
			fmt.Println("userposition表插入不成功！")
		}
	}
	/*for {
		idx := <-pBuf.chBuf
		// 如果idx＝＝1，则将buf1写入数据库
		// 如果idx＝＝2，则将buf2写入数据库
		if idx == 1 {
			fmt.Println("---------------位置 buf1写入数据库：")
			for _, userposition := range pBuf.buf1 {

				fmt.Println("-----------buf1-userposition", userposition)
				sql := "insert into userlocation(id,uid,ulongitude,ulatitude,createtime) values " +
					"('" + strconv.FormatInt(userposition.Id, 10) + "','" + strconv.Itoa(userposition.UserId) + "'," + strconv.FormatFloat(userposition.TmSpace.Pos.Longitude, 'f', -1, 64) +
					"," + strconv.FormatFloat(userposition.TmSpace.Pos.Latitude, 'f', -1, 64) + ",'" + userposition.TmSpace.CreateTime + "')"
				if !db.G_db.Insert2Table(sql) {
					fmt.Println("userposition表插入不成功！")
				}
			}
			//pBuf.buf1 = nil
		} else {
			fmt.Println("---------------位置 buf2写入数据库：")
			for _, userposition := range pBuf.buf2 {

				fmt.Println("-----------buf2-userposition", userposition)
				sql := "insert into userlocation(id,uid,ulongitude,ulatitude,createtime) values " +
					"('" + strconv.FormatInt(userposition.Id, 10) + "','" + strconv.Itoa(userposition.UserId) + "'," + strconv.FormatFloat(userposition.TmSpace.Pos.Longitude, 'f', -1, 64) +
					"," + strconv.FormatFloat(userposition.TmSpace.Pos.Latitude, 'f', -1, 64) + ",'" + userposition.TmSpace.CreateTime + "')"
				if !db.G_db.Insert2Table(sql) {
					fmt.Println("userposition表插入不成功！")
				}
			}
			//pBuf.buf2 = nil
		}
	}*/

}
