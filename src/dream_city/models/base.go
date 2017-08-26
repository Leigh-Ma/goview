package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type ITable interface {
	TableName() string
	SetId(id int64)
}


type TableCommon struct {
	Id             int64     `orm:"auto" show:"ID"`
	CreatedAt      time.Time `orm:"auto_now_add;type(datetime)" show:"-"`
	UpdatedAt      time.Time `orm:"auto_now;type(datetime)" show:"更新"`
}

func (t *TableCommon) SetId(id int64){
	t.Id = id
}

func FindById(obj ITable, id... int64) error {
	if len(id) > 0 {
		obj.SetId(id[0])
	}
	return orm.NewOrm().Read(obj)
}

func FindBy(field string, value interface{}, obj ITable) error {
	return orm.NewOrm().QueryTable(obj.TableName()).Filter(field, value).One(obj)
}

func NewQuery(obj interface{}) (orm.QuerySeter) {
	return orm.NewOrm().QueryTable(obj)
}

func Insert(v ITable) (int64, error) {
	return orm.NewOrm().Insert(v)
}