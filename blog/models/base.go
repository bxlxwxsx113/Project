package models

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	//"root:111111@tcp(127.0.0.1:3306)/HelloBeego?charset=utf8"
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"

	fmt.Println("dsn = ", dsn)

	orm.RegisterDataBase("default", "mysql", dsn, 30)

	// register mode
	orm.RegisterModel(new(Link), new(Mood), new(Post), new(Tag), new(TagPost), new(User))
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
}

//获取最新的4篇文章
func GetLatestBlog() []*Post {
	//创建切片，用于存储查询结果
	var result []*Post
	post := Post{}
	//从tb_post表中过滤出状态正常的文章
	query := orm.NewOrm().QueryTable(&post).Filter("status", 0)
	//获取满足条件的文章数量
	count, _ := query.Count()
	//判断count是否大于0
	if count > 0 {
		//如果数量大于0，则通过文章发表时间降序排序查找4篇最新文章
		query.OrderBy("-posttime").Limit(4).All(&result)
	}
	return result
}

//获取友情链接
func GetLinks() []*Link {
	var result []*Link
	link := Link{}
	query := orm.NewOrm().QueryTable(&link)
	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-rank").All(&result)
	}
	return result
}

//获取浏览量最高的4篇文章
func GetHotBlog() []*Post {
	//创建切片，用于存储查询结果
	var result []*Post
	post := Post{}
	//从tb_post表中过滤出状态正常的文章
	query := orm.NewOrm().QueryTable(&post).Filter("status", 0)
	//获取满足条件的文章数量
	count, _ := query.Count()
	//判断count是否大于0
	if count > 0 {
		//如果数量大于0，则通过文章发表时间降序排序查找4篇最新文章
		query.OrderBy("-views").Limit(4).All(&result)
	}
	return result
}

func Md5(buf []byte) string {
	mymd5 := md5.New()
	mymd5.Write(buf)
	result := mymd5.Sum(nil)
	return fmt.Sprintf("%x", result)
}
