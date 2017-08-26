package models

type BlockComponent struct {
	TableCommon

	BlockName     string    `orm:"column(block_name);size(64)"`
	ComponentType string    `orm:"column(component_type);size(32)"`
	ComponentCode string    `orm:"column(component_code);size(64)"`
	Position      string    `orm:"column(position);size(255);null"`
	Rotation      string    `orm:"column(rotation);null"`
	Scale         string    `orm:"column(scale);null"`
	Tail          string    `orm:"column(tail);size(255);null"`

	Using         int8      `orm:"column(using);null"`
}

func (t *BlockComponent) TableName() string {
	return "block_components"
}

