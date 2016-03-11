package cache

import (
	"db"
	"errors"
	"log"
	"strconv"
	"sync"
	//	"time"
)

// 位置信息
type Position struct {
	Longitude float64
	Latitude  float64
}

// 标注信息
type Label struct {
	Id          string   // 标注唯一ID
	Name        string   // 标注名称
	Pos         Position // 位置信息
	Description string   // 描述
	Category    int      // 标注种类
	CreatTime   string   //  创建时间
	Value       float64  //  标注值
}

// 所有标注信息，以及名称、位置、值以及id到标注的映射
// 是否下面定义的map是合适的？
type LabelSet struct {
	allLabels []Label
	//name2Label map[string]*Label
	//pos2Label  map[Position]*Label
	//val2Label  map[float64]*Label

	category2Label map[int][]Label
	LabelChan      chan Label

	rwLock sync.RWMutex // 读写锁
}

// 初始化函数
func InitLabelSet() *LabelSet {
	pRet := new(LabelSet)
	pRet = &LabelSet{allLabels: make([]Label, 0), category2Label: make(map[int][]Label, 0), LabelChan: make(chan Label, 1)}

	//pRet.LoadLabelSet(pDb)
	return pRet
}

// 写入数据
func (labelSet *LabelSet) WriteLabel2Cache(label Label) {
	labelSet.rwLock.Lock()
	defer func() {
		labelSet.rwLock.Unlock()
	}()

	//	idx := len(labelSet.allLabels)

	labelSet.allLabels = append(labelSet.allLabels, label)
	log.Println("-----------------cache:all labels------------------")
	log.Println(labelSet.allLabels)
	if _, isExist := labelSet.category2Label[label.Category]; isExist {
		labelSet.category2Label[label.Category] = make([]Label, 0)
	}
	labelSet.category2Label[label.Category] = append(labelSet.category2Label[label.Category], label)
	log.Println("-----------------cache:map labels------------------")
	log.Println(labelSet.category2Label[2])
}

//根据Category读取函数
func (labelSet LabelSet) GetPosByCategory(posType int) ([]Label, error) {
	labelSet.rwLock.RLock()
	defer func() {
		labelSet.rwLock.RUnlock()
	}()

	if posType == -1 { //-1，返回所有数据，无错误
		return labelSet.allLabels, nil
	}

	if pos, flg := labelSet.category2Label[posType]; flg { //返回相应类型的标注信息，无返回错误
		log.Println("-------------return the labeled position of related postype-------------")
		log.Println(pos)
		return pos, nil
	}

	return []Label{}, errors.New("not exist related label data")
}

// 加载函数
func (pLabelSet *LabelSet) LoadLabelSet(pDb *db.DbOperation) {
	var label Label
	log.Println("----------------------to load data--------------------")

	selectstr := "select * from mark"
	rows := pDb.Find(selectstr)
	log.Println("------------entering rows-----------------")
	for rows.Next() {
		//fmt.Println("-------rows.next-----------")
		err := rows.Scan(&label.Id, &label.Name, &label.Pos.Longitude, &label.Pos.Latitude, &label.Description, &label.Category, &label.CreatTime, &label.Value)
		log.Println(err)
		if err != nil {
			panic("error in:scanning the table user")
		}

		//fmt.Println("-------------------Load data to cache--------------------")
		//fmt.Println(label)
		pLabelSet.WriteLabel2Cache(label)
	}
}

//channel写入数据库
func (pLabelSet *LabelSet) Write2DbLabel() {
	var labelchan = G_CacheData.LabelSet.LabelChan
	for label := range labelchan {
		var sql = "insert into mark(id,markname,mlongitude,mlatitude,mdescribe,mtype,creattime,mvalue)" +
			" values ('" + label.Id + "','" + label.Name + "'," + strconv.FormatFloat(label.Pos.Longitude, 'f', -1, 64) + "," + strconv.FormatFloat(label.Pos.Latitude, 'f', -1, 64) + ",'" + label.Description + "'," + strconv.Itoa(label.Category) + ",'" + label.CreatTime + "'," + strconv.FormatFloat(label.Value, 'f', -1, 64) + ")"

		if !db.G_db.Insert2Table(sql) {
			log.Println("Label表插入不成功！")
		}

	}

}
