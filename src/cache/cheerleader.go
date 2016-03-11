package cache

import (
	"db"
	"errors"
	"fmt"
	//"image"
	_ "image/jpeg"
	"sync"
	"time"
)

// 啦啦队
type Cheerleader struct {
	Id          string // 唯一ID
	Name        string // 名称
	Description string // 运动员详细描述
	Img         string //image.Image // 运动员图片
	TmCreate    string //time.Time
}

// 啦啦队集合
type CheerleaderSet struct {
	allCheerLeader []Cheerleader
	id2CheerLeader map[string]*Cheerleader
	//name2Athlete map[string]*Cheerleader // 啦啦队名称是否唯一不确定
	CheerleaderChan chan Cheerleader

	rwLock sync.RWMutex // 读写锁
}

// 初始化函数
func InitCheerleaderSet() *CheerleaderSet {
	pRet := new(CheerleaderSet)
	pRet = &CheerleaderSet{allCheerLeader: make([]Cheerleader, 0), id2CheerLeader: make(map[string]*Cheerleader, 0), CheerleaderChan: make(chan Cheerleader, 0)}

	//pRet.LoadData(pDb)
	return pRet
}

// 写入数据to Cache
func (cheerleaderSet *CheerleaderSet) WriteCheerleader2Cache(cheer Cheerleader) {
	cheerleaderSet.rwLock.Lock()
	defer func() {
		cheerleaderSet.rwLock.Unlock()
	}()

	cheerleaderSet.allCheerLeader = append(cheerleaderSet.allCheerLeader, cheer)

	cheerleaderSet.id2CheerLeader[cheer.Id] = &cheer
}

//读取函数
func (cheerleaderSet *CheerleaderSet) GetCheerleaderById(id string) (Cheerleader, error) {
	cheerleaderSet.rwLock.RLock()
	defer func() {
		cheerleaderSet.rwLock.RUnlock()
	}()

	if idx, flg := cheerleaderSet.id2CheerLeader[id]; flg {
		return *idx, nil
	}
	return Cheerleader{Id: ""}, errors.New("not exist related cheerleader data")
}

//读取某个时间点后的数据
func (cheerleaderSet *CheerleaderSet) GetAllCheerleaders(time2check string) ([]Cheerleader, error) {

	var partCheerleader []Cheerleader
	const longForm = "2006-01-02 15:04:05"
	t, _ := time.Parse(longForm, time2check)
	fmt.Println("------------t---------")
	fmt.Println(t.String())

	for _, cheerleader := range cheerleaderSet.allCheerLeader {
		fmt.Println("---------------cheerleader:", cheerleader)
		ct, _ := time.Parse(longForm, cheerleader.TmCreate)
		if ct.After(t) {
			partCheerleader = append(partCheerleader, cheerleader)
		}

		fmt.Println(ct.String())
	}
	if len(partCheerleader) == 0 {
		return partCheerleader, errors.New("not exist related cheerleaders data")
	}
	return partCheerleader, nil

}

// 加载数据
func (cheerleaderSet *CheerleaderSet) LoadData(pDb *db.DbOperation) {
	var cheerleader Cheerleader

	selectstr := "select id,cname,cdescribe,cfavicon,screatetime from cheerteam"
	rows := pDb.Find(selectstr)

	for rows.Next() {
		err := rows.Scan(&cheerleader.Id, &cheerleader.Name, &cheerleader.Description, &cheerleader.Img, &cheerleader.TmCreate)
		if err != nil {
			panic("error in:scanning the table agenda")
		}
		fmt.Println("-----------------------cheerleader:", cheerleader)
		cheerleaderSet.WriteCheerleader2Cache(cheerleader)
	}

}

//channel写入数据库
func (cheerleaderSet *CheerleaderSet) Write2DbCheerleader() {
	var cheerleaderchan = G_CacheData.CheerleaderSet.CheerleaderChan
	for cheerleader := range cheerleaderchan {
		var sql = "insert into cheerteam(id,cname,cdescribe,cfavicon)" +
			" values ('" + cheerleader.Id + "','" + cheerleader.Name + "','" + cheerleader.Description + "','" + cheerleader.Img + "')"

		if !db.G_db.Insert2Table(sql) {
			fmt.Println("Cheerleader表插入不成功！")
		}

	}

}
