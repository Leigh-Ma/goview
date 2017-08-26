package models

type Meterial struct {
	TableCommon

	Name      string    `orm:"column(name);size(255);null"`
	Crc       string    `orm:"column(crc);size(255);null"`
	Url       string    `orm:"column(url);size(255);null"`
	Tags      string    `orm:"column(tags);size(255);null"`
	Attrs     string    `orm:"column(attrs);size(255);null"`
	Picture   string    `orm:"column(picture);size(255);null"`
	Thumbnail string    `orm:"column(thumbnail);size(255);null"`
}

func (t *Meterial) TableName() string {
	return "meterials"
}
