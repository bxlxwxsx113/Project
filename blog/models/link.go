package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//友情链接
type Link struct {
	Id int
	//网站名称
	Sitename string `orm:"size(80)"`
	//网址
	Url string `orm:"size(200)"`
	//排序值，越大越靠前
	Rank int
}

func (link *Link) TableName() string {
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "link"
}

//插入
func (link *Link) Insert() error {
	if _, err := orm.NewOrm().Insert(link); err != nil {
		return err
	}
	return nil
}

//删除
func (link *Link) Delete() error {
	if _, err := orm.NewOrm().Delete(link); err != nil {
		return err
	}
	return nil
}

//查询
func (link *Link) Read(fields ...string) error {
	if err := orm.NewOrm().Read(link, fields...); err != nil {
		return err
	}
	return nil
}

//更新
func (link *Link) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(link, fields...); err != nil {
		return err
	}
	return nil
}
