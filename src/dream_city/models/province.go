package models

type Province struct {
	TableCommon
	ProviceId   int       `orm:"column(provice_id);null"`
	ProviceName string    `orm:"column(provice_name);size(255);null"`

}

func (t *Province) TableName() string {
	return "provices"
}

