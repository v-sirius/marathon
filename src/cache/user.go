package cache

import (
//	"strings"
	"db"
	//	"image"
	//	"database/sql"
	"errors"
	"fmt"
	"log"
	_ "github.com/Go-SQL-Driver/MySQL"
	_ "image/jpeg"
	"strconv"
	"sync"
	//	"time"
	//"crypto/md5"
)

const USER_CHAN_NUM = 1

// 定义用户
type User struct {
	Id   int    // 唯一标识
	Code string // 用户代码

	NickName string // 昵称
	Password string //密码
	Img      string // 用户图像 －－ 存在特定的目录下

	Mobile    string // 手机号
	MobileMd5 string // 手机号生成的md5值

	Category    int    // 用户类别
	Description string // 描述
	Sexual string
	CreateTime  string //time.Time // 创建时间
}

// 集合，只有id，用户代码和手机号是唯一的
type UserSet struct {
	allUser       []User         // 所有的用户
	id2User       map[int]int    // id到用户的映射
	code2User     map[string]int // 用户唯一代码到用户的映射
	mobile2User   map[string]int // 电话到用户的映射
	mobileMd2User map[string]int // 用户电话md5值到用户的映射

	//code2Bytes	  map[string][]byte  // 用户唯一代码到byte数组的映射，其中byte数组保存用户图片

	UserChan chan User

	rwLock sync.RWMutex // 读写锁
}

func (usrSet UserSet) GetLength() int{
	usrSet.rwLock.RLock()
	defer func(){
		usrSet.rwLock.RUnlock()
	}()
	
	return len(usrSet.allUser)
}

// 初始化函数
func InitUserSet() *UserSet {
	pRet := new(UserSet)

	pRet = &UserSet{allUser: make([]User, 0), id2User: make(map[int]int, 0), code2User: make(map[string]int, 0), mobile2User: make(map[string]int, 0), UserChan: make(chan User, USER_CHAN_NUM), mobileMd2User: make(map[string]int, 0)}

	return pRet
}

// 检查用户是否存在－－根据唯一的电话号码判断
func (usr UserSet) IsExist(tel string) (User, bool) {
	if i, flag := usr.mobile2User["tel"]; flag {
		return usr.allUser[i], true
	}
	return User{Id: -1}, false
}

// 从db写入数据到cache
func (pSet *UserSet) WriteUserFromDb2Cache(usr User) {
	fmt.Println("--------write to cache--------")
	fmt.Println(usr)

	pSet.rwLock.Lock()
	defer func() {
		pSet.rwLock.Unlock()
	}()

	size := len(pSet.allUser)
	//usr.Id = size
	fmt.Println("-----------size----------")
	fmt.Println(size)
	pSet.allUser = append(pSet.allUser, usr)

	pSet.id2User[usr.Id] = size 
	pSet.code2User[usr.Code] = size 
	pSet.mobile2User[usr.Mobile] = size 
	pSet.mobileMd2User[usr.MobileMd5] = size 
}

// 更新用户信息到cache
func (pSet *UserSet) UpdateCache(usr User) {
	fmt.Println("--------UpdateCache--------")
	fmt.Println(usr)

	pSet.rwLock.Lock()
	defer func() {
		pSet.rwLock.Unlock()
	}()

	fmt.Println("-----------usr----------")
	fmt.Println(usr)
	
	pSet.allUser[usr.Id].Sexual = usr.Sexual
	pSet.allUser[usr.Id].Password = usr.Password
	pSet.allUser[usr.Id].NickName = usr.NickName
	pSet.allUser[usr.Id].Category = usr.Category
	pSet.allUser[usr.Id].Description = usr.Description
}

// 从server写入数据到cache
func (pSet *UserSet) WriteUserFromServer2Cache(usr User) {
	fmt.Println("--------write to cache--------")

	pSet.rwLock.Lock()
	defer func() {
		pSet.rwLock.Unlock()
	}()

	size := len(pSet.allUser)
	//usr.Id = size
	fmt.Println(usr)

	fmt.Println("-----------size----------")
	fmt.Println(size)
	pSet.allUser = append(pSet.allUser, usr)

	pSet.id2User[size] = size
	pSet.code2User[usr.Code] = size
	pSet.mobile2User[usr.Mobile] = size
	pSet.mobileMd2User[usr.MobileMd5] = size
}

