package controllers

import (
	"github.com/astaxie/beego"
	"gushici/libs"
	"gushici/models"
	"strconv"
	"strings"
)

type BaseController struct {
	beego.Controller
	//控制器名称
	controllerName string
	//方法名称
	actionName string
	//每页显示的记录的数量
	pageSize int
	//管理员ID
	userId int
	//用户名
	userName string
	//登录名
	loginName string
	//noLayout bool
}

func (this *BaseController) Prepare() {
	//WwwController/Index
	//获取控制器和方法名称
	controllerName, actionName := this.GetControllerAndAction()
	//去除控制器尾部的controller并将剩余部分装换为小写
	this.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10]) //www
	//将方法名称装换为小写
	this.actionName = strings.ToLower(actionName) //index
	noAuth := "www"
	//判断用户请求的是前台页面还是后台页面
	isNoAuth := strings.Contains(noAuth, this.controllerName)
	if !isNoAuth {
		this.auth()
	}
}

//验证登录权限
func (this *BaseController) auth() {
	/*authkey := libs.Md5([]byte(this.getClientIp() + "|" + user.Password + user.Salt))
	//创建cookie,设置存活时间为一周
	this.Ctx.SetCookie("auth", strconv.Itoa(user.Id) + "|" + authkey, 60*60*24*7)*/

	//获取名为auth的cookie并通过|切割
	arr := strings.Split(this.Ctx.GetCookie("auth"), "|")
	this.userId = 0
	if len(arr) == 2 {
		//分别取出用户id和authkey
		idstr, authkey := arr[0], arr[1]
		//将用户id装换为整形
		userId, _ := strconv.Atoi(idstr)
		//判断用户id是否大于0
		if userId > 0 {
			//通过用户id查询管理员
			user, err := models.AdminGetById(userId)
			//在查询过程中没有出现错误并且authkey正确
			if err == nil && authkey == libs.Md5([]byte(user.LastIp+"|"+user.Password+user.Salt)) {
				this.userId = user.Id
				this.loginName = user.LoginName
				this.userName = user.RealName

			}
		}
	}
	if this.userId == 0 && (this.controllerName != "login" && this.actionName != "loginin") {
		this.redirect(beego.URLFor("LoginController.LoginIn"))
	}
}

func (this *BaseController) display(tpl ...string) {
	var tplname string
	//判断tpl中是否有值
	if len(tpl) > 0 {
		//index   index.html
		//在模板名称后面拼接.hml
		tplname = strings.Join([]string{tpl[0], "html"}, ".")
	} else {
		//通过控制器和方法名称拼接页面的路径
		//   www/index.html
		tplname = this.controllerName + "/" + this.actionName + ".html"
	}
	/*if !this.noLayout {
		this.Layout = "public/layout.html"
	}*/

	this.TplName = tplname
}

//重定向
func (this *BaseController) redirect(url string) {
	this.Redirect(url, 302)
}

//判断是否是post提交
func (this *BaseController) isPost() bool {
	return this.Ctx.Request.Method == "POST"
}

//获取用户ip
func (this *BaseController) getClientIp() string {
	//使用冒号切割获取到的ip地址和端口号
	s := strings.Split(this.Ctx.Request.RemoteAddr, ":") //127.0.0.1:8080
	return s[0]
}
