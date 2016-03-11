package webserver

import (
	"log"
	"encoding/hex"
	"bytes"
	"cache"
	"crypto/md5"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 用户名检测
func CheckUser(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	res.Code = "ret_check_user"

	// 用户数据的唯一性，详细请见user定义
	if name, exist := content["tel"]; exist {
		fmt.Println(name)
		var err error

		_, err = cache.G_CacheData.UserSet.GetUserByMobile(name.(string))
		fmt.Println("-------------errorof get mobile-----------")
		fmt.Println(err)

		if _, err = cache.G_CacheData.UserSet.GetUserByMobile(name.(string)); err != nil { //不存在用户时，返回有错误，那么可以注册为新用户
			res.Is_success = true
			res.Content = map[string]interface{}{"is_legal": true} // 合法。该手机号未注册，可以注册

			fmt.Println("--------------------res----------------")
			fmt.Println(res)
			//返回客户端数据
			write2Client(w, res)
			return
		} else if err == nil { //存在用户时，返回无错误，不可以注册为新用户
			res.Is_success = false
			res.Content = map[string]interface{}{"is_legal": false} // 不合法。该手机号已经注册
			fmt.Println("--------------------res----------------")
			fmt.Println(res)
			//返回客户端数据
			write2Client(w, res)
			return
		}
	} else { //no exist name in content
		res.Is_success = false
		res.Content = map[string]interface{}{"is_legal": false}
		fmt.Println("--------------------res----------------")
		fmt.Println(res)
		//返回客户端数据
		write2Client(w, res)
		return
	}
}

// 处理注册
func Register(content map[string]interface{}, w http.ResponseWriter) {
	//var res ResponseInfo
	res := ResponseInfo{Code: "ret_register", Is_success: false}
	
	fmt.Println("-----------------------content--------------")
	fmt.Println(content)

	//res.Code = "ret_register"

	//	defer write2Client(w, res)

	// 必须存在账户和密码
	if _, exist := content["tel"]; !exist {
		//res.Is_success = false
		res.Content = map[string]interface{}{"uid": "", "msg": "no telephone info"}

		fmt.Println("---------------res exist tel??---------------------")
		fmt.Println(res)
		write2Client(w, res)
		return
	}
	if _, exist := content["password"]; !exist {
		//res.Is_success = false
		res.Content = map[string]interface{}{"uid": "", "msg": "no password info"}

		fmt.Println("---------------res exist password??---------------------")
		fmt.Println(res)
		write2Client(w, res)
		return
	}

	// 如果存在，则告之
	if tmpUser, bExist := cache.G_CacheData.UserSet.IsExist(content["tel"].(string)); bExist {
		res.Content = map[string]interface{}{"uid": tmpUser.Code, "msg": "此号码已注册"}
		write2Client(w, res)
		return
	}

	var user cache.User

	user.NickName = content["nickname"].(string)
	user.Password = content["password"].(string)
	log.Println("*********************user.NickName*****************",user.NickName)

      var chars_illegal string
	chars_illegal="'&/*^%$#@\\=-~!"
	bool_islegal:=strings.ContainsAny(user.NickName,chars_illegal)
	if bool_islegal{
		res.Content = map[string]interface{}{"uid": "", "msg": "nickname is illegal,can't have ('&/*^%$#@\\=-~!)"}
		write2Client(w, res)
		return
	}

	// ???这一个需要查看相应的包如何处理
	//user.Img, _ = image.Decode(content["avatar"].([]byte))

	user.Mobile = content["tel"].(string)
	user.Sexual=content["sex"].(string)
	md5Value := md5.Sum([]byte(user.Mobile))
	var tmpBytes []byte;
	for _, v := range md5Value{
		tmpBytes = append(tmpBytes, v)
	}
	user.MobileMd5 = hex.EncodeToString(tmpBytes)
	//user.MobileMd5 = string(md5.Sum(user.Mobile))
	//user.Code = user.Mobile
	user.Category = int(content["u_type"].(float64))

	user.Description = content["description"].(string)
	user.Id = cache.G_CacheData.UserSet.GetLength()

	fmt.Println(time.Now().String())
	createtime := strconv.Itoa(time.Now().Year()) + strconv.Itoa(time.Now().YearDay())
	user.CreateTime = createtime //system time

	fmt.Println("-------------------user createtime-----------------------")
	fmt.Println(user.CreateTime)

	//*********************需要处理下用户码**************
	user.Code = user.Mobile + time.Now().String()
	fmt.Println("-------------------user code-----------------------")
	fmt.Println(user.Code)

	fmt.Println("---------------user---------------------")
	fmt.Println(user)

	//从server写入到cache
	cache.G_CacheData.UserSet.WriteUserFromServer2Cache(user)

	//	var userchan = cache.G_CacheData.UserSet.UserChan
	cache.G_CacheData.UserSet.UserChan <- user
	fmt.Println("--------------------userchan---------------")
	fmt.Println(cache.G_CacheData.UserSet.UserChan)

	// 返回数据
	res.Is_success = true
	res.Content = map[string]interface{}{"uid": user.Code, "msg": "success"}

	fmt.Println("---------------res for register---------------------")
	fmt.Println(res)
	write2Client(w, res)

	//写入DB
	//fmt.Println("starting to  Write2DbUser......")
	//cache.G_CacheData.UserSet.Write2DbUser()
	//fmt.Println("end to  Write2DbUser......")
	return
}

// 处理更新
func UpdateUser(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	//	var user cache.User

	fmt.Println("---------------register content----------")
	fmt.Println(content)

	res.Code = "ret_user_update"
	//uid := content["uid"]
	//user.Id, _ = strconv.Atoi(content["uid"].(string)) // strconv.Atoi(uid)

	// 如果不存在，则直接返回不存在此用户
	/*tmpUser, bExist := cache.G_CacheData.UserSet.IsExist(content["tel"].(string))
	if !bExist{
		//res.Content = map[string]string{"uid": ucode, "msg":"此号码已注册"}
		write2Client(w, res)
		return
	}*/

	// 获取当前用户
	tmpUser, err := cache.G_CacheData.UserSet.GetUserByCode(content["uid"].(string))
	if err != nil {
		//res.Content = map[string]string{"uid": ucode, "msg":"此号码已注册"}
		write2Client(w, res)
		return
	}

	if content["nickname"] == nil {
		tmpUser.NickName = ""
	} else {
		tmpUser.NickName = content["nickname"].(string)
	}
	
	if content["password"] == nil {
		tmpUser.Password = ""
	} else {
		tmpUser.Password = content["password"].(string)
	}
	
	if content["sex"] == nil {
		tmpUser.Sexual = ""
	} else {
		tmpUser.Sexual = content["sex"].(string)
	}
	
	//user.Mobile = content["tel"].(string)   //tel不能更新，唯一标志
	tmpUser.Category, _ = content["type"].(int) //strconv.ParseFloat(content["type"].(string), 64)

	if content["description"] == nil {
		tmpUser.Description = ""
	} else {
		tmpUser.Description = content["description"].(string)
	}
	
	log.Println("####################",tmpUser)
	cache.G_CacheData.UserSet.UpdateCache(tmpUser)
	
	//err = cache.G_CacheData.UserSet.Update2DbUser(tmpUser)
	sql := "update user set unickname='" + tmpUser.NickName +
		"',udescribe='" + tmpUser.Description +
		"',password='" + tmpUser.Password +
		"',usex='" + tmpUser.Sexual +
		"',utype=" + strconv.Itoa(tmpUser.Category) + " where ucode='" + tmpUser.Code + "';"
	
	
	cache.G_CacheData.ChString <- sql
	
	if err == nil {
		res.Is_success = true
		res.Content = map[string]interface{}{"uid": tmpUser.Code, "nickname": tmpUser.NickName, "u_type": tmpUser.Category, "description": tmpUser.Description,"avatar":tmpUser.Img,"sex":tmpUser.Sexual}
	} else {
		res.Is_success = false
		res.Content = map[string]interface{}{}
	}

	fmt.Println("--------------response of user update-----------")
	fmt.Println(res)
	//返回数据到客户端
	write2Client(w, res)
	return
}

//处理登录
func UserLogin(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	fmt.Println("------------request body of user login----------")
	fmt.Println(content)

	res.Code = "ret_login"

	if tel, exist := content["tel"]; !exist {
		fmt.Println(tel)
		res.Is_success = false
		res.Content = map[string]interface{}{"is_successed": false, "msg": "no nickname info", "data": nil}

		fmt.Println("------------response of user login-----------")
		fmt.Println(res)

		write2Client(w, res)
		return
	}
	if pw, exist := content["password"]; !exist {
		fmt.Println(pw)
		res.Is_success = false
		res.Content = map[string]interface{}{"is_successed": false, "msg": "no password info", "data": nil}

		fmt.Println("------------response of user login-----------")
		fmt.Println(res)
		//返回数据到客户端
		write2Client(w, res)
		return
	}

	this_user, err := cache.G_CacheData.UserSet.GetUserByMobile(content["tel"].(string))
	//upload_avatar
	fmt.Println(this_user)
	fmt.Println(this_user.Id)
	fmt.Println("--------------------check??---------------------")
	if err == nil {
		if this_user.Password == content["password"] {
			//			data := map[string]interface{}{"uid": this_user.Id, "code": this_user.Code, "nickname:": this_user.NickName, "img": this_user.Img, "mobile": this_user.Mobile, "u_type": this_user.Category, "description": this_user.Description, "creattm": this_user.CreateTime}
			res.Is_success = true
			res.Content = map[string]interface{}{"msg": "登陆成功", "data": map[string]interface{}{"uid": this_user.Code, "nickname": this_user.NickName, "u_type": this_user.Category, "description": this_user.Description,"avatar":this_user.Img,"sex":this_user.Sexual}}
			

			fmt.Println("------------response of user login-----------")
			fmt.Println(res)
			//返回数据到客户端
			write2Client(w, res)
			return

		} else {
			res.Is_success = false
			res.Content = map[string]interface{}{"msg": "login failed,wrong password!", "data": nil}

			fmt.Println("------------response of user login-----------")
			fmt.Println(res)
			//返回数据到客户端
			write2Client(w, res)
			return
		}
	} else {
		res.Is_success = false
		res.Content = map[string]interface{}{"msg": "login failed,wrong tel!", "data": nil}

		fmt.Println("------------response of user login-----------")
		fmt.Println(res)
		//返回数据到客户端
		write2Client(w, res)
		return
	}

}

//处理avatar upload
func UploadAvatar(content map[string]interface{}, w http.ResponseWriter) {
	var res ResponseInfo
	fmt.Println("------------request body of user upload avatar----------")
	fmt.Println(content)

	res.Code = "ret_upload_avatar"

	if avatar, exist := content["avatar"]; !exist {
		fmt.Println(avatar)
		res.Is_success = false
		res.Content = map[string]interface{}{"url": nil}

		fmt.Println("------------response of user upload avatar-----------")
		fmt.Println(res)

		write2Client(w, res)
		return
	}

	//avatar 类型为byte[],转换后存储到服务器的一个文件夹中，返回url地址
	ff := content["avatar"].([]byte) //convert  to byte[]
	bbb := bytes.NewBuffer(ff)
	m, _, _ := image.Decode(bbb) //decode

	f, _ := os.Create("output.jpg") //output as the related user's avatar
	fmt.Println("f name:", f.Name())
	defer f.Close()
	jpeg.Encode(f, m, nil)

	str := getCurrentPath()
	fmt.Println("end:", str)
	urlStr := str + f.Name()
	fmt.Println("end:", urlStr)
	res.Is_success = true
	res.Content = map[string]interface{}{"url": urlStr}

	fmt.Println("------------response of user upload avatar-----------")
	fmt.Println(res)

	write2Client(w, res)
	return
}

func getCurrentPath() string {
	file, _ := exec.LookPath(os.Args[0])
	fmt.Println("file:", file)
	path, _ := filepath.Abs(file)
	fmt.Println("path:", path)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	return ret
}
