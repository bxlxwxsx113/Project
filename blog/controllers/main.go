package controllers

import (
	"blog/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type MainController struct {
	beego.Controller
	Pager *models.Pager
}

func (this *MainController) Prepare() {
	var page int
	var err error

	if page, err = strconv.Atoi(this.Ctx.Input.Param(":page")); err != nil {
		page = 1
	}

	pagesize := 2
	//   /index%d.html
	this.Pager = models.NewPager(page, pagesize, 0, "")
}

//首页
func (this *MainController) Index() {
	//创建文章切片，用于存储查询结果
	var list []*models.Post
	//创建文章结构体
	post := models.Post{}
	//获得tb_post表的句柄
	query := orm.NewOrm().QueryTable(&post).Filter("status", 0)
	//获得符合条件的记录数量
	count, _ := query.Count()

	//判断符合条件的记录的数量是否大于0
	if count > 0 {
		//计算偏移量
		offset := (this.Pager.Page - 1) * this.Pager.Pagesize
		//按是否置顶，浏览量降序排序
		_, err := query.OrderBy("-istop", "-views").Limit(this.Pager.Pagesize, offset).All(&list)
		if err != nil {
			this.Redirect("/404", 302)
		}
	}
	this.Data["list"] = list

	this.Pager.SetTotalnum(int(count))
	this.Pager.SetUrlpath("/index%d.html")

	this.Data["pagebar"] = this.Pager.ToString()

	//设置头部公共信息
	this.setHeadMeater()
	//设置右侧公共信息
	this.setRight()
	this.display("index")

}

//通过文章id查询文章
func (this *MainController) Show() {
	//获取文章id
	id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err != nil {
		this.Redirect("/404", 302)
	}
	//创建文章结构体
	post := new(models.Post)
	post.Id = id
	//按照文章id查询
	err = post.Read()
	if err != nil {
		this.Redirect("/404", 302)
	}
	//浏览量加一
	post.Views++
	//更新浏览量
	post.Update("views")
	this.Data["post"] = post
	this.Data["smalltitle"] = "文章详情"
	//获取上一篇文章和下一篇文章
	pre, next := post.GetPreAndNext()
	this.Data["pre"] = pre
	this.Data["next"] = next

	this.setHeadMeater()
	this.display("article")
}

func (this *MainController) display(tplname string) {
	theme := "double"
	//布局
	this.Layout = theme + "/layout.html"
	//模版名称
	this.TplName = theme + "/" + tplname + ".html"
	//创建map，用于存储布局子木块
	this.LayoutSections = make(map[string]string)
	//页头
	this.LayoutSections["head"] = theme + "/head.html"
	//页脚
	this.LayoutSections["foot"] = theme + "/foot.html"

	if tplname == "index" {
		this.LayoutSections["banner"] = theme + "/banner.html"
		this.LayoutSections["middle"] = theme + "/middle.html"
		this.LayoutSections["right"] = theme + "/right.html"
	} else if tplname == "life" {
		this.LayoutSections["right"] = theme + "/right.html"
	}
}

//设置头部公共信息
func (this *MainController) setHeadMeater() {
	this.Data["title"] = beego.AppConfig.String("title")
	this.Data["keywords"] = beego.AppConfig.String("keywords")
	this.Data["description"] = beego.AppConfig.String("description")
}

//设置右侧公共信息
func (this *MainController) setRight() {
	//查询最新的4篇文章
	this.Data["latestblog"] = models.GetLatestBlog()
	//查询点击量最高的4篇文章
	this.Data["hotblog"] = models.GetHotBlog()
	//查询友情链接
	this.Data["links"] = models.GetLinks()
}

//关于我
func (this *MainController) About() {
	this.setHeadMeater()
	this.display("about")
}

//成长录
func (this *MainController) Life() {
	//创建文章结构体，用于存储查询结果
	var list []*models.Post
	query := orm.NewOrm().QueryTable(new(models.Post)).Filter("status", 0)

	//获取数量
	count, _ := query.Count()
	if count > 0 {
		//计算偏移量
		offset := (this.Pager.Page - 1) * this.Pager.Pagesize
		query.OrderBy("-istop", "-posttime").Limit(this.Pager.Pagesize, offset).All(&list)
	}
	this.Data["list"] = list
	this.Pager.SetTotalnum(int(count))

	//   "/index%d.html"
	this.Pager.SetUrlpath("/life%d.html")
	this.Data["pagebar"] = this.Pager.ToString()
	this.setHeadMeater()
	this.setRight()
	this.display("life")
}

func (this *MainController) Mood() {
	//创建切片，用于存储查询结果
	var list []*models.Mood
	//获得tb_post表的句柄
	query := orm.NewOrm().QueryTable(new(models.Mood))
	count, _ := query.Count()
	if count > 0 {
		//计算偏移量
		offset := (this.Pager.Page - 1) * this.Pager.Pagesize
		query.OrderBy("-posttime").Limit(this.Pager.Pagesize, offset).All(&list)
	}
	this.Data["list"] = list
	//设置总数量
	this.Pager.SetTotalnum(int(count))
	//设置rootpath
	this.Pager.SetUrlpath("/mood%d.html")
	//设置分页导航栏
	this.Data["pagebar"] = this.Pager.ToString()
	//设置头部公共信息
	this.setHeadMeater()
	this.display("mood")
}
