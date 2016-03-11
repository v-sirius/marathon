package webserver

import (
	"log"
	"time"
	"cache"
//	"fmt"
	"net/http"
	//	"strconv"
)

//const longForm = "2006-01-02 15:04:05"

// 位置上传
func UpdatePosition(content map[string]interface{}, w http.ResponseWriter) {
	currentTime := time.Now().Format(longForm)
	res := ResponseInfo{Code: "ret_up_pos", Is_success: false, Content:map[string]interface{}{"time": currentTime}}
	
	ucode, exist := content["uid"]
	if !exist{
		log.Println("UpdatePosition: no ucode!")
		write2Client(w, res)
	}
	
	pos, exist := content["pos"].([]interface{})
	if !exist{
		log.Println("UpdatePostion: no pos info")
		write2Client(w, res)
	}
	
	user, err := cache.G_CacheData.UserSet.GetUserByCode(ucode.(string))
	if err != nil{
		log.Println("UpdatePosition: no user")
		write2Client(w, res)
	}
	
	var tm string
	tm = time.Now().Format(longForm)
	// 逐个将位置信息发往管道
	for _, v := range pos {
		log.Println("**********************v of up_pos",v)
		var tmpPos cache.UserPosition

		tmpV, _ := v.(map[string]interface{})
		
		log.Println("*****************tmpV*********",tmpV)
		lon:= tmpV["lon"].(float64)
		log.Println("*****************lon*****",lon)
		lat:= tmpV["lat"].(float64)
		tm:= tmpV["time"].(string)
		log.Println("************tm********",tm)

		tmpPos.UserId = user.Id
		tmpPos.TmSpace.Pos = cache.Position{Longitude: lon, Latitude: lat}
		log.Println("********************tmpPos.Pos********",tmpPos.TmSpace.Pos)

		tmpPos.TmSpace.CreateTime = tm

		log.Println("****************************************tmpPos******",tmpPos)
		cache.G_CacheData.ChPos <- tmpPos
		
	}
	
	res.Is_success = true
	res.Content = map[string]interface{}{"time": tm}	
	
	write2Client(w, res)
	log.Println("***********UpdateLocation-- res:****", res)
	return
}

//用户轨迹
func UserPath(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	res.Code = "ret_trace"

	//读取前端信息
	code := content["uid"].(string)

	res.Is_success = true
	res_t, res_c := cache.G_CacheData.RealTimePositionSet.CheckUserTrace(code)

	res.Content["time"] = res_t
	res.Content["data"] = res_c
	write2Client(w, res)
	return

}

//查看某个人的轨迹
func CheckPath(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	res.Code = "ret_someone_trace"
	//读取前端信息
	code := content["uid"].(string)

	res.Is_success = true
	res_t, res_c := cache.G_CacheData.RealTimePositionSet.CheckUserTrace(code)

	res.Content["time"] = res_t
	res.Content["data"] = res_c
	write2Client(w, res)
	return
}

//查看好友位置
func ViewFriendPos(content map[string]interface{}, w http.ResponseWriter) {
	res := ResponseInfo{Code: "ret_friend_pos", Is_success: false, Content: nil}
	
	ucode, exist := content["uid"]
	if !exist{
		log.Println("ViewFriendPos: no uid info!")
		write2Client(w, res)
		return
	}
	
	usr, err := cache.G_CacheData.UserSet.GetUserByCode(ucode.(string))
	if err != nil{
		log.Println("ViewFriendPos: wrong uid info")
		write2Client(w, res)
		return
	}
	
	friends := cache.G_CacheData.FriendSet.GetFriendsById(usr.Id)
	if friends == nil{
		log.Println("ViewFriendPos: no friend!")
		write2Client(w, res)
		return
	}
	
	dataValue := make([]map[string]interface{}, 0)
	var tm string
	// 遍历所有朋友，获取相关信息
	for _, v := range friends{
		log.Println("************************loop1*****************")
		log.Println("************************loop1*****************")
		var uid string
		var lon, lat float64
		
		tmpUser, err := cache.G_CacheData.UserSet.GetUserById(v.FriendId)
		if err == nil{
			uid = tmpUser.Code
		}
		
		track, err := cache.G_CacheData.RealTimePositionSet.GetNewestPosByFrendId(v.FriendId)
		if err != nil{
			log.Println("ViewFriendPos: no friend position data, friend id:", v)
			continue
		}else{
			if v.Share == true{
				log.Println("********************sharer or not******",v.Share)
				lon = track.TmSpace.Pos.Longitude
				lat = track.TmSpace.Pos.Latitude
				tm = track.TmSpace.CreateTime
			
				dataValue = append(dataValue, map[string]interface{}{"uid": uid, "lon": lon, "lat": lat})
			}			
		}		
	}
	res.Is_success = true
	res.Content = map[string]interface{}{"time": tm, "data": dataValue}
	
	log.Println("****res****res****res****",res)
	write2Client(w,res)
	return
	}
