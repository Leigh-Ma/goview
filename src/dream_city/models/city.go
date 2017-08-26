package models

import "github.com/astaxie/beego/orm"

type City struct {

	TableCommon
	Name          string  `orm:"size(32)"`
	CityCode      string  `orm:"size(16);column(city_id)"`
	Province      string  `orm:"size(64);column(prov_id)"`

	CreateUserId  int64

	BirthCoord    string  `orm:"column(birth_coordinate)"`
	CenterLng     string  `orm:"size(32)"`
	CenterLat     string  `orm:"size(32)"`
	LngRatio      string  `orm:"size(32)"`
	LatRatio      string  `orm:"size(32)"`

	BlockFile     string
	BlockFileCrc  string  `orm:"size(32)"`

	IsOld         bool
	IsReleased    bool    `orm:"column(does_released)"`
	IsHot         bool    `orm:"column(is_hotcity)"`

	Config        orm.JSONField `orm:"null"`
}

func (t *City) TableName() string{
	return "cities"
}

func (t *City) Show(name... string) ([]interface{}){
	return []interface{}{
		t.Id,
		t.UpdatedAt,
		t.Name,
		t.CityCode,
		t.Province,
		t.CreatorName(),
		t.BirthCoord,
		t.CenterLng + ", " + t.CenterLat,
		t.LngRatio + ", " + t.LatRatio,
		t.Config,
	}
}

func (t *City) Titles(name... string)[]string{
	return []string{
		"ID",
		"更新",
		"城市",
		"编码",
		"省份",
		"用户",
		"降生点",
		"中心点",
		"比例",
		"配置",
	}
}

func (t *City) CreatorName() string{
	return mUser.GetUserName(t.CreateUserId)
}