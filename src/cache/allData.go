package cache

import (
	"strings"
	"db"
	"fmt"
	"log"
	"strconv"
)

// 位置管道的长度
const POSITION_CHANNEL_LEN = 10000

// cache中所有的数据结构
// main程序启动后，须初始化这些结构
type AllCacheData struct {
	AgendaSet
	AthleteSet
	Buffer4CachePosition
	Buffer4Db
	CheerleaderSet
	FriendSet
	LabelSet
	OnlineRecords
	PartnerSet
	UserSet
	RealTimePositionSet
	BehaviourSet

	ChPos    chan UserPosition
	ChPos4Db chan UserPosition
	
	ChString chan string

	Stop bool
}

// 所有数据结构
var G_CacheData *AllCacheData

// 初始化函数－－初始化所有的缓存数据结构
func InitAllCacheData() *AllCacheData {
	pRet := new(AllCacheData)

	pRet.AgendaSet = *InitAgendaSet()
	pRet.AthleteSet = *InitAthleteSet()
	pRet.Buffer4CachePosition = *InitBuffer4CachePosition()
	pRet.Buffer4Db = *initBuffer4Db()
	pRet.CheerleaderSet = *InitCheerleaderSet()
	pRet.FriendSet = *InitFriendSet()
	pRet.LabelSet = *InitLabelSet()
	pRet.OnlineRecords = *InitOnlineRecords()
	pRet.PartnerSet = *InitPartnerSet()
	pRet.UserSet = *InitUserSet()
	pRet.RealTimePositionSet = *initRealTimePositionSet()
	pRet.BehaviourSet = *InitBehaviourSet()

	pRet.ChPos = make(chan UserPosition, POSITION_CHANNEL_LEN)
	pRet.ChPos4Db = make(chan UserPosition, POSITION_CHANNEL_LEN)
	pRet.ChString = make(chan string, 1000)

	pRet.Stop = false

	return pRet
}

// 加载所有用户数据
func LoadAllData(pRet *AllCacheData, pDb *db.DbOperation) {
	pRet.AgendaSet.LoadData(pDb)
	pRet.AthleteSet.LoadData(pDb)
	pRet.UserSet.LoadData(pDb)
	pRet.PartnerSet.LoadData(pDb)
	pRet.CheerleaderSet.LoadData(pDb)
	pRet.FriendSet.LoadData(pDb)
	pRet.LabelSet.LoadLabelSet(pDb)
	
	//***********load RealTime Position,zcy*******
	pRet.RealTimePositionSet.LoadData(pDb)
}

// 开启goroute
func (pAll *AllCacheData) GoCache() {
	fmt.Println("entering GoChe!")
	//log.Println(pAll.RealTimePositionSet)
	log.Println(pAll.ChPos)

	//go pAll.Buffer4CachePosition.WriteData2Buffer(&pAll.RealTimePositionSet, pAll.ChPos)
	//go pAll.Buffer4CachePosition.Write2CacheAndDbBuffer(&pAll.Buffer4Db, &pAll.RealTimePositionSet)
	go pAll.Write2CacheAndDbBuffer()
	go pAll.WritePosDb()
	//go WriteDb() //位置写入数据库
	log.Println("out Gocache")
}

