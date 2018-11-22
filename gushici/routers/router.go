package routers

import (
	"github.com/astaxie/beego"
	"gushici/controllers"
)

func init() {
	//首页
	beego.Router("/", &controllers.WwwController{}, "*:Index")
	//古诗详情   /show/{{$elem.id}}   /show/2
	beego.Router("/show/:id", &controllers.WwwController{}, "*:Show")
	//分类查询  /list/2
	beego.Router("/list/:class_id", &controllers.WwwController{}, "*:List")
	//登录   http://localhost:8080/login
	beego.Router("/login", &controllers.LoginController{}, "*:LoginIn")

	beego.Router("/home", &controllers.HomeController{}, "*:Index")
}
