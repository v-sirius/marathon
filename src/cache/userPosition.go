package cache

import (
	"log"
	"errors"
	"strconv"
	"sync"
	"db"
	"fmt"
)

// 时空结构 －－ 通常我们查询时指定特定的时间和空间
type TimeSpace struct {
	CreateTime string
	Pos        Position
}

type TimeSpaces []*TimeSpace

// 用户位置的实时位置
type UserPosition struct {
	Id      int64
	UserId  int
	TmSpace TimeSpace
}

type Positions []Position
type UserPositions []UserPosition

// 用于存储用户id的切片
type UserIds []int

// 用于存储用户位置唯一标识id的切片
type UserPositionIds []int64

// 存储所有用户的位置信息
type RealTimePositionSet struct {
	allPositions    []UserPosition          // 所有的位置记录
	userId2Postions map[int]UserPositionIds // 记录单个人的运行轨迹
	//position2Users	map[TimeSpace]UserIds 	// 记录在某个时刻某个位置周边的人
	time2Positions map[string]UserPositionIds // 记录某个时刻所有的位置记录

	idx chan int // 写标识

	RwLock sync.RWMutex //读写锁，存储和查询分开，查询－读，存储－写
}

// 初始化
func initRealTimePositionSet() *RealTimePositionSet {
	pRet := new(RealTimePositionSet)

	pRet = &RealTimePositionSet{allPositions: make([]UserPosition, 0), userId2Postions: make(map[int]UserPositionIds),
		time2Positions: make(map[string]UserPositionIds), idx: make(chan int), RwLock: sync.RWMutex{}}

	return pRet
}


//******************zcy,add here**************
func (rt *RealTimePositionSet)LoadData(pDb *db.DbOperation){
	//select real time position from db
	//var userPos UserPosition
	fmt.Println("begining:load RealTimePositionSet........")
	selectstr := "select id,uid,ulongitude,ulatitude,createtime from userlocation"
	rows := pDb.Find(selectstr)
	
	size:= len(rt.allPositions)
	fmt.Println("************size*********")
	fmt.Println("*******",size)
	for rows.Next() {
		var userPos UserPosition
		err := rows.Scan(&userPos.Id, &userPos.UserId,&userPos.TmSpace.Pos.Longitude,&userPos.TmSpace.Pos.Latitude,&userPos.TmSpace.CreateTime)

		if err != nil {
			panic("error in:scanning the table userlocation")
		}

		fmt.Println("*************Load userlocation data to cache*********")
		fmt.Println(userPos)

		//size=len(rt.allPositions)
		rt.WriteUserLocationFromDb2Cache(&size,userPos)
		size = len(rt.allPositions)
		log.Println("size: ", size)
	}
}

// ***********zcy,从db写入数据到cache*********
func (rt *RealTimePositionSet) WriteUserLocationFromDb2Cache(size *int,usrPos UserPosition) {
	//fmt.Println("************write to cache***********")
	//fmt.Println(usrPos)

	rt.RwLock.Lock()
	defer func() {
		rt.RwLock.Unlock()
	}()
	
     // fmt.Println("************size*********")
	//fmt.Println(size)
	string_size:=strconv.Itoa(*size)
	int64_size, _ := strconv.ParseInt(string_size, 10, 64)
	//usr.Id = size
	fmt.Println("************int64_size*********")
	fmt.Println(int64_size)
	rt.allPositions = append(rt.allPositions, usrPos)

	rt.userId2Postions[usrPos.UserId]=append(rt.userId2Postions[usrPos.UserId],int64_size)
//rt.time2Position[usrPos.TmSpace.CreateTime]=append(rt.time2Position[usrPos.TmSpace.CreateTime],size)
	
}

// 获取朋友的轨迹
func (rt RealTimePositionSet) GetTrackByFriendId(id int) UserPositionIds {
	rt.RwLock.RLock()
	defer func() {
		rt.RwLock.RUnlock()
	}()

	ret, flg := rt.userId2Postions[id]
	if !flg {
		return nil
	}

	return ret
}

// 获取某个人的最新记录
func (rt RealTimePositionSet) GetNewestPosByFrendId(id int) (UserPosition, error) {
	rt.RwLock.RLock()
	defer func() {
		rt.RwLock.RUnlock()
	}()

	track, flg := rt.userId2Postions[id]
	log.Println("********************id of userId2Postions********",id)
	//log.Println("********************track of userId2Postions********",track)
	//log.Println("********************track[len(track)-1 of userId2Postions********",track[len(track)-1])
	log.Println("********************flag of track********",flg)
	if !flg {
		
		return UserPosition{Id: -1, UserId: -1, TmSpace: TimeSpace{CreateTime: "", Pos: Position{Latitude: -1, Longitude: -1}}}, errors.New("wrong position")
	}
//	if track==nil{
//		return UserPosition{Id: -1, UserId: id, TmSpace: TimeSpace{CreateTime: "", Pos: Position{Latitude: 0, Longitude: 0}}}, errors.New("no  position")
	
//	}

	return rt.allPositions[int(track[len(track)-1])], nil
}

//查看用户的轨迹
func (pSet *RealTimePositionSet) CheckUserTrace(uid string) (string, Positions) {
	uid_int, _ := strconv.Atoi(uid)
	userpositionIds := pSet.userId2Postions[uid_int]
	var p Positions
	var nt string //最新时间
	for _, userpositionId := range userpositionIds {
		t := pSet.allPositions[userpositionId].TmSpace.CreateTime //位置创建的时间
		nt = t                                                    //格式化后的时间
		position := pSet.allPositions[userpositionId].TmSpace.Pos //位置
		p = append(p, position)
	}
	return nt, p

}

type FriendPosition struct {
	UserId        int
	TimeFriendPos TimeSpace
}

//读取好友位置数据
func (realTimePosSet *RealTimePositionSet) GetFriendPos(uid int) FriendPosition {
	//var friendPositions []FriendPosition
	var friend_pos FriendPosition
	realTimePosSet.RwLock.RLock()
	defer func() {
		realTimePosSet.RwLock.RUnlock()
	}()

	f_positions := realTimePosSet.userId2Postions[uid] //UserPositionIds
	f_positionsLen := len(f_positions)                 //UserPositionIds的长度
	f_idIndex := f_positions[f_positionsLen-1]
	friend_pos.UserId = uid
	friend_pos.TimeFriendPos.CreateTime = realTimePosSet.allPositions[f_idIndex].TmSpace.CreateTime
	friend_pos.TimeFriendPos.Pos.Longitude = realTimePosSet.allPositions[f_idIndex].TmSpace.Pos.Longitude
	friend_pos.TimeFriendPos.Pos.Latitude = realTimePosSet.allPositions[f_idIndex].TmSpace.Pos.Latitude

	//friendPositions = append(friendPositions, friend_pos)

	return friend_pos
}
