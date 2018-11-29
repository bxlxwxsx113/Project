package admin

import (
	"blog/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"strings"
)

type UserController struct {
	baseController
}

//用户列表
func (this *UserController) List() {
	//创建切片
	var list []*models.User
	//获得tb_user表的句柄
	query := orm.NewOrm().QueryTable(new(models.User))
	count, _ := query.Count()
	if count > 0 {
		offset := (this.pager.Page - 1) * this.pager.Pagesize
		query.OrderBy("-id").Limit(this.pager.Pagesize, offset).All(&list)
	}
	//设置总数量
	this.pager.SetTotalnum(int(count))
	//设置rootpath
	this.pager.SetUrlpath("/admin/user/list?page=%d")
	this.Data["pagebar"] = this.pager.ToString()
	this.Data["list"] = list
	this.display()
}

//添加用户
func (this *UserController) Add() {
	//创建map,用于存储错误提示
	errmsg := make(map[string]string)
	//创建map，用于回显用户输入的数据
	input := make(map[string]string)
	if this.Ctx.Request.Method == "POST" {
		//username  password  password2   email
		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		password2 := strings.TrimSpace(this.GetString("password2"))
		email := strings.TrimSpace(this.GetString("email"))
		active, _ := this.GetInt("active")

		input["username"] = username
		input["password"] = password
		input["password2"] = password2
		input["email"] = email

		//创建Validation对象，用于校验数据是否合法
		valid := validation.Validation{}
		//校验用户名
		if result := valid.Required(username, "username"); !result.Ok {
			errmsg["username"] = "请输入用户名!"
		} else if result := valid.MaxSize(username, 15, "username"); !result.Ok {
			errmsg["username"] = "用户名长度不能大于15个字符!"
		}

		//校验密码
		if result := valid.Required(password, "password"); !result.Ok {
			errmsg["password"] = "请输入密码!"
		}

		//校验重复密码
		if result := valid.Required(password2, "password2"); !result.Ok {
			errmsg["password2"] = "请再次输入密码!"
		} else if password != password2 {
			errmsg["password2"] = "两次输入的密码不一致!"
		}

		//校验邮箱
		if result := valid.Required(email, "email"); !result.Ok {
			errmsg["email"] = "请输入电子邮箱!"
		} else if result := valid.Email(email, "email"); !result.Ok {
			errmsg["email"] = "电子邮箱非法!"
		}

		//错误提示为空，说明没有出现错误
		if len(errmsg) == 0 {
			//创建用户结构体
			var user = &models.User{Username: username, Email: email, Active: active}
			user.Password = models.Md5([]byte(password))
			if err := user.Insert(); err != nil {
				this.showmsg(err.Error())
			}
			this.Redirect("/admin/user/list", 302)
		}
	}
	this.Data["input"] = input
	this.Data["errmsg"] = errmsg
	this.display()
}

//删除用户
func (this *UserController) Delete() {
	//获取用户id
	id, _ := this.GetInt("id")
	//不能删除超级管理员
	if id == 7 {
		this.showmsg("不能删除超级管理员!")
	}
	//创建用户结构体，并初始化id
	user := &models.User{Id: id}
	//判断数据库中是否存在该用户
	if user.Read() == nil {
		//删除用户
		user.Delete()
	}
	this.Redirect("/admin/user/list", 302)
}

//编辑用户
func (this *UserController) Edit() {
	//获取用户id
	id, _ := this.GetInt("id")
	//创建用户结构体并初始化id
	user := &models.User{Id: id}
	if err := user.Read(); err != nil {
		this.showmsg("用户不存在!")
	}
	errmsg := make(map[string]string)
	if this.Ctx.Input.Method() == "POST" {
		//password   password2   email  active
		password := strings.TrimSpace(this.GetString("password"))
		password2 := strings.TrimSpace(this.GetString("password2"))
		email := strings.TrimSpace(this.GetString("email"))
		active, _ := this.GetInt("active")
		valid := validation.Validation{}
		if password != "" {
			if result := valid.Required(password2, "password2"); !result.Ok {
				errmsg["password2"] = "请再次输入密码!"
			} else if password != password2 {
				errmsg["password2"] = "两次输入的面不一致!"
			} else {
				//设置密码
				user.Password = models.Md5([]byte(password))
			}
		}

		if result := valid.Required(email, "email"); !result.Ok {
			errmsg["email"] = "请输入Email地址!"
		} else if result := valid.Email(email, "email"); !result.Ok {
			errmsg["email"] = "邮箱不合法!"
		} else {
			user.Email = email
		}

		user.Active = active
		if len(errmsg) == 0 {
			user.Update()
			this.Redirect("/admin/user/list", 302)
		}
	}
	this.Data["user"] = user
	this.Data["errmsg"] = errmsg
	this.display()

}
