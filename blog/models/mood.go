package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

//说说
type Mood struct {
	Id int
	//说说内容
	Content string `orm:"type(text)"`
	//封面
	Cover string `orm:"size(70)"`
	//发布时间
	Posttime time.Time `orm:"type(datetime)"`
}

func (mood *Mood) TableName() string {
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "mood"
}

//插入
func (mood *Mood) Insert() error {
	if _, err := orm.NewOrm().Insert(mood); err != nil {
		return err
	}
	return nil
}

//删除
func (mood *Mood) Delete() error {
	if _, err := orm.NewOrm().Delete(mood); err != nil {
		return err
	}
	return nil
}

//查询
func (mood *Mood) Read(fields ...string) error {
	if err := orm.NewOrm().Read(mood, fields...); err != nil {
		return err
	}
	return nil
}

//更新
func (mood *Mood) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(mood, fields...); err != nil {
		return err
	}
	return nil
}
