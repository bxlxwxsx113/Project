package admin

import (
	"blog/models"
	"github.com/astaxie/beego/orm"
	"strings"
)

type LinkController struct {
	baseController
}

func (this *LinkController) List() {
	//创建切片，用于存储查询结果
	var list []*models.Link
	orm.NewOrm().QueryTable(new(models.Link)).OrderBy("-rank").All(&list)
	this.Data["list"] = list
	this.display()
}

func (this *LinkController) Add() {
	if this.Ctx.Request.Method == "POST" {
		//sitename  url  rank
		//获取网站的名称
		sitename := this.GetString("sitename")
		//获取网站的url
		url := this.GetString("url")
		//获取网站的排序值
		rank, err := this.GetInt("rank")
		if err != nil {
			rank = 0
		}
		//创建模型并初始化
		var link = &models.Link{Sitename: sitename, Url: url, Rank: rank}
		//插入数据库
		if err = link.Insert(); err != nil {
			this.showmsg(err.Error())
		}
		this.Redirect("/admin/link/list", 302)
	}
	this.display()
}

//删除友链
func (this *LinkController) Delete() {
	//获取友链id
	id, err := this.GetInt("id")
	if err != nil {
		this.showmsg("删除失败!")
	}
	//创建友链结构体并初始化id
	link := &models.Link{Id: id}
	//通过id查询友链
	if err = link.Read(); err == nil {
		//删除友链
		link.Delete()
	}
	this.Redirect("/admin/link/list", 302)
}

//编辑友链
func (this *LinkController) Edit() {
	//获取友情链接id
	id, _ := this.GetInt("id")
	//创建友情链接结构体并初始化id
	link := &models.Link{Id: id}
	if err := link.Read(); err != nil {
		this.showmsg("友链不存在，无法删除!")
	}
	if this.Ctx.Request.Method == "POST" {
		//获取网站名称
		sitename := strings.TrimSpace(this.GetString("sitename"))
		//获取网站地址
		url := strings.TrimSpace(this.GetString("url"))
		//获取排序等级
		rank, err := this.GetInt("rank")
		if err != nil {
			rank = 0
		}
		link.Sitename = sitename
		link.Url = url
		link.Rank = rank
		if err = link.Update(); err != nil {
			this.showmsg("更新失败，稍后再试!")
		}
		this.Redirect("/admin/link/list", 302)
	}
	this.Data["link"] = link
	this.display()
}
