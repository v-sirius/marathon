package webserver

import (
	"strconv"
	"net/url"
	"cache"
	"fmt"
	"log"
	"net/http"
	//	"reflect"
	//	"strconv"

	"time"
)

const longForm = "2006-01-02 15:04:05"

//好友个数---好友在cache中的读取
func FriendsCount(content map[string]interface{}, w http.ResponseWriter) {
	res := ResponseInfo{Code: "ret_friend", Is_success: false}

	ucode, exist := content["uid"]
	if !exist {
		res.Content = map[string]interface{}{"msg": "wrong uid", "data": nil}
		write2Client(w, res)
		return
	}

	user, err := cache.G_CacheData.UserSet.GetUserByCode(ucode.(string))
	if err != nil {
		res.Content = map[string]interface{}{"msg": "wrong uid", "data": nil}
		write2Client(w, res)
		return
	}

	friendSet := cache.G_CacheData.FriendSet.GetFriendsById(user.Id)
	friendNum := 0

	if friendSet != nil {
		friendNum = len(friendSet)
	}
	res.Content = map[string]interface{}{"msg": "return friends num", "data": map[string]interface{}{"data": friendNum}}
	res.Is_success = true
	write2Client(w, res)
	return
}

//上传好友---好友写入到cache中
func UploadFriend(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	res.Code = "ret_friend_upload"

	/*if _, exist := content["uid"]; !exist { //判断request body 内是否有id字段
		res.Is_success = false
		res.Content = map[string]interface{}{"time": time.Now().Format(longForm)}
		write2Client(w, res)
		return
	}
	if _, exist := content["friend"]; !exist { //判断request body 内是否有friend字段
		res.Is_success = false
		res.Content = map[string]interface{}{"time": time.Now().Format(longForm)}
		write2Client(w, res)
		return
	}*/

	// 用户code
	ucode := content["uid"].(string)
	usr, err := cache.G_CacheData.UserSet.GetUserByCode(ucode)
	if err != nil{
		log.Println("UploadFriend: wrong user")
		write2Client(w, res)
		return
	}
	// 手机md5
	friendTel := content["friends"].([]interface{})
	//log.Println("friendTel:", friendTel)
	// 获取留言信息
	msg := content["message"].(string)
	// 当前时间
	currentTime := time.Now().Format(longForm)

	data := make(url.Values)
	var destuid string = ""
	// 构造推送消息，包括请求者的nickname，tel，ucode以及时间等
	for _, v := range friendTel {
		//var pushInfo cache.PushInfo
		//pushInfo.Topic = v.(string)

		//tmpUser, err := cache.G_CacheData.UserSet.GetUserByCode(ucode)
		tmpUser, err := cache.G_CacheData.UserSet.GetUserByMd5(v.(string))
		if err != nil{
			log.Println("UploadFriend: wrong user for md5")	
			continue		
		}
		
		if destuid == ""{
			destuid = tmpUser.Code
		}else{
			destuid = destuid + "," + tmpUser.Code
		}
		
		/*log.Println("make friend, usrId:", usr.Id, "," , "tmpUsrId:", tmpUser.Id)
		cache.G_CacheData.FriendSet.FriendChan <- cache.Friend{UserId: usr.Id, FriendId: tmpUser.Id, MakeTime: time.Now().Format(longForm), Share: false}
		cache.G_CacheData.FriendSet.CreateFriend(usr.Id, tmpUser.Id, false)*/
		/*if err == nil {
			pushInfo.Content = map[string]interface{}{"type": "add_friend",
					"content": map[string]string{"uid": ucode,
					"nickname": tmpUser.NickName,
					"tel":      tmpUser.Mobile,
					"message":  msg,
					"time":     currentTime}}
			cache.G_PushInfo <- pushInfo
		}*/
	}
	/*{
“type”:”add_friend”,
“destuid”:“id1,id2...”,#从用户注册信息表里面找到被加好友的id，然后一次性的发过来，如果只是对一个人加好友，则这里的长度为一
“uid”:””,
“nickname”:””,
“message”:””,
“time”:”yyyy-MM-dd 24hh:mm:ss”
}*/
	if destuid != ""{
		data["type"] = []string{"add_friend"}
		data["destuid"] = []string{destuid}
		data["uid"] = []string{ucode}
		data["nickname"] = []string{usr.NickName}
		data["message"] = []string{msg}
		data["time"] = []string{currentTime}
		
		log.Println("****************data time of add friend*******",data["time"])
		
		cache.G_PushInfo <- data
	}
	
	
	res.Is_success = true
	res.Content = map[string]interface{}{"time": currentTime}
	write2Client(w, res)
	return
}

