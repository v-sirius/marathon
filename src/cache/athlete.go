package cache

import (
	"db"
	"errors"
	"log"
	//	"image"
	_ "image/jpeg"
	"sync"
	"time"
)

// 运动员
type Athlete struct {
	Id          string // 唯一ID
	Name        string // 名称
	Description string // 运动员详细描述
	Img         string // 运动员图片
	TmCreate    string
}

// 运动员集合
type AthleteSet struct {
	allAthlete []Athlete
	id2Athlete map[string]*Athlete
	//name2Athlete map[string]*Athlete //姓名是否唯一不确定

	AthleteChan chan Athlete

	rwLock sync.RWMutex // 读写锁
}

// 初始化函数
func InitAthleteSet() *AthleteSet {
	pRet := new(AthleteSet)
	pRet = &AthleteSet{allAthlete: make([]Athlete, 0), id2Athlete: make(map[string]*Athlete), AthleteChan: make(chan Athlete, 1)}

	//pRet.LoadData(pDb)
	return pRet
}

// 写入数据
func (athleteSet *AthleteSet) WriteAthlete2Cache(athlete Athlete) {
	athleteSet.rwLock.Lock()
	defer func() {
		athleteSet.rwLock.Unlock()
	}()

	athleteSet.allAthlete = append(athleteSet.allAthlete, athlete)

	athleteSet.id2Athlete[athlete.Id] = &athlete

}

//根据ID读取函数
func (athleteSet *AthleteSet) GetAthleteById(id string) (Athlete, error) {
	athleteSet.rwLock.RLock()
	defer func() {
		athleteSet.rwLock.RUnlock()
	}()

	if idx, flg := athleteSet.id2Athlete[id]; flg {
		return *idx, nil
	}
	return Athlete{Id: ""}, errors.New("not exist related athlete data")
}

//读取所有数据
func (athleteSet *AthleteSet) GetAllAthletes(time2check string) ([]Athlete, error) {
	var partAthlete []Athlete
	const longForm = "2006-01-02 15:04:05"
	t, _ := time.Parse(longForm, time2check)

	for _, athlete := range athleteSet.allAthlete {

		ct, _ := time.Parse(longForm, athlete.TmCreate)
		if ct.After(t) {
			partAthlete = append(partAthlete, athlete)
		}
		log.Println("------------ct---------")
		log.Println(ct.String())
	}
	if len(partAthlete) == 0 {
		return partAthlete, errors.New("not exist related athletes data")
	}
	return partAthlete, nil
}

// 加载数据
func (athleteSet *AthleteSet) LoadData(pDb *db.DbOperation) {
	var athlete Athlete

	selectstr := "select id,sname,sdescribe,sfavicon,screatetime from sports"
	rows := pDb.Find(selectstr)

	for rows.Next() {
		err := rows.Scan(&athlete.Id, &athlete.Name, &athlete.Description, &athlete.Img, &athlete.TmCreate)
		if err != nil {
			panic("error in:scanning the table agenda")
		}
		log.Println("-----------------------athlete:", athlete)
		athleteSet.WriteAthlete2Cache(athlete)
	}
}

//channel写入数据库
func (athleteSet *AthleteSet) Write2DbAthlete() {
	//var athletechan = G_CacheData.AthleteSet.AthleteChan
	for athlete := range athleteSet.AthleteChan {
		var sql = "insert into sports(id,sname,sdescribe,sfavicon)" +
			" values ('" + athlete.Id + "','" + athlete.Name + "','" + athlete.Description +
			"','" + athlete.Img + "')"

		if !db.G_db.Insert2Table(sql) {
			log.Println("athlete表插入不成功！")
		}

	}

}
