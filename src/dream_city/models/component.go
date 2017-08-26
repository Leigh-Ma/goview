package models

type Component struct {
	TableCommon
	Name      string    `orm:"column(name);size(255);null"`
	Info      string    `orm:"column(info);null"`

}

func (t *Component) TableName() string {
	return "components"
}