//接受好友请求
func AcceptFriendUpload(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	res.Code = "ret_accept_friend_upload"
	res.Is_success = false
	res.Content = map[string]interface{}{"time": time.Now().Format(longForm)}

	if _, exist := content["uid"]; !exist { //判断request body 内是否有id字段
		log.Println("AcceptFriendUpload: no uid")
		write2Client(w, res)
		return
	}
	if _, exist := content["fid"]; !exist { //判断request body 内是否有friend字段
		log.Println("AcceptFriendUpload: no friend")
		write2Client(w, res)
		return
	}
	if _, exist := content["accept"]; !exist { //判断request body 内是否有accept字段
		log.Println("AcceptFriendUpload: accept field")
		write2Client(w, res)
		return
	}

	uId, _ := content["uid"].(string)
	friendId, _ :=  content["fid"].(string)//*******content["friend"].(string) 
	accept_bool, _ := content["accept"].(string)
	
	log.Println("********************accept_bool of add friend**************",accept_bool)
//	var usr1 cache.User
var nickname_now string
	if accept_bool == "1" {
		usr1, err1 := cache.G_CacheData.UserSet.GetUserByCode(uId)
		nickname_now=usr1.NickName
		
		log.Println("*****************nickname_now***********",nickname_now)
		
		usr2, err2 := cache.G_CacheData.UserSet.GetUserByCode(friendId)
		
		log.Println("*********error of UserSet.GetUserByCode(uId)******",err1)
		log.Println("*********error of UserSet.GetUserByCode(friendId)******",err2)
		
		if err1 == nil && err2 == nil {
			cache.G_CacheData.FriendSet.CreateFriend(usr1.Id, usr2.Id, false)
			log.Println("********out creatfriend ******")
			res.Is_success = true
			cache.G_CacheData.FriendSet.FriendChan <- cache.Friend{UserId: usr1.Id, FriendId: usr2.Id, MakeTime: time.Now().Format(longForm), Share: false}
		//****zcy modified here 1****	cache.G_CacheData.FriendSet.CreateFriend(usr1.Id, usr2.Id, false)
			//res.Content = map[string]interface{}{"time": time.Now().Format(longForm)}
		}
	} else { //拒绝加为好友
		fmt.Println("//拒绝加为好友")
		res.Is_success = false
		//res.Content = map[string]interface{}{"time": time.Now().Format(longForm)}
	}
	
	/*{
“type”:”add_friend_echo”,
“destuid”:“id1”,#发起加好友请求用户的uid 
“echo”:YES/NO,
“uid”:””,
“nickname”:””,
“message”:””,
“time”:”yyyy-MM-dd 24hh:mm:ss”
}*/
	data := make(url.Values)
	data["type"] = []string{"add_friend_echo"}
	data["destuid"] = []string{friendId}
	if accept_bool == "1"{
		data["echo"]=[]string{"YES"}
	}else{
		data["echo"]=[]string{"NO"}
	}
	data["uid"] = []string{uId}
	data["nickname"] = []string{nickname_now}
	data["message"] = []string{""}
	data["time"] = []string{time.Now().Format(longForm)}
	
	//test
	log.Println("************data********",data)
	log.Println("************usr1---->friendId********",friendId)
	log.Println("************destuid********",data["destuid"])
	log.Println("************nickname_now********",nickname_now)
	log.Println("************data time of response add friend********",data["time"])
	
	cache.G_PushInfo <- data

	log.Println("--------------------response of add friend-----------")
	log.Println(res)
	write2Client(w, res)
}

