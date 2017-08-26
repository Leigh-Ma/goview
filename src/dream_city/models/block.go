package models

import (
	"github.com/astaxie/beego/orm"
	"strings"
)

type Block struct {
	TableCommon
	UserId     int64  `orm:"size(64)" show:"用户"`
	Name       string `orm:"size(64)" show:"地块"`
	CityName   string `orm:"size(64)" show:"城市"`
	Meta       orm.TextField `show:"-"`
	Status     string `show:"状态"`
	UpUserId   string `show:"上游用户"`
	DownUserId string `show:"下游用户"`

	InUse      bool    `orm:"column(using)"`
}

func init() {
	orm.RegisterModel(new(Block))
}

func (b *Block) TableName() string{
	return "blocks"
}

func (b *Block) Show(name... string) ([]interface{}){
	return []interface{}{
		b.Id,
		b.UpdatedAt,
		b.OwnerName(),
		b.Name,
		b.CityName,
		b.Status,
		b.UpNames(),
		b.DownNames(),
	}
}

func (b *Block) Titles(name... string)[]string{
	return []string{
		"ID",
		"更新",
		"用户",
		"地块",
		"城市",
		"状态",
		"上游用户",
		"下游用户",
	}
}

func (b *Block) OwnerName() string{
	return mUser.GetUserName(b.UserId)
}

func (b *Block) UpNames() []string{
	names := make([]string, 0)
	for _, id := range strings.Split(b.UpUserId, ",") {
		names = append(names, mUser.GetUserName(id))
	}

	return names
}


func (b *Block) DownNames() []string{
	names := make([]string, 0)
	for _, id := range strings.Split(b.DownUserId, ",") {
		names = append(names, mUser.GetUserName(id))
	}

	return names
}