package controllers

type HomeController struct {
	BaseController
}

func (this *HomeController) Index() {
	this.TplName = "public/main.html"
	//this.Ctx.WriteString("登录成功!")
}
