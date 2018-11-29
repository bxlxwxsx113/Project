package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id int
	//用户名
	Username string `orm:"size(15)"`
	//密码
	Password string `orm:"size(32)"`
	//邮箱
	Email string `orm:"size(50)"`
	//登录次数
	Logincount int
	Authkey    string `orm:"size(10)"`
	//是否激活
	Active int
}

func (user *User) TableName() string {
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "user"
}

//插入
func (user *User) Insert() error {
	if _, err := orm.NewOrm().Insert(user); err != nil {
		return err
	}
	return nil
}

//删除
func (user *User) Delete() error {
	if _, err := orm.NewOrm().Delete(user); err != nil {
		return err
	}
	return nil
}

//查询
func (user *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(user, fields...); err != nil {
		return err
	}
	return nil
}

//更新
func (user *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(user, fields...); err != nil {
		return err
	}
	return nil
}
