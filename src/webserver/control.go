// 这个文件中只定义控制通用结构，具体的处理放在各个处理文件中
// 对于请求处理，对需求应该作响应调整：
// 1.所有的请求目录地址都是一样，没有必要每个申请一个地址
// 2.所有的返回信息中都必须包含一个是否成功处理字段，客户端方便获知结果，否则，客户端判断很麻烦

package webserver

import (
	"fmt"
	//	"cache"
	"io/ioutil"
	"net/http"

	//	"database/sql"
	"encoding/json"
)

// 接收request的基本结构
type RequestInfo struct {
	Category string                 `json:"type"`    // 类别
	Content  map[string]interface{} `json:"content"` // 具体内容
}

// 返回response的基本结构
type ResponseInfo struct {
	Code       string                 `json:"code"`       // 返回类型
	Is_success bool                   `json:"is_success"` //
	Content    map[string]interface{} `json:"content"`    // 返回的具体内容
}

// 处理接收到的所有请求
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)

		fmt.Println("request body:", string(body))
		//fmt.Println("")

		if err != nil {
			// todo:需添加日志
			panic("error in reading request's body when handling request")
		}

		// 定义一个接收的结构变量
		var requestBody RequestInfo
		err = json.Unmarshal(body, &requestBody)

		fmt.Println("------------error------:", err)
		fmt.Println("---------------------json Unmashal----------------")
		fmt.Println(requestBody)

		if err != nil {
			// todo:需要添加日志
			panic("format error in request body")
		}

		// 根据type种类分别进行处理
		switch requestBody.Category {

		//1. 用户名检测
		case "check_user":
			fmt.Println("entering check user name............")
			CheckUser(requestBody.Content, w)

			// 2.注册
		case "register":
			Register(requestBody.Content, w)

			// 用户信息完善
		case "user_update":
			UpdateUser(requestBody.Content, w)

			//3.登录
		case "login":
			UserLogin(requestBody.Content, w)

			//4.上传位置
		case "up_pos":
			UpdatePosition(requestBody.Content, w)

		//5.标注信息
		case "identify":
			PositionLabled(requestBody.Content, w)

			//6.运动员信息
		case "athe_info":
			AthletesInfo(requestBody.Content, w)

			//7.啦啦队信息
		case "cheer":
			CheersInfo(requestBody.Content, w)

			//8.合作伙伴信息
		case "partner":
			PartnerInfo(requestBody.Content, w)

			//9.日程信息
		case "calendar":
			CalendarInfo(requestBody.Content, w)

			//10.好友个数
		case "friend":
			FriendsCount(requestBody.Content, w)

			//12.查看自己轨迹
		case "trace":
			UserPath(requestBody.Content, w)

			//13.查看某人轨迹
		case "someone_trace":
			CheckPath(requestBody.Content, w)

			//14.查看好友位置
		case "friend_pos":
			ViewFriendPos(requestBody.Content, w)

			//15.应急路线推送 -----未完成
		case "emergency":
			EmPathPush(requestBody.Content, w)

			//16.上传好友、加好友
		case "friend_upload":
			UploadFriend(requestBody.Content, w)

		// 共享位置给好友／群组
		//case "share_pos":
		//	SharePos2Friend(requestBody.Content, w)

		//17.响应好友请求
		case "accept_friend_upload":
			AcceptFriendUpload(requestBody.Content, w)

			//18 约人请求
		case "invite_friend":
			InviteFriend(requestBody.Content, w)

			//19 响应约人请求
		case "accept_invite_friend":
			AcceptInvite(requestBody.Content, w)

			//20.处理头像上传
		case "upload_avatar":
			UploadAvatar(requestBody.Content, w)

			//21.查看好友信息
		case "friend_info":
			ViewFriendInfo(requestBody.Content, w)
			
			// 请求共享位置
		case "loc_share":
			SharePosWithFriend(requestBody.Content, w)
			
			// 响应位置共享请求
		case "res_loc_share":
			ResponsePosShare(requestBody.Content, w)
			
			// 删除好友
		case "remove_friend":
			RemoveFriend(requestBody.Content, w)
		}
	}
}

// 写返回给客户端的数据
func write2Client(w http.ResponseWriter, val ResponseInfo) {
	bytes, err := json.Marshal(val)
	if err != nil {
		// todo: add log
		panic("error in encoding json -- func:write2Client")
	}
	if _, err = w.Write(bytes); err != nil {
		// todo: add log
		panic("error in writing response -- func:write2Client")
	}
}