// 约人请求
func InviteFriend(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	res.Code = "ret_invite_friend"
	res.Is_success = false
	currentTime := time.Now().Format(longForm)
	res.Content = map[string]interface{}{"time": currentTime}

	if _, exist := content["uid"]; !exist { //判断request body 内是否有id字段
		//res.Is_success = false
		log.Println("InviteFriend: no uid")
		write2Client(w, res)
		return
	}
	if _, exist := content["fid"]; !exist { //判断request body 内是否有friend字段
		log.Println("InviteFriend: no friend")
		write2Client(w, res)
		return
	}

	uId, _ := content["uid"].(string)
	friendId, _ := content["fid"].(string)
//	msg, msgExist := content["message"].(string)
//	lon, lonExist := content["lon"].(float64)
//	lat, latExist := content["lat"].(float64)
	lon:=content["lon"].(float64)
	lat:=content["lat"].(float64)
      msg:= content["message"].(string)
	
	log.Println("*******************lon of invite friend*******",lon)
	log.Println("*******************lat of invite friend*******",lat)
	log.Println("*******************msg of invite friend*******",msg)
	
//	string_lon:=content["lon"].(string)
//	string_lat:=content["lat"].(string)
//	float64_lon,_:= strconv.ParseInt(string_lon,10,64)
//	float64_lat,_:=strconv.ParseInt(string_lat,10,64)
//	log.Println("*******************lon of invite friend*******",float64_lon)
//	log.Println("*******************lat of invite friend*******",float64_lat)
	
	// 如果有错，则再服务器端报错
//	if msgExist != true || lonExist != true || latExist != true {
//		write2Client(w, res)
//		log.Println("InviteFriend: error in msg , lon or lat")
//		return
//	}
	
	usr, err := cache.G_CacheData.UserSet.GetUserByCode(uId)
	if err != nil{
		log.Println("InviteFriend: wrong uid")
		write2Client(w, res)
		return
	}
	
	/*{
“type”:”meeting”,
destuid:“id1”, 
“uid”:””,
“nickname”:””,
“message”:””,
“lon”:float,
“lat”:float,
“time”:”yyyy-MM-dd 24hh:mm:ss” 
}*/

	// 如果没错，则正常处理：构造消息发送至推送服务器
	data_meeting := make(url.Values)
	data_send_message := make(url.Values)

	// 获取被约者的手机md5
	if _, err := cache.G_CacheData.UserSet.GetUserByCode(friendId); err != nil {
		write2Client(w, res)
		return
	} else {
		log.Println("*********************invitor user nickname*************")
		log.Println("**********",usr.NickName)
		
		if lon==0 && lat==0{
			data_send_message["type"] = []string{"send_message"}
			data_send_message["destuid"] = []string{friendId}//data["fuid"] = []string{friendId}
			data_send_message["uid"] = []string{uId}
			data_send_message["nickname"] = []string{usr.NickName}
			data_send_message["message"] = []string{msg}
			data_send_message["time"] = []string{time.Now().Format(longForm)}
			
			log.Println("##################################")
			log.Println("#################",data_send_message)
			log.Println("##################################")
			
			cache.G_PushInfo <- data_send_message
		}else{
			data_meeting["type"] = []string{"meeting"}
			data_meeting["destuid"] = []string{friendId}//data["fuid"] = []string{friendId}
			data_meeting["uid"] = []string{uId}
			data_meeting["nickname"] = []string{usr.NickName}
			data_meeting["message"] = []string{msg}
			data_meeting["lon"] = []string{strconv.FormatFloat(lon, 'f', -1, 64)}
			data_meeting["lat"] = []string{strconv.FormatFloat(lat, 'f', -1, 64)}
			data_meeting["time"] = []string{time.Now().Format(longForm)}
			
			log.Println("##################################")
			log.Println("#################",data_meeting)
			log.Println("##################################")
			
			cache.G_PushInfo <- data_meeting
		}
		/*pushInfo.Topic = invitor.MobileMd5
		pushInfo.Content = map[string]interface{}{"uid": uId,
			"nickname": invitor.NickName,
			"tel":      invitor.Mobile,
			"message":  msg,
			"lon":      lon,
			"lat":      lat,
			"time":     currentTime}*/
		

		res.Is_success = true
		log.Println("^^^^")
		log.Println("^^^^")
		log.Println("^^",res.Is_success)
		log.Println("^^^^")
		log.Println("^^^^")
		write2Client(w, res)
		return
	}
	return
}

