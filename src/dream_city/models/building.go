package models


type Buildings struct {
	TableCommon

	BlockName  string    `orm:"column(block_name);size(64)"`
	Guid       string    `orm:"column(guid);size(64);null"`
	IsBox      int8      `orm:"column(is_box)"`
	Tag        string    `orm:"column(tag);size(32);null"`
	Name       string    `orm:"column(name);size(128);null"`
	FileCrc    string    `orm:"column(file_crc);size(64)"`
	FileName   string    `orm:"column(file_name);size(64)"`
	FileType   string    `orm:"column(file_type);size(32);null"`
	Position   string    `orm:"column(position);size(255);null"`
	Rotation   string    `orm:"column(rotation);size(255);null"`
	Scale      string    `orm:"column(scale);size(255);null"`
	PicRouteId string    `orm:"column(pic_route_id);size(255);null"`
	Material   string    `orm:"column(material);null"`
	Children   string    `orm:"column(children);null"`

	Using      int8      `orm:"column(using);null"`
}

func (t *Buildings) TableName() string {
	return "buildings"
}

