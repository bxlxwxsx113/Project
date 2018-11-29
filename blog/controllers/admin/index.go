package admin

import (
	"blog/models"
	"github.com/astaxie/beego/orm"
	"os"
	"runtime"
)

type IndexController struct {
	baseController
}

//后台首页
func (this *IndexController) Index() {
	// hostname  gover   os  arch  cpunum  postnum  tagnum  usernum
	//主机名称
	this.Data["hostname"], _ = os.Hostname()
	//Go语言版本
	this.Data["gover"] = runtime.Version()
	//操作系统
	this.Data["os"] = runtime.GOOS
	//处理器架构
	this.Data["arch"] = runtime.GOARCH
	//cpu数量
	this.Data["cpunum"] = runtime.NumCPU()
	//文章数量
	this.Data["postnum"], _ = orm.NewOrm().QueryTable(new(models.Post)).Count()
	//分类数量
	this.Data["tagnum"], _ = orm.NewOrm().QueryTable(new(models.Tag)).Count()
	//用户数量
	this.Data["usernum"], _ = orm.NewOrm().QueryTable(new(models.User)).Count()
	this.display()
}
