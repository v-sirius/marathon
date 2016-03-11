package cache

import (
	"strconv"
	"db"
	"log"

	//	"errors"
	//	"fmt"
	//	"strconv"
	"sync"
	//	"time"
)

const FRIEND_CHAN_NUM = 100

// 好友信息
type Friend struct {
	Id       int    // 数据id
	UserId   int    // 用户id
	FriendId int    // 朋友id
	MakeTime string //time.Time //交好友时间
	Share 	 bool   // 是否共享
}

// 考虑是否共享
type FriendShare struct{
	FriendId int
	Share	 bool
}

// 包括用户的共享信息
type FriendIds []FriendShare

// 内存中好友格式信息，好友是否是相互的，因此，下面这个结构不再需要
/*type FriendMem struct {
	UserId    int
	FriendIds []int
}*/

// 好友集合
type FriendSet struct {
	allFriend     []Friend
	userId2Friend map[int]FriendIds
	//bFlg          chan bool
	rwLock     sync.RWMutex // 读写锁
	FriendChan chan Friend
}

// 初始化函数
func InitFriendSet() *FriendSet {
	pRet := new(FriendSet)
	pRet = &FriendSet{allFriend: make([]Friend, 0), userId2Friend: make(map[int]FriendIds, 0), FriendChan: make(chan Friend, FRIEND_CHAN_NUM)}

	return pRet
}

// 根据用户的id获取所有的用户
func (friendSet FriendSet) GetFriendsById(id int) FriendIds {
	friendSet.rwLock.Lock()
	defer func(){
		friendSet.rwLock.Unlock()
	}()
	
	if _, flg := friendSet.userId2Friend[id]; !flg {
		return nil
	}
	return friendSet.userId2Friend[id]
}

// 写入数据to Cache
/*func (friendSet *FriendSet) WriteFriend2Cache(friendMem FriendMem) {
	friendSet.rwLock.Lock()
	defer func() {
		friendSet.rwLock.Unlock()
	}()

	//friendSet.allFriend = append(friendSet.allFriend, friendMem)

	friendSet.userId2Friend[friendMem.UserId] = &friendMem
	//friendSet.userId2Friend[friendMem.UserId] = friendSet.userId2Friend
}*/

// 加载数据
func (pFriendSet *FriendSet) LoadData(pDb *db.DbOperation) {
	//查找数据库中好友表，加载到cache的pFriendSet.allFriend中
	selectstr1 := "select * from friend"
	rows1 := pDb.Find(selectstr1)
	for rows1.Next() {
		var friend Friend
		err := rows1.Scan(&friend.Id, &friend.UserId, &friend.FriendId, &friend.MakeTime, &friend.Share)
		if err != nil {
			panic("error in:scanning the table friend")
		}
		pFriendSet.allFriend = append(pFriendSet.allFriend, friend)
	}
	
	log.Println("all friends:", pFriendSet.allFriend)

	// 根据数据库中的数据构建好友关系
	for _, v := range pFriendSet.allFriend {
		userId1, userId2, bShare := v.UserId, v.FriendId, v.Share
		log.Println("userid1:", userId1, "userid2:", userId2)
		if _, flg := pFriendSet.userId2Friend[userId1]; !flg {
			pFriendSet.userId2Friend[userId1] = make(FriendIds, 0)
		}
		if _, flg := pFriendSet.userId2Friend[userId2]; !flg {
			pFriendSet.userId2Friend[userId2] = make(FriendIds, 0)
		}

		if !isInSlice(pFriendSet.userId2Friend[userId1], userId2) {
			pFriendSet.userId2Friend[userId1] = append(pFriendSet.userId2Friend[userId1], FriendShare{userId2, bShare})
		}

		if !isInSlice(pFriendSet.userId2Friend[userId2], userId1) {
			pFriendSet.userId2Friend[userId2] = append(pFriendSet.userId2Friend[userId2], FriendShare{userId1, bShare})
		}
		
		log.Println("friends:", pFriendSet.userId2Friend[userId1], "friend2:", pFriendSet.userId2Friend[userId2])
	}
}

// 判断一个整数是否存在于一个slice中
func isInSlice(slice FriendIds, iVar int) bool {
	for _, v := range slice {
		if v.FriendId == iVar {
			return true
		}
	}
	return false
}

