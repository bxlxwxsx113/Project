package models

type Auth struct {
	Id       int
	Pid      int    //上级ID
	AuthName string //权限名称
	AuthUrl  string //URL地址
	Icon     string
}

func (auth *Auth) TableName() string {
	return TableName("uc_auth")
}
