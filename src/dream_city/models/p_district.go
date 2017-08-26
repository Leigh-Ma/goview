package models

type PDistrict struct {
	TableCommon
	PCityId      string    `orm:"column(city_id);size(255);null"`
	Code         string    `orm:"column(country_id);size(255);null"`
	Name         string    `orm:"column(country_name);size(255);null"`

}

func (t *PDistrict) TableName() string {
	return "position_countries"
}

