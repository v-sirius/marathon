package webserver

import (
	"log"
	"cache"
	"fmt"
	"net/http"
)

//标注信息
func PositionLabled(content map[string]interface{}, w http.ResponseWriter) {
	res := ResponseInfo{Code:"ret_identify", Is_success: false, Content: nil}
	
	category, exist  := content["type"]
	if !exist{
		log.Println("PositionLabeled: no type field!")
		write2Client(w, res)
		return
	}
	
	categoryValue, flg := category.(int64)
	if !flg{
		log.Println("PostionLabeled: wrong type's value")
		write2Client(w, res)
		return
	}
	
	allLabel, err := cache.G_CacheData.LabelSet.GetPosByCategory(int(categoryValue))
	if err != nil{
		log.Println("PositionLabeled: wrong type!")
		write2Client(w, res)
		return
	}
	
	// 返回所有信息
	allDatas := make([]interface{}, 0)
	for _, v := range allLabel{
		tmpData := map[string]interface{}{"name": v.Name,
										"lon": v.Pos.Longitude,
										"lat": v.Pos.Latitude,
										"description": v.Description,
										"type": v.Category,
										"value": v.Value}
		allDatas = append(allDatas, tmpData)
	}
	res.Content = map[string]interface{}{"msg": "return specified information",
										"data": allDatas}
	res.Is_success = true
	
	write2Client(w, res)
	return
}

//运动员信息
func AthletesInfo(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	res.Code = "ret_athe_info"

	// 返回运动员信息，数组形式的运动员列表
	if time2check, exist := content["time"]; exist {
		if athletes, err := cache.G_CacheData.GetAllAthletes(time2check.(string)); err != nil { //
			fmt.Println(athletes)
			res.Is_success = false
			res.Content = map[string]interface{}{"msg": "no athletes data", "data": nil}

			fmt.Println(res)
			write2Client(w, res)
			return
		} else {
			res.Is_success = true
			res.Content = map[string]interface{}{"msg": "successful", "data": athletes}

			fmt.Println(res)
			write2Client(w, res)
			return
		}
	} else {
		res.Is_success = false
		res.Content = map[string]interface{}{"msg": "not exiest time in reaquest", "data": nil}

		fmt.Println(res)
		write2Client(w, res)
		return
	}
}

//拉拉队信息
func CheersInfo(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	res.Code = "ret_cheer"

	// 返回啦啦队信息，数组形式的啦啦队列表
	if time2check, exist := content["time"]; exist {

		if cheerleaders, err := cache.G_CacheData.GetAllCheerleaders(time2check.(string)); err != nil {
			res.Is_success = false
			res.Content = map[string]interface{}{"msg": "no cheerleader data", "data": nil}

			fmt.Println(res)
			write2Client(w, res)
			return
		} else {
			res.Is_success = true
			res.Content = map[string]interface{}{"msg": "successful", "data": cheerleaders}

			fmt.Println(res)
			write2Client(w, res)
			return
		}
	} else {
		res.Is_success = false
		res.Content = map[string]interface{}{"is_success": false, "msg": "not exiest time in reaquest", "data": nil}

		fmt.Println(res)
		write2Client(w, res)
		return
	}
}

//合作伙伴信息
func PartnerInfo(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	res.Code = "ret_partner"

	// 返回合作伙伴信息，数组形式的合作伙伴信息列表
	if time2check, exist := content["time"]; exist {
		if partners, err := cache.G_CacheData.GetAllPartners(time2check.(string)); err != nil { //
			res.Is_success = false
			res.Content = map[string]interface{}{"msg": "no partners data", "data": nil}

			fmt.Println(res)
			write2Client(w, res)
			return
		} else {
			res.Is_success = true
			res.Content = map[string]interface{}{"msg": "successful", "data": partners}

			fmt.Println(res)
			write2Client(w, res)
			return
		}
	} else {
		res.Is_success = false
		res.Content = map[string]interface{}{"msg": "not exiest time in reaquest", "data": nil}

		fmt.Println(res)
		write2Client(w, res)
		return
	}
}

//日程信息
func CalendarInfo(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	res.Code = "ret_calendar"

	// 返回日程信息，数组形式的日程信息列表
	if time2Check, exist := content["time"]; exist {
		if agendas, err := cache.G_CacheData.GetAllAgendas(time2Check.(string)); err != nil {
			res.Is_success = false
			res.Content = map[string]interface{}{"msg": "no agendas data", "data": nil}

			fmt.Println(res)
			write2Client(w, res)
			return
		} else {
			res.Is_success = true
			res.Content = map[string]interface{}{"msg": "successful", "data": agendas}

			fmt.Println(res)
			write2Client(w, res)
			return
		}
	} else {
		res.Is_success = false
		res.Content = map[string]interface{}{"msg": "not exiest time in reaquest", "data": nil}

		fmt.Println(res)
		write2Client(w, res)
		return
	}
}

//应急路线推送
func EmPathPush(content map[string]interface{}, w http.ResponseWriter) {
	
}
