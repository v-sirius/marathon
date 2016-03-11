package cache

import (
	"db"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// 用户上线记录
type OnlineRecord struct {
	Id        int64
	UserId    int
	Pos       Position
	LoginTime time.Time
}

// 所有的用户上线记录
type OnlineRecords struct {
	allOnlineRecord []OnlineRecord
	id2OnlineRecord map[int64]*OnlineRecord
	//	bFlg            chan bool
	OnlineRecordChan chan OnlineRecord

	rwLock sync.RWMutex // 读写锁
}

// 初始化函数
func InitOnlineRecords() *OnlineRecords {
	pRet := new(OnlineRecords)
	pRet = &OnlineRecords{allOnlineRecord: make([]OnlineRecord, 0), id2OnlineRecord: make(map[int64]*OnlineRecord, 0), OnlineRecordChan: make(chan OnlineRecord, 0)}

	return pRet
}

// 写入数据
func (pSet *OnlineRecords) WriteOnlineRecord2Cache(onlinerecord OnlineRecord) {
	pSet.rwLock.Lock()
	defer func() {
		pSet.rwLock.Unlock()
	}()

	//idx := len(pSet.allOnlineRecord)

	pSet.allOnlineRecord = append(pSet.allOnlineRecord, onlinerecord)

	pSet.id2OnlineRecord[onlinerecord.Id] = &onlinerecord
}

// 读取数据
func (pSet *OnlineRecords) GetOnlineRecordById(id int64) (OnlineRecord, error) {
	pSet.rwLock.RLock()
	defer func() {
		pSet.rwLock.RUnlock()
	}()

	if idx, flg := pSet.id2OnlineRecord[id]; flg {
		return *idx, nil
	}
	return OnlineRecord{Id: -1}, errors.New("not exist related OnlineRecord data")
}

// 加载数据
func (pOnline *OnlineRecords) LoadData(db *db.DbOperation) {

	rows, err := db.Db.Query("select * from onlinerecord")
	var onlinerecord OnlineRecord
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err = rows.Scan(&onlinerecord.Id, &onlinerecord.UserId, &onlinerecord.Pos, &onlinerecord.LoginTime)
		if err != nil {
			panic(err)
		}

	}
	pOnline.WriteOnlineRecord2Cache(onlinerecord)

}

//channel写入数据库
func (pLabelSet *LabelSet) Write2DbOnlineRecord() {
	var onlinerecordchan = G_CacheData.OnlineRecords.OnlineRecordChan
	for onlinerecord := range onlinerecordchan {
		var sql = "insert into onlinerecord(id,uid,loginlongitude,loginlatitude,logintime)" +
			" values ('" + strconv.FormatInt(onlinerecord.Id, 36) + "','" +
			strconv.Itoa(onlinerecord.UserId) + "','" +
			strconv.FormatFloat(onlinerecord.Pos.Longitude, 'e', 10, 64) + "','" +
			strconv.FormatFloat(onlinerecord.Pos.Latitude, 'e', 10, 64) + "','" +
			onlinerecord.LoginTime.String() + "')"

		if !db.G_db.Insert2Table(sql) {
			fmt.Println("onlinerecord表插入不成功！")
		}

	}

}
