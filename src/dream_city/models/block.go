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

func (b *Block) HiddenCols(auth string) ([]string) {
	return []string{"Status", "fake"}
}

func (b *Block) WrapCols(col string) string {
	switch col {
	case "UserId":
		return b.OwnerName()
	case "UpUserId":
		return strings.Join(b.UpNames(), ", ")
	case "DownUserId":
		return strings.Join(b.DownNames(), ", ")
	}
	return ""
}

func (b *Block) TableName() string{
	return "blocks"
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