// 从缓存写入数据－－同时写入cache和db缓存
func (pAll *AllCacheData) Write2CacheAndDbBuffer(/*chPos <-chan UserPosition*/) {
	for elm := range pAll.ChPos {
		log.Println("**************************************************")
		log.Println("*****************Write2CacheAndDbBuffer***************")
		log.Println("**************************************************")
		// 写cache缓存
		pAll.RealTimePositionSet.RwLock.Lock()


//**************zcy,modified here,zhu shi xia mian xin xi **************
		//pRealSet := pAll.RealTimePositionSet
		//size := len(pRealSet.allPositions)
/*
		elm.Id = int64(len(pRealSet.allPositions))
		pRealSet.allPositions = //append(pRealSet.allPositions, elm)

		// 添加用户id到轨迹的映射
		if _, exist := pRealSet.userId2Postions[elm.UserId]; !exist {
			pRealSet.userId2Postions[elm.UserId] = make(UserPositionIds, 0)
		}
		pRealSet.userId2Postions[elm.UserId] = append(pRealSet.userId2Postions[elm.UserId], elm.Id)

		// 添加时间到位置的映射
		if _, exist := pRealSet.time2Positions[elm.TmSpace.CreateTime]; !exist {
			pRealSet.time2Positions[elm.TmSpace.CreateTime] = make(UserPositionIds, 0)
		}
		pRealSet.time2Positions[elm.TmSpace.CreateTime] = append(pRealSet.time2Positions[elm.TmSpace.CreateTime], elm.Id)
*/


	elm.Id = int64(len(pAll.RealTimePositionSet.allPositions))
		pAll.RealTimePositionSet.allPositions = append(pAll.RealTimePositionSet.allPositions, elm)

		// 添加用户id到轨迹的映射
		if _, exist := pAll.RealTimePositionSet.userId2Postions[elm.UserId]; !exist {
			pAll.RealTimePositionSet.userId2Postions[elm.UserId] = make(UserPositionIds, 0)
		}
		pAll.RealTimePositionSet.userId2Postions[elm.UserId] = append(pAll.RealTimePositionSet.userId2Postions[elm.UserId], elm.Id)

		// 添加时间到位置的映射
		if _, exist := pAll.RealTimePositionSet.time2Positions[elm.TmSpace.CreateTime]; !exist {
			pAll.RealTimePositionSet.time2Positions[elm.TmSpace.CreateTime] = make(UserPositionIds, 0)
		}
		pAll.RealTimePositionSet.time2Positions[elm.TmSpace.CreateTime] = append(pAll.RealTimePositionSet.time2Positions[elm.TmSpace.CreateTime], elm.Id)


log.Println("************out write cacheandbuffer********")

		pAll.RealTimePositionSet.RwLock.Unlock()
log.Println("************elm********",elm)
		// 写db缓存
		pAll.ChPos4Db <- elm
	}
}

// 写入数据库
func (pAll *AllCacheData)WritePosDb(){
	log.Println("**************************************************")
		log.Println("*****************WritePosDb*********************")
		log.Println("**************************************************")
		for userposition := range pAll.ChPos4Db {
		sql := "insert into userlocation(uid,ulongitude,ulatitude,createtime) values " +
			"('"  + strconv.Itoa(userposition.UserId) + "'," + strconv.FormatFloat(userposition.TmSpace.Pos.Longitude, 'f', -1, 64) +
			"," + strconv.FormatFloat(userposition.TmSpace.Pos.Latitude, 'f', -1, 64) + ",'" + userposition.TmSpace.CreateTime + "')"
		if !db.G_db.Insert2Table(sql) {
			fmt.Println("userposition表插入不成功！")
		}
	}
}

// 处理数据库更新或删除
func (pAll *AllCacheData) UpdateDb(){
	for v := range pAll.ChString{
		if db.G_db.Insert2Table(v) == false {
			log.Println("update error: ", v)
		}
	}
}

//接口：sql语句中特殊字符处理--encode
func SqlEncode(msg string)string{
	log.Println("----------input of sqlEncode():",msg)
	
	msgAfterEncode:=strings.Replace(msg,"'","&#39",-1)
	
	msgAfterEncode=strings.Replace(msgAfterEncode,"\"","&#34",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,"=","&#61",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,"-","&#45",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,";","&#59",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,"exec","ＥＸＥＣ",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,"or","ＯＲ",-1)
	msgAfterEncode=strings.Replace(msgAfterEncode,"and","ＡＮＤ",-1)
	
	log.Println("----------ouput after encode:",msgAfterEncode)
	return msgAfterEncode
}

//接口：sql语句中特殊字符处理--decode
func SqlDecode(msg string)string{
	log.Println("----------input of sqlDecode():",msg)
	
	msgAfterDecode:=strings.Replace(msg,"&#39","'",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"&#34","\"",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"&#61","=",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"&#45","-",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"&#59",";",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"ＥＸＥＣ","exec",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"ＯＲ","or",-1)
	msgAfterDecode=strings.Replace(msgAfterDecode,"ＡＮＤ","and",-1)
	
	log.Println("----------output after decode:",msgAfterDecode)
	return msgAfterDecode
}

