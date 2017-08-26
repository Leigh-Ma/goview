package models


type MeterialTag struct {
	TableCommon
	Name      string    `orm:"column(name);size(255);null"`
	CreaterId int       `orm:"column(creater_id);null"`

}

func (t *MeterialTag) TableName() string {
	return "meterial_tags"
}

