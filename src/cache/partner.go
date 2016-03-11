package cache

import (
	"errors"
	//	"image"
	"db"
	"fmt"
	_ "image/jpeg"
	"strconv"
	"sync"
	"time"
)

// 合作伙伴
type Partner struct {
	Id          string
	Name        string
	Description string
	Category    int
	Img         string //image.Image
	TmCreate    string //time.Time
}

// 合作伙伴内存集合
type PartnerSet struct {
	allPartners []Partner
	id2Partner  map[string]*Partner
	//Name2Partner map[string]*Partner

	PartnerChan chan Partner

	rwLock sync.RWMutex // 读写锁
}

// 初始化函数
func InitPartnerSet() *PartnerSet {
	pRet := new(PartnerSet)
	pRet = &PartnerSet{allPartners: make([]Partner, 0), id2Partner: make(map[string]*Partner, 0), PartnerChan: make(chan Partner, 0)}
	return pRet
}

// 写入数据
func (partnerSet *PartnerSet) WritePartner2Cache(partner Partner) {
	partnerSet.rwLock.Lock()
	defer func() {
		partnerSet.rwLock.Unlock()
	}()

	partnerSet.allPartners = append(partnerSet.allPartners, partner)

	partnerSet.id2Partner[partner.Id] = &partner
}

//读取函数
func (partnerSet *PartnerSet) GetPartnerById(id string) (Partner, error) {
	partnerSet.rwLock.RLock()
	defer func() {
		partnerSet.rwLock.RUnlock()
	}()

	if idx, flg := partnerSet.id2Partner[id]; flg {
		return *idx, nil
	}
	return Partner{Id: ""}, errors.New("not exist related partner data")
}

//读取所有数据
func (partnerSet *PartnerSet) GetAllPartners(time2check string) ([]Partner, error) {

	var partPartner []Partner
	const longForm = "2006-01-02 15:04:05"
	t, _ := time.Parse(longForm, time2check)

	for _, partner := range partnerSet.allPartners {

		ct, _ := time.Parse(longForm, partner.TmCreate)
		if ct.After(t) {
			partPartner = append(partPartner, partner)
		}
		fmt.Println("------------ct---------")
		fmt.Println(ct.String())
	}
	if len(partPartner) == 0 {
		return partPartner, errors.New("not exist related partner data")
	}
	return partPartner, nil

}

// 加载数据
func (partnerSet *PartnerSet) LoadData(pDb *db.DbOperation) {
	var pPartner Partner

	selectstr := "select id,coname,codescribe,cotype,cofavicon,cocreatetime from copartner"
	rows := pDb.Find(selectstr)

	for rows.Next() {
		err := rows.Scan(&pPartner.Id, &pPartner.Name, &pPartner.Description, &pPartner.Category, &pPartner.Img, &pPartner.TmCreate)
		if err != nil {
			panic("error in:scanning the table partner")
		}

		partnerSet.WritePartner2Cache(pPartner)
	}

}

//channel写入数据库
func (partnerSet *PartnerSet) Write2DbPartner() {
	var partnerchan = G_CacheData.PartnerSet.PartnerChan
	for partner := range partnerchan {
		var sql = "insert into copartner(id,coname,codescribe,cotype,cofavicon)" +
			" values ('" + partner.Id +
			"','" + partner.Name + "','" + partner.Description + "','" + strconv.Itoa(partner.Category) + "','" +
			partner.Img + ")"

		if !db.G_db.Insert2Table(sql) {
			fmt.Println("partner表插入不成功！")
		}

	}

}
