package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

//文章
type Post struct {
	Id int
	//用户id
	Userid int
	//作者
	Author string `orm:"size(15)"`
	//标题
	Title string `orm:"size(100)"`
	//标题颜色
	Color string `orm:"size(7)"`
	//文章内容
	Content string `orm:"type(text)"`
	//所属标签
	Tags string `orm:"size(100)"`
	//浏览量
	Views int
	//文章状态
	Status int
	//发布时间
	Posttime time.Time `orm:"type(datetime)"`
	//更新时间
	Updated time.Time `orm:"type(datetime)"`
	//是否置顶
	Istop int
	//封面
	Cover string `orm:"size(70)"`
}

func (post *Post) TableName() string {
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "post"
}

//插入
func (post *Post) Insert() error {
	if _, err := orm.NewOrm().Insert(post); err != nil {
		return err
	}
	return nil
}

//查询
func (post *Post) Read(fields ...string) error {
	if err := orm.NewOrm().Read(post, fields...); err != nil {
		return err
	}
	return nil
}

//更新
func (post *Post) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(post, fields...); err != nil {
		return err
	}
	return nil
}

func (this *Post) Link() string {
	//   /article/2
	return "/article/" + strconv.Itoa(this.Id)
}

//设置标题内容和颜色
func (this *Post) ColorTitle() string {
	if this.Color != "" {
		//  <span style="color=''">this.Title</span>
		return fmt.Sprintf("<span style='color:%s'>%s</span>", this.Color, this.Title)
	}
	return this.Title
}

func (this *Post) Excerpt() string {
	return this.Content
}

func (this *Post) TagsLink() string {
	if this.Tags == "" {
		return ""
	}
	//去除标签前后的逗号
	return strings.Trim(this.Tags, ",")
}

//获取上一篇文章和下一篇文章
func (this *Post) GetPreAndNext() (pre, next *Post) {
	//创建文章结构体
	pre = &Post{}
	next = &Post{}
	//上一篇文章
	err := orm.NewOrm().QueryTable(new(Post)).OrderBy("-id").Filter("id__lt", this.Id).Filter("status", 0).Limit(1).One(pre)
	if err != nil {
		pre = nil
	}
	//下一篇文章
	err = orm.NewOrm().QueryTable(new(Post)).OrderBy("id").Filter("id__gt", this.Id).Filter("status", 0).Limit(1).One(next)
	if err != nil {
		next = nil
	}
	return
}

//删除文章
func (post *Post) Delete() error {
	//判断被删除的文章是否属于某一个分类
	if post.Tags != "" {
		orm := orm.NewOrm()
		//获得标签文章表的句柄并通过文章id过滤
		query := orm.QueryTable(&TagPost{}).Filter("postid", post.Id)
		//因为一篇文章可能属于多个分类，所以可能在标签文章表中查询出多条记录
		var tagpost []*TagPost
		//查询
		if n, err := query.All(&tagpost); n > 0 && err == nil {
			//遍历切片，取出每一条记录
			for i := 0; i < len(tagpost); i++ {
				//创建标签结构体并初始化标签id
				var tag = &Tag{Id: tagpost[i].Tagid}
				//数据库中有对应记录并且其count字段大于0
				if err = tag.Read(); err == nil && tag.Count > 0 {
					tag.Count--
					//更新count字段
					tag.Update("count")
				}
			}
		}
		//在中间表中删除原来的记录
		orm.QueryTable(&TagPost{}).Filter("postid", post.Id).Delete()
	}
	//删除文章
	if _, err := orm.NewOrm().Delete(post); err != nil {
		return err
	}
	return nil
}