// 响应约人请求
func AcceptInvite(content map[string]interface{}, w http.ResponseWriter) {
	res := ResponseInfo{Code: "ret_accept_invite_friend", Is_success: false}
	currentTime := time.Now().Format(longForm)
	res.Content = map[string]interface{}{"time": currentTime}

	// 如果不存在，则返回错误
	uid, exist := content["uid"]
	if !exist { //判断request body 内是否有id字段
		//res.Is_success = false
		log.Println("不存在的用户ID")
		write2Client(w, res)
		return
	}
	
	fid, exist := content["fid"]
	if !exist { //判断request body 内是否有friend字段
		log.Println("不存在的朋友ID")
		write2Client(w, res)
		return
	}
	
	lat, exist := content["lat"]
	if !exist{
		log.Println("不存在lat")
		write2Client(w, res)
		return
	}
	lon, exist := content["lon"]
	if !exist{
		log.Println("不存在lon")
		write2Client(w, res)
		return
	}	
	
	user1, err := cache.G_CacheData.UserSet.GetUserByCode(uid.(string))
	if err != nil{
		log.Println("AcceptInvite: no usr1")
		write2Client(w, res)
		return
	}
	
	_, err = cache.G_CacheData.UserSet.GetUserByCode(fid.(string))
	if err != nil{
		log.Println("AcceptInvite: no usr2")
		write2Client(w, res)
		return
	}
	
	/*{
“type”:”meeting_echo”,
destuid:“id1”,
“echo”:YES/NO,
“uid”:””,
“nickname”:””,
“lon”:float,
“lat”:float,
“time”:”yyyy-MM-dd 24hh:mm:ss” 
}*/


	if accept_bool, exist := content["accept"]; !exist {
		log.Println("不存在accept字段")
		write2Client(w, res)
		return
	} else {
		data := make(url.Values)
		data["type"] = []string{"meeting_echo"}
		data["destuid"] = []string{fid.(string)}
		data["uid"]=[]string{uid.(string)}
		data["nickname"]=[]string{user1.NickName}
		data["lon"] = []string{strconv.FormatFloat(lon.(float64), 'f', -1, 64)}
		data["lat"] = []string{strconv.FormatFloat(lat.(float64), 'f', -1, 64)}
		data["time"] = []string{currentTime}
		if accept_bool.(string) == "1" {
			res.Is_success = true
			data["echo"]=[]string{"YES"}
		}else{
			data["echo"]=[]string{"NO"}
		}
		
		cache.G_PushInfo <- data
		write2Client(w, res)
		return
	}
	return
}

// 位置共享给好友(群组暂时未实现)
/*func SharePos2Friend(content map[string]interface{}, w http.ResponseWriter){
	var res ResponseInfo
	res = ResponseInfo{Code: "share_pos", Is_success:false}
	res.Content = map[string]interface{}{"time": time.Now().Format(longForm)}

	ucode, _ := content["uid"].(string)
	usr, err := cache.G_CacheData.UserSet.GetUserByCode(ucode)

	// 如果没有对应的用户code
	if err != nil{
		write2Client(w, res)
		return
	}

	// 遍历好友列表发送
	friends, exist != content["friends"].(map[string]interface{})
	if exist == nil{
		for _, v := range friends{

		}
	}
}*/

//查看好友信息返回数据
/*type FriendResp struct {
	FriendRespId int // 唯一标识

	FriendRespNickName string // 昵称

	//Img      string //image.Image // 图像

	FriendRespCategory    int // 用户类别
	FriendRespDescription string  // 描述
	FriendRespCreateTime  string  //time.Time // 创建时间
	FriendRespPos         cache.Position
}*/

