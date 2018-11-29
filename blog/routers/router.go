package routers

import (
	"blog/controllers"
	"blog/controllers/admin"
	"github.com/astaxie/beego"
)

func init() {
	// http://localhost:8080
	beego.Router("/", &controllers.MainController{}, "*:Index")

	//  /index3.html
	beego.Router("/index:page:int.html", &controllers.MainController{}, "*:Index")

	//  /article/2
	beego.Router("/article/:id:int", &controllers.MainController{}, "*:Show")

	//关于我
	beego.Router("/about.html", &controllers.MainController{}, "*:About")

	//成长录
	beego.Router("/life.html", &controllers.MainController{}, "*:Life")
	//  /life2.html  成长录分页路由
	beego.Router("/life:page:int.html", &controllers.MainController{}, "*:Life")

	//碎言碎语
	beego.Router("/mood.html", &controllers.MainController{}, "*:Mood")
	//碎言碎语分页路由
	beego.Router("/mood:page:int.html", &controllers.MainController{}, "*:Mood")

	//--------------------------------------后台------------------------------------------------
	//首页
	beego.Router("/admin", &admin.IndexController{}, "*:Index")

	//------------------------------说说管理----------------------------------
	//添加说说
	beego.Router("/admin/mood/add", &admin.MoodController{}, "*:Add")
	//说说列表
	beego.Router("/admin/mood/list", &admin.MoodController{}, "*:List")
	//删除说说
	beego.Router("/admin/mood/delete", &admin.MoodController{}, "*:Delete")

	//------------------------------友链管理----------------------------------
	//友链列表
	beego.Router("/admin/link/list", &admin.LinkController{}, "*:List")
	//添加友链
	beego.Router("/admin/link/add", &admin.LinkController{}, "*:Add")
	//删除友链
	beego.Router("/admin/link/delete", &admin.LinkController{}, "*:Delete")
	//编辑友链
	beego.Router("/admin/link/edit", &admin.LinkController{}, "*:Edit")

	//------------------------------用户管理----------------------------------
	//用户列表
	beego.Router("/admin/user/list", &admin.UserController{}, "*:List")
	//添加用户
	beego.Router("/admin/user/add", &admin.UserController{}, "*:Add")
	//删除用户
	beego.Router("/admin/user/delete", &admin.UserController{}, "*:Delete")
	//编辑用户
	beego.Router("/admin/user/edit", &admin.UserController{}, "*:Edit")

	//------------------------------文章管理----------------------------------
	//
	beego.Router("/admin/article/add", &admin.ArticleController{}, "*:Add")
	beego.Router("/admin/article/save", &admin.ArticleController{}, "*:Save")
	//文章列表
	beego.Router("/admin/article/list", &admin.ArticleController{}, "*:List")
	//删除文章
	beego.Router("/admin/article/delete", &admin.ArticleController{}, "*:Delete")
	//编辑文章(跳转到编辑页面)
	beego.Router("/admin/article/edit", &admin.ArticleController{}, "*:Edit")
	//编辑文章
	beego.Router("/admin/article/update", &admin.ArticleController{}, "*:Update")
	//文章批量操作
	beego.Router("/admin/article/batch", &admin.ArticleController{}, "*:Batch")
	//------------------------------标签管理----------------------------------
	//标签列表
	beego.Router("/admin/tag", &admin.TagController{}, "*:List")
}
