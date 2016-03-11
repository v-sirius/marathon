package db

import (
	"database/sql"
	"fmt"
	//"fmt"
	"log"
)

// 定义数据库及各类操作
type DbOperation struct {
	Db      *sql.DB // 数据库
	DataSrc string  // 数据源
}

var G_db *DbOperation

// 初始化数据库
func InitDbOperation(dsr string) *DbOperation {
	pRet := new(DbOperation)
	pRet.DataSrc = dsr
	return pRet
}

// 打开数据库
func (db *DbOperation) Open() {
	database, err := sql.Open("mysql", db.DataSrc)
	db.Db = database
	if err != nil {
		panic("cannot open database")
	}
}

// 关闭数据库
func (db *DbOperation) Close() {
	if db.Db != nil {
		if db.Db.Close() != nil {
			panic("cannot close database")
		}
	}
}

// 创建数据库表格
func (db *DbOperation) CreateTable() {
	// 标注表
	marksql := "create table mark (id int(20) primary key auto_increment,markname varchar(50)," +
		"mlongitude float,mlatitude float,mdescribe varchar(255),mtype int,creattime varchar(50),mvalue float);"
	smt, err := db.Db.Prepare(marksql)
	checkErr(err)
	smt.Exec()

	//运动员表
	sportsql := "create table sports(id int(20) primary key auto_increment,sname char(50),sdescribe char(255),sfavicon char(50),screatetime char(50) );"
	smt, err = db.Db.Prepare(sportsql)
	checkErr(err)
	smt.Exec()

	//啦啦队
	cheerteam := "create table cheerteam(id int(20) primary key auto_increment,cname char(50),cdescribe char(255),cfavicon char(50),screatetime char(50) );"
	smt, err = db.Db.Prepare(cheerteam)
	checkErr(err)
	smt.Exec()

	//合作伙伴
	copartner := "create table copartner(id int(20) primary key auto_increment,coname char(50)," +
		"codescribe char(255),cotype int,cofavicon char(50),cocreatetime char(50) );"
	smt, err = db.Db.Prepare(copartner)
	checkErr(err)
	smt.Exec()

	//关注表
	focustable := "create table focustable(id int(20) primary key auto_increment,uid char(20)," +
		"followerid char(20),focustime char(50) );"
	smt, err = db.Db.Prepare(focustable)
	checkErr(err)
	smt.Exec()

	//好友表
	friend := "create table friend(id int(20) primary key auto_increment,uid int," +
		"friendid int, addftime varchar(50), shared bool);"
	smt, err = db.Db.Prepare(friend)
	checkErr(err)
	smt.Exec()

	//用户位置表
	userlocation := "create table userlocation(id int(20) primary key auto_increment,uid char(20)," +
		"ulongitude float,ulatitude float,createtime char(50) );"
	smt, err = db.Db.Prepare(userlocation)
	checkErr(err)
	smt.Exec()

	//上线记录
	onlinerecord := "create table onlinerecord(id int(20) primary key auto_increment,uid char(20)," +
		"loginlongitude float,loginlatitude float,logintime char(50) );"
	smt, err = db.Db.Prepare(onlinerecord)
	checkErr(err)
	smt.Exec()

	//日程信息表
	scheduleinfo := "create table scheduleinfo(id int(20) primary key auto_increment,uid char(20)," +
		"stitle char(50),stime datetime,scontent char(255),screattime char(50),stype int );"
	smt, err = db.Db.Prepare(scheduleinfo)
	checkErr(err)
	smt.Exec()

	//第三方登录
	thirdlogin := "create table thirdlogin(id int(20) primary key auto_increment,ttype int," +
		"tuserid char,logintime char(50) );"
	smt, err = db.Db.Prepare(thirdlogin)
	checkErr(err)
	smt.Exec()

	//照片
	photo := "create table photo(pid int(20) primary key auto_increment,uid char(20),plongitude float," +
		"platitude float,uploadphoto char(50),pnote char(20),pcreattime char(50) );"
	smt, err = db.Db.Prepare(photo)
	checkErr(err)
	smt.Exec()

	//用户
	user := "create table user(uid int(20),ucode varchar(50)," +
		"unickname varchar(50),password varchar(50),ufavicon varchar(50),utel char(11)," +
		"utype int,udescribe varchar(255),ucreattime varchar(50), umd5 varchar(50));"
	smt, err = db.Db.Prepare(user)
	checkErr(err)
	smt.Exec()

	//个数统计
	numstatis := "create table numstatis(id int(20) primary key auto_increment,uid char(20),upphotonum int," +
		"focusnum int,fansnum int,updatetime date);"
	smt, err = db.Db.Prepare(numstatis)
	checkErr(err)
	smt.Exec()

	//用户轨迹
	usertraject := "create table usertraject(id int(20) primary key auto_increment,uid char(20),utraject varchar(100)," +
		"updatet datetime);"
	smt, err = db.Db.Prepare(usertraject)
	checkErr(err)
	smt.Exec()

	//定义好友fk1外键
	/*idx_userid := "alter table user add index idx_userid(uid);"
	smt, err = db.Db.Prepare(idx_userid)
	checkErr(err)
	smt.Exec()

	idx_userid = "alter table friend add index idx_friuid(uid);"
	smt, err = db.Db.Prepare(idx_userid)
	checkErr(err)
	smt.Exec()

	idx_userid = "alter table friend add index idx_friendid(friendid);"
	smt, err = db.Db.Prepare(idx_userid)
	checkErr(err)
	smt.Exec()

	frifk1sql := "alter table friend add constraint fk1_fri_user foreign key(friendid) references user(uid);"
	smt, err = db.Db.Prepare(frifk1sql)
	checkErr(err)
	smt.Exec()

	//定义好友fk2外键
	frifk2sql := "alter table friend add constraint fk2_fri_user foreign key(uid) references user(uid);"
	smt, err = db.Db.Prepare(frifk2sql)
	checkErr(err)
	smt.Exec()

	//定义关注表fk1外键
	idx_userid = "alter table focustable add index idx_focusuid(uid);"
	smt, err = db.Db.Prepare(idx_userid)
	checkErr(err)
	smt.Exec()

	idx_userid = "alter table focustable add index idx_focusfollid(followerid);"
	smt, err = db.Db.Prepare(idx_userid)
	checkErr(err)
	smt.Exec()

	focfk1sql := "alter table focustable add constraint fk1_focus_user " +
		"foreign key(uid) references user(uid);"
	smt, err = db.Db.Prepare(focfk1sql)
	checkErr(err)
	smt.Exec()

	//定义关注表fk2外键
	focfk2sql := "alter table focustable add constraint fk2_focus_user " +
		"foreign key(followerid) references user(uid);"
	smt, err = db.Db.Prepare(focfk2sql)
	checkErr(err)
	smt.Exec()

	//定义用户位置表外键
	idx_userid = "alter table userlocation add index idx_userlouid(uid);"
	smt, err = db.Db.Prepare(idx_userid)
	checkErr(err)
	smt.Exec()

	userlofk1sql := "alter table userlocation add constraint fk1_userlo_user foreign key(uid) references user(uid);"
	smt, err = db.Db.Prepare(userlofk1sql)
	checkErr(err)
	smt.Exec()

	//定义上线记录外键

	idx_userid = "alter table onlinerecord add index idx_onlineuid(uid);"
	smt, err = db.Db.Prepare(idx_userid)
	checkErr(err)
	smt.Exec()

	onlinefk1sql := "alter table onlinerecord add constraint fk1_onlinere_user foreign key(uid) references user(uid);"
	smt, err = db.Db.Prepare(onlinefk1sql)
	checkErr(err)
	smt.Exec()

	//定义日程信息外键
	idx_userid = "alter table scheduleinfo add index idx_scheuid(uid);"
	smt, err = db.Db.Prepare(idx_userid)
	checkErr(err)
	smt.Exec()

	schefk1sql := "alter table scheduleinfo add constraint fk1_sche_user foreign key(uid) references user(uid);"
	smt, err = db.Db.Prepare(schefk1sql)
	checkErr(err)
	smt.Exec()

	//定义照片外键
	idx_userid = "alter table photo add index idx_phouid(uid);"
	smt, err = db.Db.Prepare(idx_userid)
	checkErr(err)
	smt.Exec()

	phofk1sql := "alter table photo add constraint fk1_pho_user foreign key(uid) references user(uid);"
	smt, err = db.Db.Prepare(phofk1sql)
	checkErr(err)
	smt.Exec()

	//定义个数统计
	idx_userid = "alter table numstatis add index idx_numsuid(uid);"
	smt, err = db.Db.Prepare(idx_userid)
	checkErr(err)
	smt.Exec()

	numsfk1sql := "alter table numstatis add constraint fk1_nums_user foreign key(uid) references user(uid);"
	smt, err = db.Db.Prepare(numsfk1sql)
	checkErr(err)
	smt.Exec()

	//定义用户轨迹
	idx_userid = "alter table usertraject add index idx_usertra(uid);"
	smt, err = db.Db.Prepare(idx_userid)
	checkErr(err)
	smt.Exec()

	usertfk1sql := "alter table usertraject add constraint fk1_usertra_user foreign key(uid) references user(uid);"
	smt, err = db.Db.Prepare(usertfk1sql)
	checkErr(err)
	smt.Exec()*/

}

//增加数据到数据库表中
func (db *DbOperation) Insert2Table(str string) bool {

	stmt, err := db.Db.Prepare(str)
	checkErr(err)
	res, err := stmt.Exec()
	defer stmt.Close()

	checkErr(err)
	//可以获得插入的id
	id,err2id:=res.LastInsertId()
	if err2id==nil{
	fmt.Println("插入后的id:",id)
	}
	
	//可以获得影响的行数
	i, err := res.RowsAffected()
	if i > 0 && err == nil {
		return true
	} else {
		return false
	}

}

//从数据库表中删除数据
func (db *DbOperation) DelFromTable(str string) bool {

	stmt, err := db.Db.Prepare(str)
	checkErr(err)
	res, err := stmt.Exec()
	defer stmt.Close()

	checkErr(err)
	i, err := res.RowsAffected()
	if i > 0 && err == nil {
		return true
	} else {
		return false
	}

}

func (db *DbOperation) Find(str string) *sql.Rows {
	fmt.Println(str)
	rows, err := db.Db.Query(str)
	fmt.Println("--------------------err", err)
	if err != nil {
		panic("error in: selecting in table")
	}
	//rows.Close()

	return rows
}

//错误检查
func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