//查看好友信息
func ViewFriendInfo(content map[string]interface{}, w http.ResponseWriter) {
	// 构造返回的对象
	currentTime := time.Now().Format(longForm)
	res := ResponseInfo{Code: "ret_friend_info",
		Is_success: false}
	res.Content = map[string]interface{}{"time": currentTime, "friends": nil}

	uidObj, exist := content["uid"]
	if !exist {
		write2Client(w, res)
		return
	}

	// 检查uid是否是字符串
	uid, err := uidObj.(string)
	if err != true {
		log.Println("ViewFriendInfo: uid is not string")
		write2Client(w, res)
		return
	}

	usr, errGet := cache.G_CacheData.UserSet.GetUserByCode(uid)
	if errGet != nil {
		log.Println("ViewFriendInfo: cannot get user by code")
		write2Client(w, res)
		return
	}

	friends := cache.G_CacheData.FriendSet.GetFriendsById(usr.Id)
	log.Println(usr.Id, "'s friends", friends)
	if friends == nil {
		log.Println("ViewFriendInfo: no friends")
		res.Is_success = true
		write2Client(w, res)
		return
	}

	friendsData := make([]map[string]interface{}, 0)
	// 逐个数据加入
	for _, v := range friends {
		usrTmp, err := cache.G_CacheData.UserSet.GetUserById(v.FriendId)
		if err != nil {
			log.Println("ViewFriendInfo: fault id")
			continue
		}

		tmpData := map[string]interface{}{"uid": usrTmp.Code,
			"nickname":    usrTmp.NickName,
			"u_type":      usrTmp.Category,
			"description": usrTmp.Description,
			"sex":usrTmp.Sexual,
			"share":v.Share,
		}
		friendsData = append(friendsData, tmpData)
	}

	res.Is_success = true
	res.Content = map[string]interface{}{"time": currentTime, "friends": friendsData}
	write2Client(w, res)
	return
}

// 好友位置共享请求
func SharePosWithFriend(content map[string]interface{}, w http.ResponseWriter){
	currentTime := time.Now().Format(longForm)
	res := ResponseInfo{Code: "ret_loc_share", Is_success: true, Content: map[string]interface{}{"time": currentTime}}
	
	ucode, flg := content["uid"]
	if !flg{
		res.Is_success = false
		log.Println("SharePosWithFriend: no user")
		write2Client(w, res)
		return
	}	
	
	tmpUser, err := cache.G_CacheData.UserSet.GetUserByCode(ucode.(string))
	if err != nil{
		res.Is_success = false
		log.Println("SharePosWithFriend: wrong user")
		write2Client(w, res)
		return
	}
	
	shares, flg := content["shares"]
	if !flg{
		log.Println("SharePosWithFriend: no friends")
		write2Client(w, res)
		return
	}
	
	allFriendsReq, ok := shares.([]interface{})
	if !ok{
		res.Is_success = false
		log.Println("SharePosWithFriend: wrong friend format1")
		write2Client(w, res)
		return
	}
	
	shareSlice := make([]string, 0)
	unshareSlice := make([]string, 0)
	
	for _, v := range allFriendsReq{
		fid, flg := v.(map[string]interface{})["fid"]
		if !flg{
			res.Is_success = false
			log.Println("SharePosWithFriend: wrong friend format2")
			write2Client(w, res)
			return
		}
		
		bShare, flg := v.(map[string]interface{})["share"]
		if !flg{
			res.Is_success = false
			log.Println("SharePosWithFriend: wrong friend format3")
			write2Client(w, res)
			return
		}
		
		uFriendCode := fid.(string)
		/*friendUser, err := cache.G_CacheData.UserSet.GetUserByCode(uFriendCode)
		if err != nil{
			res.Is_success = false
			log.Println("SharePosWithFriend: wrong friend format")
			write2Client(w, res)
			return
		}*/
		
		shared := bShare.(string)
		
		if shared == "1" {
			shareSlice = append(shareSlice, uFriendCode)
		}else if shared == "0"{
			unshareSlice = append(unshareSlice, uFriendCode)
		}else{
			
		}
		
		
		// todo: 如何推送
/*		{
“type”:”position_share”,
“destuid”:“id1,id2...”, 
“share”:YES/NO,
“uid”:””,
“nickname”:””,
“message”:””,
“time”:”yyyy-MM-dd 24hh:mm:ss”
}位置共享的请求与取消分开推送*/

	}
	
	log.Println("shared slice:", shareSlice)
	log.Println("unshared slice:", unshareSlice)
	
	data := make(url.Values)
	data["type"] = []string{"position_share"}
	
	var destuid string = ""
	for _, v := range shareSlice{
		if destuid == ""{
			destuid = destuid + v
		}else{
			destuid = destuid + "," + v
		}		
	}
	
	data["destuid"] = []string{destuid}
	data["share"] = []string{"YES"}
	data["uid"] = []string{tmpUser.Code}
	data["nickname"]=[]string{tmpUser.NickName}
	data["message"] = []string{""}
	data["time"] = []string{currentTime}
	
	data1 := make(url.Values)
	data1["type"] = []string{"position_share"}
	var destuid1 string
	for _, v := range unshareSlice{
		if destuid1 == ""{
			destuid1 = destuid1 + v
		}else{
			destuid1 = destuid1 + "," + v
		}
	}
	data1["destuid"] = []string{destuid1}
	data1["share"] = []string{"YES"}
	data1["uid"] = []string{tmpUser.Code}
	data1["nickname"]=[]string{tmpUser.NickName}
	data1["message"] = []string{""}
	data1["time"] = []string{currentTime}
	
	if destuid != ""{
		cache.G_PushInfo <- data
	}
	
	if destuid1 != "" {
		cache.G_PushInfo <- data1
	}
	
	write2Client(w, res)
	return
}

