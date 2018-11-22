package libs

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

//获取哈希
func Md5(buf []byte) string {
	//创建哈希对象
	hash := md5.New()
	//把待哈希的数据写入哈希对象
	hash.Write(buf)
	//将哈希之后的结果以16进制字符串返回
	return fmt.Sprintf("%x", hash.Sum(nil))
}

//生成随机字符串
func GetRandomString(lens int) string {
	//准备母串 6
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	//存储每次获取到的随机字符
	result := []byte{} //sj
	//生成随机数对象
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//i指的是遍历的次数
	for i := 0; i < lens; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//返回密码加密之后的结果
func Password(len int, pwd0 string) (pwd, salt string) {
	//生成长度为4的随机字符串
	salt = GetRandomString(len)
	//默认密码
	defaultPwd := "george518"
	//判断用户是否传递了密码
	if pwd0 != "" {
		//将默认密码修改为用户传递过来的密码
		defaultPwd = pwd0
	}
	//将密码和密码盐拼接并做哈希处理
	pwd = Md5([]byte(defaultPwd + salt))
	return pwd, salt
}
