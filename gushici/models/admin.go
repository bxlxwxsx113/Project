package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Admin struct {
	Id         int
	LoginName  string //登录名
	RealName   string //真实姓名
	Password   string //密码
	Phone      string //手机号码
	Email      string //邮箱
	Salt       string //密码盐
	LastLogin  int64  //最后登录时间
	LastIp     string //最后登录的ip地址
	Status     int    //状态
	CreateId   int    //创建者id
	UpdateId   int    //修改者id
	CreateTime int64  //创建时间
	UpdateTime int64  //修改时间
}

//获取表名
func (admin *Admin) TableName() string {
	return TableName("uc_admin")
}

//根据用户名查询管理员
func AdminGetByName(loginName string) (*Admin, error) {
	admin := new(Admin)
	err := orm.NewOrm().QueryTable(TableName("uc_admin")).Filter("login_name", loginName).One(admin)
	if err != nil {
		fmt.Println("err = ", err)
		return nil, err
	}
	return admin, nil
}

//更新管理员表
func (admin *Admin) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(admin, fields...); err != nil {
		return err
	}
	return nil
}

//根据id查询管理员
func AdminGetById(id int) (*Admin, error) {
	admin := new(Admin)
	err := orm.NewOrm().QueryTable(TableName("uc_admin")).Filter("id", id).One(admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}
