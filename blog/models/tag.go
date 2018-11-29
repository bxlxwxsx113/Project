package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

//标签(分类)
type Tag struct {
	Id int
	//标签名称
	Name string `orm:"size(20)"`
	//文章数量
	Count int
}

func (tag *Tag) TableName() string {
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "tag"
}

//插入
func (tag *Tag) Insert() error {
	if _, err := orm.NewOrm().Insert(tag); err != nil {
		return err
	}
	return nil
}

//查询
func (tag *Tag) Read(fields ...string) error {
	if err := orm.NewOrm().Read(tag, fields...); err != nil {
		return err
	}
	return nil
}

//更新
func (tag *Tag) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(tag, fields...); err != nil {
		return err
	}
	return nil
}

/*
思路：根据需要被删除的标签的id在标签文章表中找到相关记录，从这些记录中获取对应的
文章id，然后根据文章id在文章表中找到对应的记录，将标签名替换成逗号，最后删除标签文章表中的记录
*/
func (tag *Tag) Delete() error {
	//创建标签文章切片
	var list []*TagPost
	//获得文章表的表名
	table := new(Post).TableName()
	orm.NewOrm().QueryTable(new(TagPost)).Filter("tagid", tag.Id).All(&list)
	//判断查询结果中是否有值
	if len(list) > 0 {
		//创建切片，用于存储文章id
		ids := make([]string, 0, len(list))
		//遍历list
		for _, tagpost := range list {
			//将文章id存储到ids中
			ids = append(ids, strconv.Itoa(tagpost.Postid))
		}
		//[11,17,18]
		//UPDATE tb_post SET tags = REPLACE(tags, ':', ',') WHERE id IN(11, 17, 18);
		orm.NewOrm().Raw("UPDATE "+table+" set tags = REPLACE(tags, ?, ',') where id in ("+
			strings.Join(ids, ",")+")", ","+tag.Name+",").Exec()
		sql := "UPDATE " + table + " set tags = REPLACE(tags, ?, ',') where id in (" + strings.Join(ids, ",") + ")"
		fmt.Println("sql = ", sql)
		//删除中间表的记录
		orm.NewOrm().QueryTable(&TagPost{}).Filter("tagid", tag.Id).Delete()
	}
	//删除标签
	if _, err := orm.NewOrm().Delete(tag); err != nil {
		return err
	}
	return nil
}

func (tag *Tag) MergeTo(to *Tag) {
	//创建标签文章切片
	var list []*TagPost
	query := orm.NewOrm().QueryTable(new(TagPost))
	//在标签文章表中过滤出和原标签相关的记录
	query.Filter("tagid", tag.Id).All(&list)
	//判断list中是否有值
	if len(list) > 0 {
		//创建切片，用于存储文章id
		ids := make([]string, 0, len(list))
		//遍历list
		for _, v := range list {
			//取出每一文章id并追加到ids中
			ids = append(ids, strconv.Itoa(v.Postid))
		}
		//在标签文章表中将原标签id更新为目标标签的id
		query.Filter("tagid", tag.Id).Update(orm.Params{"tagid": to.Id})
		//UPDATE tb_post SET tags = REPLACE(tags, ':', ',') WHERE id IN(11, 17, 18);
		//获取文章表的表名
		table := new(Post).TableName()
		fmt.Println("sql = ", "UPDATE "+table+" set tags = REPLACE(tags, ?, ?) where id in ("+
			strings.Join(ids, ",")+")")
		orm.NewOrm().Raw("UPDATE "+table+" set tags = REPLACE(tags, ?, ?) where id in ("+
			strings.Join(ids, ",")+")", ","+tag.Name+",", ","+to.Name+",").Exec()
	}
}

func (tag *Tag) UpCount() {
	//获取标签中文章的数量
	count, err := orm.NewOrm().QueryTable(&TagPost{}).Filter("tagid", tag.Id).Count()
	newcount := int(count)
	//查询没有出错并且查询出来的文章数量和原数量不等
	if err == nil && newcount != tag.Count {
		tag.Count = newcount
		tag.Update("count")
	}
}
