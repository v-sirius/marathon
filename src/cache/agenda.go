package cache

import (
	"db"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// 日程表
type Agenda struct {
	Id         int    // 日程唯一标识
	Title      string // 标题
	AgendaTime string // 日程时间 "Jan 2, 2006 at 3:04pm (MST)"
	CreateTime string // 创建日程的时间
	Content    string // 日程内容
	//Category   int    // 日程类别
}

// 所有日程表
type AgendaSet struct {
	allAgenda []Agenda
	//id2Agenda map[string]*Agenda
	//title2Agenda map[string]*Agenda // titile是否唯一不确定
	AgendaChan chan Agenda

	rwLock sync.RWMutex // 读写锁
}

// 初始化日程表缓存
func InitAgendaSet() *AgendaSet {
	pRet := new(AgendaSet)
	pRet = &AgendaSet{allAgenda: make([]Agenda, 0), AgendaChan: make(chan Agenda, 0)}

	//pRet.LoadData(pDb)
	return pRet
}

// 写入数据
func (agendaSet *AgendaSet) WriteAgenda2Cache(agenda Agenda) {
	agendaSet.rwLock.Lock()
	agendaSet.allAgenda = append(agendaSet.allAgenda, agenda)

	func() {
		agendaSet.rwLock.Unlock()
	}()
	//agendaSet.id2Agenda[agenda.Id] = &agenda
}

//读取函数
//func (agendaSet *AgendaSet) GetPartnerById(id string) (Agenda, error) {
//	agendaSet.rwLock.RLock()
//	defer func() {
//		agendaSet.rwLock.RUnlock()
//	}()

//	if idx, flg := agendaSet.id2Agenda[id]; flg {
//		return *idx, nil
//	}
//	return Agenda{Id: ""}, errors.New("not exist related agenda data")
//}

//读取所有数据
func (agendaSet *AgendaSet) GetAllAgendas(time2check string) ([]Agenda, error) {
	var partAgenda []Agenda
	const longForm = "2006-01-02 15:04:05"
	t, _ := time.Parse(longForm, time2check)
	fmt.Println("------------t---------")
	fmt.Println(t.String())

	for _, agenda := range agendaSet.allAgenda {
		ct, _ := time.Parse(longForm, agenda.CreateTime)
		if ct.After(t) {
			partAgenda = append(partAgenda, agenda)
		}
	}
	if len(partAgenda) == 0 {
		return partAgenda, errors.New("not exist related agenda data")
	}
	return partAgenda, nil
}

// 加载数据
func (agendaSet *AgendaSet) LoadData(pDb *db.DbOperation) {
	var agenda Agenda
	selectstr := "select id,stitle,stime,scontent,screattime from scheduleinfo"
	rows := pDb.Find(selectstr)

	for rows.Next() {
		err := rows.Scan(&agenda.Id, &agenda.Title, &agenda.AgendaTime, &agenda.Content, &agenda.CreateTime)
		if err != nil {
			panic("error in:scanning the table agenda")
		}

		agendaSet.WriteAgenda2Cache(agenda)
	}
}

//channel写入数据库
func (agendaSet *AgendaSet) Write2DbAgenda() {
	var agendachan = G_CacheData.AgendaSet.AgendaChan
	for agenda := range agendachan {
		var sql = "insert into scheduleinfo(id,uid,stitle,stime,scontent,screattime)" +
			" values ('" + strconv.Itoa(agenda.Id) + "','" + agenda.Title + "','" + agenda.AgendaTime + "','" +
			agenda.Content + "','" + agenda.CreateTime + "')"

		if !db.G_db.Insert2Table(sql) {
			fmt.Println("agenda表插入不成功！")
		}

	}

}
