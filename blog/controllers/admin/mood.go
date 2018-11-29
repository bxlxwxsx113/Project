package admin

import (
	"blog/models"
	"fmt"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"time"
)

type MoodController struct {
	baseController
}

//添加说说
func (this *MoodController) Add() {
	//判断是否是post请求
	if this.Ctx.Request.Method == "POST" {
		//获取用户提交的说说内容
		content := this.GetString("content")
		//创建说说结构体
		var mood models.Mood
		mood.Content = content
		//初始化随机数种子
		rand.Seed(time.Now().Unix())
		//生成[0,11)之间的随机数
		var r = rand.Intn(11)
		mood.Cover = "/static/upload/blog" + fmt.Sprintf("%d", r) + ".jpg"
		//设置说说发布时间
		mood.Posttime = time.Now()
		//插入数据库
		if err := mood.Insert(); err != nil {
			this.showmsg(err.Error())
		}
		this.Redirect("/admin/mood/list", 302)
	}
	this.display()
}

//说说列表
func (this *MoodController) List() {
	//创建说收切片，用于存储查询结果
	var list []*models.Mood
	//获得tb_mood表的句柄
	query := orm.NewOrm().QueryTable(new(models.Mood))
	count, _ := query.Count()
	if count > 0 {
		//计算偏移量
		offset := (this.pager.Page - 1) * this.pager.Pagesize
		//分页查询
		query.OrderBy("-id").Limit(this.pager.Pagesize, offset).All(&list)
	}
	this.pager.SetTotalnum(int(count))
	this.pager.SetUrlpath("/admin/mood/list?page=%d")
	//设置分页导航栏
	this.Data["pagebar"] = this.pager.ToString()
	this.Data["list"] = list
	this.display()
}

//删除说说
func (this *MoodController) Delete() {
	//获取说说id
	id, err := this.GetInt("id")
	//判断是否出现错误
	if err != nil {
		this.showmsg("删除失败!")
	}
	//创建说说结构体，并使用id进行初始化
	mood := models.Mood{Id: id}
	//通过id查询说说，没有出现错误，则删除
	if err = mood.Read(); err == nil {
		mood.Delete()
	}
	this.Redirect("/admin/mood/list", 302)
}
