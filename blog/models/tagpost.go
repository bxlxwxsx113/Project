package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

//标签文章表(中间表)
type TagPost struct {
	Id int
	//标签id
	Tagid int
	//文章id
	Postid int
	//文章的状态
	Poststatus int
	//发布时间
	Posttime time.Time `orm:"type(datetime)"`
}

func (tagpost *TagPost) TableName() string {
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "tag_post"
}

//插入
func (tagpost *TagPost) Insert() error {
	if _, err := orm.NewOrm().Insert(tagpost); err != nil {
		return err
	}
	return nil
}

//删除
func (tagpost *TagPost) Delete() error {
	if _, err := orm.NewOrm().Delete(tagpost); err != nil {
		return err
	}
	return nil
}

//查询
func (tagpost *TagPost) Read(fields ...string) error {
	if err := orm.NewOrm().Read(tagpost, fields...); err != nil {
		return err
	}
	return nil
}

//更新
func (tagpost *TagPost) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(tagpost, fields...); err != nil {
		return err
	}
	return nil
}
