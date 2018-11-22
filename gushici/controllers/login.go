package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"gushici/libs"
	"gushici/models"
	"strconv"
	"strings"
	"time"
)

type LoginController struct {
	BaseController
}

func (this *LoginController) LoginIn() {
	//用户如果已经登录则重定向到后台首页
	if this.userId > 0 {
		//  /home
		this.redirect(beego.URLFor("HomeController.Index"))
	}
	fmt.Println("=============", beego.URLFor("HomeController.Index"), "=================")
	if this.isPost() {
		//获取用户名并去除两边的空格
		username := strings.TrimSpace(this.GetString("username"))
		//获取密码并去除两边的空格
		password := strings.TrimSpace(this.GetString("password"))
		if username != "" && password != "" {
			//根据用户名查询管理员
			user, err := models.AdminGetByName(username)
			//查询没有出错并且密码正确以及状态正常
			if err == nil && user.Password == libs.Md5([]byte(password+user.Salt)) && user.Status == 1 {
				//获取用户的ip
				user.LastIp = this.getClientIp()
				//获取登录时间
				user.LastLogin = time.Now().UnixNano()
				//更新用户
				user.Update()
				authkey := libs.Md5([]byte(this.getClientIp() + "|" + user.Password + user.Salt))
				//创建cookie,设置存活时间为一周
				this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 60*60*24*7)
				this.redirect(beego.URLFor("HomeController.Index"))
			}
			this.redirect(beego.URLFor("LoginController.LoginIn"))
		}
	}

	this.TplName = "login/login.html"
}