// 响应好友位置请求
func ResponsePosShare(content map[string]interface{}, w http.ResponseWriter){
	currentTime := time.Now().Format(longForm)
	var res ResponseInfo
	res = ResponseInfo{Code: "ret_res_loc_share", Is_success: false, Content: map[string]interface{}{"time": currentTime}}
	
	ucode, flg := content["uid"]
	if !flg {
		log.Println("ResponsePosShare: wrong user code")
		//res = ResponseInfo{Code: "ret_res_loc_share", Is_success: false, Content: map[string]interface{}{"time": currentTime}}
		write2Client(w, res)
		return
	}
	
	fcode, flg := content["fid"]
	if !flg {
		log.Println("ResponsePosShare: wrong friend code")
		//res = ResponseInfo{Code: "ret_res_loc_share", Is_success: false, Content: map[string]interface{}{"time": currentTime}}
		write2Client(w, res)
		return
	}
	
	accept, flg := content["accept"]
	if !flg {
		log.Println("ResponsePosShare: wrong request format")
		//res = ResponseInfo{Code: "ret_res_loc_share", Is_success: false, Content: map[string]interface{}{"time": currentTime}}
		write2Client(w, res)
		return
	}
	
//	var usr1 cache.User
	log.Println("******************entering accept 1*************")
	// 共享位置
	bool_accept,_:=strconv.ParseBool(accept.(string))
	log.Println("******************entering accept 1:bool_accept*************",bool_accept)
	if bool_accept==true{
		// 共享位置--修改内存， 修改数据库
		log.Println("******************entering accept 2*************")
		usr1, err := cache.G_CacheData.UserSet.GetUserByCode(ucode.(string))
		if err != nil{
			log.Println("ResponsePosShare: wrong user1")
		//	res = ResponseInfo{Code: "ret_res_loc_share", Is_success: false, Content: map[string]interface{}{"time": currentTime}}
			write2Client(w, res)
			return
		}
		usr2, err := cache.G_CacheData.UserSet.GetUserByCode(fcode.(string))
		if err != nil{
			log.Println("ResponsePosShare: wrong user2")
			//res = ResponseInfo{Code: "ret_res_loc_share", Is_success: false, Content: map[string]interface{}{"time": currentTime}}
			write2Client(w, res)
			return
		}
		//******zcy ,add res ****
		res = ResponseInfo{Code: "ret_res_loc_share", Is_success: true, Content: map[string]interface{}{"time": currentTime}}
	log.Println("******************res of accept*************",res)
	
		cache.G_CacheData.FriendSet.ShareFriendPos(usr1.Id, usr2.Id, true)
		
		//***test,zcy****
		log.Println("*************out ShareFriendPos******")
		
		sql := "update friend set shared=true where (uid=" + strconv.Itoa(usr1.Id) + " and friendid=" + strconv.Itoa(usr2.Id) + ") or (uid=" + strconv.Itoa(usr2.Id) + " and friendid=" + strconv.Itoa(usr1.Id) + ");"
		
		log.Println("*************sql******",sql)
		
		cache.G_CacheData.ChString <- sql
		write2Client(w, res)
		
			/*{
“type”:”position_share_echo”,
“destuid”:“id1”,#发起加好友请求用户的uid 
“echo”:YES/NO,
“uid”:””,
“nickname”:””,
“time”:”yyyy-MM-dd 24hh:mm:ss”
}*/

	data := make(url.Values)
	data["type"] = []string{"position_share_echo"}//*****{"add_friend_echo"}
	data["destuid"]= []string{fcode.(string)}
	data["uid"] = []string{ucode.(string)}
	
	//*********zcy modified here*********
	data["echo"] = []string{accept.(string)}//*****{strconv.FormatBool(accept.(bool))}
	data["nickname"] = []string{usr1.NickName}
	//data["message"] = []string{"result:" +accept.(string) }//****strconv.FormatBool(accept.(bool))}
	data["time"] = []string{currentTime}
	
	cache.G_PushInfo <- data
	}else{
		// 共享位置--修改内存， 修改数据库
		log.Println("******************entering accept 3*************")
		usr1, err := cache.G_CacheData.UserSet.GetUserByCode(ucode.(string))
		if err != nil{
			log.Println("ResponsePosShare: wrong user1")
		//	res = ResponseInfo{Code: "ret_res_loc_share", Is_success: false, Content: map[string]interface{}{"time": currentTime}}
			write2Client(w, res)
			return
		}
		usr2, err := cache.G_CacheData.UserSet.GetUserByCode(fcode.(string))
		if err != nil{
			log.Println("ResponsePosShare: wrong user2")
			//res = ResponseInfo{Code: "ret_res_loc_share", Is_success: false, Content: map[string]interface{}{"time": currentTime}}
			write2Client(w, res)
			return
		}
		//******zcy ,add res ****
		res = ResponseInfo{Code: "ret_res_loc_share", Is_success: true, Content: map[string]interface{}{"time": currentTime}}
	log.Println("******************res of accept*************",res)
	
		cache.G_CacheData.FriendSet.ShareFriendPos(usr1.Id, usr2.Id, false)
		
		//***test,zcy****
		log.Println("*************out ShareFriendPos******")
		
		sql := "update friend set shared=false where (uid=" + strconv.Itoa(usr1.Id) + " and friendid=" + strconv.Itoa(usr2.Id) + ") or (uid=" + strconv.Itoa(usr2.Id) + " and friendid=" + strconv.Itoa(usr1.Id) + ");"
		cache.G_CacheData.ChString <- sql
		
		write2Client(w, res)
	}
	log.Println("******************out accept*************")
	
	return
}

