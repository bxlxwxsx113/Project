package models

type InfoClass struct {
	Id        int
	ClassName string //分类名称
}

func (infoclass *InfoClass) TableName() string {
	return TableName("info_class")
}
