package admin

import (
	"blog/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type TagController struct {
	baseController
}

func (this *TagController) List() {

	act := this.GetString("act")
	switch act {
	//批处理
	case "batch":
		this.batch()
	default:
		this.TagList()
	}
}

func (this *TagController) batch() {
	//获取用户所选择的id
	ids := this.GetStrings("ids[]")
	//获取用户所选择的操作
	op := this.GetString("op")
	idarr := make([]int, 0)
	for _, v := range ids {
		//取出每一个id并转换为整数
		if id, _ := strconv.Atoi(v); id > 0 {
			idarr = append(idarr, id)
		}
	}
	switch op {
	//合并
	/*
		1.获取用户输入的目标标签名称，并去除两边的空格，根据该名称去标签表中查找记录，
		如果没有查找到则需要插入新的记录
		2.根据标签id在标签文章表中查找对应记录，将对应的记录标签id更新为目标标签的id
		3.在文章表中将原标签替换成目标标签的名称
		4.更新目标标签中的count字段
	*/
	case "merge":
		//获取目标标签的名称并去除两边的空格
		toname := strings.TrimSpace(this.GetString("toname"))
		if toname != "" && len(idarr) > 0 {
			//创建标签结构体
			tag := new(models.Tag)
			//设置标签名称
			tag.Name = toname
			//在标签表中根据标签名称查询
			if tag.Read("name") != nil {
				tag.Count = 0
				//插入标签
				tag.Insert()
			}
			//遍历idarr
			for _, id := range idarr {
				obj := models.Tag{Id: id}
				//根据标签id查询数据库
				if obj.Read() == nil {
					obj.MergeTo(tag)
					obj.Delete()
				}
			}
			//更新新标签
			tag.UpCount()
		}

	//删除
	case "delete":
		for _, id := range idarr {
			//创建标签结构体并初始化id
			obj := models.Tag{Id: id}
			if obj.Read() == nil {
				obj.Delete()
			}
		}
	}
	this.Redirect("/admin/tag", 302)
}

func (this *TagController) TagList() {
	//创建切片
	var list []*models.Tag
	//获得标签表的句柄
	query := orm.NewOrm().QueryTable(new(models.Tag))
	//获取标签的数量
	count, _ := query.Count()
	//数量大于0则查询
	if count > 0 {
		offset := (this.pager.Page - 1) * this.pager.Pagesize
		query.Limit(this.pager.Pagesize, offset).All(&list)
	}
	//设置总数量
	this.pager.SetTotalnum(int(count))
	//设置rootpath
	this.pager.SetUrlpath("/admin/tag?page=%d")
	this.Data["pagebar"] = this.pager.ToString()
	this.Data["list"] = list
	this.display("tag_list")
}