// 删除好友
func RemoveFriend(content map[string]interface{}, w http.ResponseWriter){
	currentTime := time.Now().Format(longForm)
	res := ResponseInfo{Code: "ret_remove_friend", Is_success: true, Content: map[string]interface{}{"time": currentTime}}
	
	ucode, flg := content["uid"]
	if !flg {
		log.Println("RemoveFriend: wrong uid")
		res.Is_success = false
		write2Client(w, res)
		return
	}
	
	fcodes, flg := content["fids"]
	if !flg {
		log.Println("RemoveFriend: wrong fids")
		res.Is_success = false
		write2Client(w, res)
		return
	}
	
	user, err := cache.G_CacheData.UserSet.GetUserByCode(ucode.(string))
	if err != nil{
		log.Println("RemoveFriend: wrong user")
		res.Is_success = false
		write2Client(w, res)
		return
	}
	
	destuid := ""
	for _, v := range fcodes.([]interface{}){
		tmpUser, err := cache.G_CacheData.UserSet.GetUserByCode(v.(string))
		if err != nil{
			log.Println("RemoveFriend: wrong friend")
			res.Is_success = false
			write2Client(w, res)
			return
		}
		if destuid == ""{
			destuid = v.(string)
		}else{
			destuid = destuid + "," + v.(string)
		}
		//destuid = destuid + "," + v
		cache.G_CacheData.FriendSet.DeleteFriends(user.Id, tmpUser.Id)
		updateString := "delete from friend where uid='" + strconv.Itoa(user.Id) + "'and friendid='" + strconv.Itoa(tmpUser.Id) + "';"
		cache.G_CacheData.ChString <- updateString
	}	
	
	/*{
“type”:”delete_friend”,
“destuid”:“id1,id2,id3...” 
“uid”:””,
“nickname”:””,
“time”:”yyyy-MM-dd 24hh:mm:ss”
}*/
	data := make(url.Values)
	data["type"] = []string{"delete_friend"}
	data["destuid"] = []string{destuid}
	data["uid"] = []string{ucode.(string)}
	data["nickname"] = []string{user.NickName}
	data["time"] = []string{currentTime}
	
	cache.G_PushInfo <- data
	
	write2Client(w, res)
	return
}