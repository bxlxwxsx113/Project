package admin

import (
	"blog/models"
	"github.com/astaxie/beego"
	"strings"
)

type baseController struct {
	beego.Controller
	userid         int    //用户id
	username       string //用户姓名
	controllerName string //控制器名称
	actionName     string //处理函数名称
	pager          *models.Pager
}

func (this *baseController) Prepare() {
	//获取控制器名称和处理函数的名称
	controllerName, actionName := this.GetControllerAndAction()
	//去除控制器尾部的controller，并将剩余部分转换为小写
	this.controllerName = strings.ToLower(controllerName[:len(controllerName)-10])
	//将处理函数名称转换为小写
	this.actionName = strings.ToLower(actionName)

	page, err := this.GetInt("page")
	if err != nil {
		page = 1
	}
	pagesize := 2
	this.pager = models.NewPager(page, pagesize, 0, "")
}

func (this *baseController) display(tplname ...string) {
	moduleName := "admin/"
	this.Layout = moduleName + "layout.html"
	if len(tplname) > 0 {
		//   admin/tag_list.html
		this.TplName = moduleName + tplname[0] + ".html"
	} else {
		this.TplName = moduleName + this.controllerName + "_" + this.actionName + ".html"
	}
}

func (this *baseController) showmsg(msg ...string) {
	if len(msg) == 0 {
		msg = append(msg, "出错了!")
	}
	this.Data["msg"] = msg[0]
	this.Data["redirect"] = this.Ctx.Request.Referer()
	this.display("showmsg")
	this.Render()
	this.StopRun()
}