// 读取数据
func (pSet *UserSet) GetUserById(id int) (User, error) {
	pSet.rwLock.RLock()
	defer func() {
		pSet.rwLock.RUnlock()
	}()

	if idx, flg := pSet.id2User[id]; flg {
		return pSet.allUser[idx], nil
	}

	return User{Id: -1}, errors.New("not exist related user data")
}

// 根据code获取用户
func (pSet *UserSet) GetUserByCode(code string) (User, error) {
	pSet.rwLock.RLock()
	defer func() {
		pSet.rwLock.RUnlock()
	}()

	if idx, flg := pSet.code2User[code]; flg {
		return pSet.allUser[idx], nil
	}

	return User{Id: -1}, errors.New("not exist related user data")
}

// 根据用户md5获取用户
func (pSet *UserSet) GetUserByMd5(md string)(User, error){
	pSet.rwLock.RLock()
	defer func() {
		pSet.rwLock.RUnlock()
	}()

	if idx, flg := pSet.mobileMd2User[md]; flg {
		log.Println("idx:", idx)
		return pSet.allUser[idx], nil
	}

	return User{Id: -1}, errors.New("not exist related user data")
}

// 根据mobile获取用户
func (pSet *UserSet) GetUserByMobile(mobile string) (User, error) {
	pSet.rwLock.RLock()
	defer func() {
		pSet.rwLock.RUnlock()
	}()
	fmt.Println("-------------------------mobile-------------")
	fmt.Println(mobile)
	fmt.Println("-------------------------pSet.mobile2User[mobile]-------------")
	fmt.Println(pSet.mobile2User[mobile])
	if idx, flg := pSet.mobile2User[mobile]; flg { //已经存在用户，返回错误为空
		fmt.Println("------------flg----------------")
		fmt.Println(flg)
		fmt.Println("-------------------GetUserByMobile-----idx--------")
		fmt.Println(idx)
		fmt.Println("------------------check2???-------")
		return pSet.allUser[idx], nil
	} else {
		//不存在用户，有返回错误
		fmt.Println("---------------check3--------------------")

		return User{Id: -1}, errors.New("not exist related user data")
	}
}

// 加载用户数据
func (pUsersSet *UserSet) LoadData(pDb *db.DbOperation) {
	var user User
	fmt.Println("begining:load userset........")
	selectstr := "select uid,ucode,unickname,password,ufavicon,utel,utype,udescribe,ucreattime, umd5 from user"
	rows := pDb.Find(selectstr)

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Code, &user.NickName, &user.Password, &user.Img, &user.Mobile, &user.Category, &user.Description, &user.CreateTime, &user.MobileMd5)

		if err != nil {
			panic("error in:scanning the table user")
		}
		//md5Value := md5.Sum([]byte(user.Mobile))
		//user.MobileMd5 = string(md5Value[:])

		fmt.Println("-------------------Load data to cache--------------------")
		fmt.Println(user)
		pUsersSet.WriteUserFromDb2Cache(user)
	}
}

// byte 2 string

//注册 从管道写入数据库
func (pUsersSet *UserSet) Write2DbUser() {
	fmt.Println("------------entering write2db-----------")

	for user := range pUsersSet.UserChan {

		fmt.Println("------------------user in db-------------------")
		fmt.Println(user)
		
		//user.NickName:=strings.Replace(user.NickName,"'","",-1)
		var sql = "insert into user(uid,ucode,unickname,password,ufavicon,utel,utype,udescribe,ucreattime, umd5,usex) values (" + strconv.Itoa(user.Id) + ",'" + user.Code + "','" + user.NickName + "','" + user.Password + "','" + user.Img + "','" + user.Mobile + "'," + strconv.Itoa(user.Category) + ",'" + user.Description + "','" + user.CreateTime + "','" + user.MobileMd5 + "','"+user.Sexual+"');"
		log.Println("sql:", sql)
		if !db.G_db.Insert2Table(sql) {
			fmt.Println("user表插入不成功！")
		} else {
			fmt.Println("entering weite2db OK ! ! !")
		}
	}
}

//更新数据库
func (pUsersSet *UserSet) Update2DbUser(user User) error {
	var sql = "update user set ufavicon='" + user.Img +
		"',udescribe='" + user.Description +
		"',utype=" + strconv.Itoa(user.Category) + " where ucode='" + user.Code + "';"

	fmt.Println("更新字符串:", sql)
	if !db.G_db.Insert2Table(sql) {
		fmt.Println("user表更新不成功！")
		return errors.New("update error!!!")
	}
	return nil
}
