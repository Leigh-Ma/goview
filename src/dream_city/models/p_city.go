package models

type PositionCities struct {
	TableCommon
	ProvinceId int64     `orm:"column(provice_id);null"`
	CityId     int64     `orm:"column(city_id);null"`
	CityName   string    `orm:"column(city_name);size(255);null"`

}

func (t *PositionCities) TableName() string {
	return "position_cities"
}