// 建立好友关系
func (pFriendSet *FriendSet) CreateFriend(userId1, userId2 int, bShare bool) {
	log.Println("********entering creatfriend ******")
	
	pFriendSet.rwLock.Lock()
	defer func(){
		pFriendSet.rwLock.Unlock()
	}()
	
	if _, flg := pFriendSet.userId2Friend[userId1]; !flg {
		pFriendSet.userId2Friend[userId1] = make(FriendIds, 0)
	}
	if _, flg := pFriendSet.userId2Friend[userId2]; !flg {
		pFriendSet.userId2Friend[userId2] = make(FriendIds, 0)
	}

	if !isInSlice(pFriendSet.userId2Friend[userId1], userId2) {
		pFriendSet.userId2Friend[userId1] = append(pFriendSet.userId2Friend[userId1], FriendShare{userId2, bShare})
	}

	if !isInSlice(pFriendSet.userId2Friend[userId2], userId1) {
		pFriendSet.userId2Friend[userId2] = append(pFriendSet.userId2Friend[userId2], FriendShare{userId1, bShare})
	}
}

// 解除好友关系
func (pFriendSet *FriendSet) DeleteFriends(usrId1, usrId2 int){
	var i, j int = -1, -1
	for k, v := range pFriendSet.allFriend{
		if (v.UserId == usrId1 && v.FriendId == usrId2) || (v.UserId == usrId2 && v.FriendId == usrId1){
			if i == -1 {
				i = k
				continue
			}
			if j == -1 {
				j = k
				break
			}
		}
	}
	
	if j != -1 {
		pFriendSet.allFriend = append(pFriendSet.allFriend[:j], pFriendSet.allFriend[j+1:]...)
	}
	if i != -1 {
		pFriendSet.allFriend = append(pFriendSet.allFriend[:i], pFriendSet.allFriend[:i+1]...)
	}
	
	i = -1
	for k, v := range pFriendSet.userId2Friend[usrId1]{
		if v.FriendId == usrId2{
			i = k
		}
	}
	if i != -1{
		pFriendSet.userId2Friend[usrId1] = append(pFriendSet.userId2Friend[usrId1][:i], pFriendSet.userId2Friend[usrId1][i+1:]...)
	}
	
	i = -1 
	for k, v := range pFriendSet.userId2Friend[usrId2]{
		if v.FriendId == usrId1{
			i = k
		}
	}
	if i != -1{
		pFriendSet.userId2Friend[usrId2] = append(pFriendSet.userId2Friend[usrId2][:i], pFriendSet.userId2Friend[usrId2][i+1:]...)
	}
}

// 共享/取消好友位置
func (pFriendSet *FriendSet) ShareFriendPos(userId1, userId2 int, bShare bool){
	pFriendSet.rwLock.Lock()
	defer func(){
		pFriendSet.rwLock.Unlock()
	}()
	
	//***test zcy****
	log.Println("*************entering ShareFriendPos******")
	for _, v := range pFriendSet.allFriend{
		if v.UserId == userId1 && v.FriendId == userId2{
			v.Share = bShare
			break
		}
	}
	
	for _, v := range pFriendSet.userId2Friend[userId1]{
		if v.FriendId == userId2{
			v.Share = bShare
			break
		}
	}
	
	//********zcy modified here (add as follows)*****
	for _, v := range pFriendSet.userId2Friend[userId2]{
		if v.FriendId == userId1{
			v.Share = bShare
			break
		}
	}
	
}

//注册 从管道写入数据库
func (pFriend *FriendSet) Write2DbFriend() {
	log.Println("------------entering write2db-----------")

	for friend := range pFriend.FriendChan {

		//fmt.Println("------------------user in db-------------------")
		//fmt.Println(user)

		var sql = "insert into friend(uid, friendid, addftime, shared) values (" + strconv.Itoa(friend.UserId) + "," + strconv.Itoa(friend.FriendId) + ",'"  + friend.MakeTime +"',"+ strconv.FormatBool(friend.Share) + ");"

		if !db.G_db.Insert2Table(sql) {
			log.Println("user表插入不成功！")
		} else {
			log.Println("entering weite2db OK ! ! !")
		}
	}
}

