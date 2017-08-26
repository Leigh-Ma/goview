package models

type Entity struct {
	TableCommon

	Name      string    `orm:"column(name);size(255);null"`
	Crc       string    `orm:"column(crc);size(255);null"`
	Url       string    `orm:"column(url);size(255);null"`

	FileType  string    `orm:"column(file_type);null"`
}

func (t *Entity) TableName() string {
	return "entities"
}