//写入数据库
/*func Write2DbFriend(friendmem FriendMem) {
	const longForm = "2006-01-02 15:04:05"
	t := time.Now().Format(longForm)
	i := friendmem.UserId
	for j, friend := range friendmem.FriendIds {
		var sql = "insert into friend(uid,friendid,addftime)" +
			" values (" + strconv.Itoa(i) + "," + strconv.Itoa(friend) + ",'" + t + "')"

		if !db.G_db.Insert2Table(sql) {
			fmt.Println(strconv.Itoa(i) + strconv.Itoa(j) + "friend表插入不成功！")
		}
	}
}*/

//获取好友个数
/*func (pFriendSet *FriendSet) GetFriendCount(uid string) int {
	//	var fmen *FriendMem
	var idIndex, _ = strconv.Atoi(uid)
	//fmen = pFriendSet.userId2Friend[idIndex]
	if _, ok := pFriendSet.userId2Friend[idIndex]; ok {
		count := len(pFriendSet.userId2Friend[idIndex].FriendIds)
		fmt.Println(idIndex, pFriendSet.userId2Friend[idIndex].FriendIds)
		return count
	} else {
		return 0
	}
}*/

//获取所有好友id
/*func (pFriendSet *FriendSet) GetFriendIds(uid int) []int {
	var fmen []int
	//var idIndex, _ = strconv.Atoi(uid)

	if _, ok := pFriendSet.userId2Friend[uid]; ok {
		fmen = pFriendSet.userId2Friend[uid].FriendIds
		return fmen
	} else {
		return nil
	}
}*/

//上传好友，加好友
func (pFriendSet *FriendSet) Request4AddFriend(id string, mobile []interface{}) bool {
	/*var idIndex, _ = strconv.Atoi(id)
	//var tel, _ = strconv.Atoi(mobile.(string))

	fmt.Println("--------------entering request4AddFriend......---------")

	//判断用户好友列表中是否已经存在所上传好友，
	//好友列表若无则查看所有用户里是否存在上传的好友手机号，若有则返回消息提示已经为好友
	//好友列表中无好友信息，还需要查看所有用户里是否存在该手机号用户
	//存在则服务器推送给好友手机号加好友请求，若不存在，则返回用户不存在消息
	if v, ok := pFriendSet.userId2Friend[idIndex]; ok {

		for _, m := range mobile { //遍历电话列表
			//valueOfField, _ := strconv.Atoi(value.Field(i).String())
			mo := m.(string)

			if uididx, exist := G_CacheData.UserSet.mobile2User[mo]; exist {
				uid := G_CacheData.UserSet.allUser[uididx].Id
				for j := 0; j < len(v.FriendIds); j++ {
					if uid == v.FriendIds[j] {
						//TODO推送消息:提示已经为好友
						fmt.Println("已经是好友了。")
						break
					} else {
						//TODO推送消息到friend ：加好友请求
						fmt.Println("//TODO推送消息到friend ：加好友请求")
						break
					}

				}

			} else {
				fmt.Println("此用户未使用app")
				continue
			}

		}

		return true
	} else {

		//TODO推送消息到friend ：加好友请求
		fmt.Println("好友表内无，//TODO推送消息到friend ：加好友请求")

		return false
	}*/
	return true

}

//同意加为好友后，对加好友请求的响应
func (pFriendSet *FriendSet) Response4AddFriend() {
	/*var friendchan = G_CacheData.FriendSet.FriendChan
	for friendmem := range friendchan {
		Write2DbFriend(friendmem) //写入数据库
		uid := friendmem.UserId
		for _, friendid := range friendmem.FriendIds { //写入cache

			if friendM, ok := pFriendSet.userId2Friend[uid]; ok {
				friendM.FriendIds = append(friendM.FriendIds, friendid)
				pFriendSet.WriteFriend2Cache(*friendM)
			}

			if friendM, ok := pFriendSet.userId2Friend[friendid]; ok {
				friendM.FriendIds = append(friendM.FriendIds, uid)
				pFriendSet.WriteFriend2Cache(*friendM)
			}
		}
	}
	*/
}
