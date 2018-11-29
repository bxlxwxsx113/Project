package admin

import (
	"blog/models"
	"fmt"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type ArticleController struct {
	baseController
}

//跳转到保存页面
func (this *ArticleController) Add() {
	this.display()
}

//保存文章
/*
1.获取用户输入的文章信息，插入数据库
2.用户有可能会输入多个标签，所以我们可以取出每个标签并去除两边的空格，去除重复数据
3.判断标签表中是否存在这些标签，如果不存在要插入标签，否则更新count字段
4.在标签文章表中插入相应的记录
*/
func (this *ArticleController) Save() {
	//创建文章结构体
	var post models.Post
	// title  color  istop  tags  posttime  status  content
	//获取文章标题
	post.Title = strings.TrimSpace(this.GetString("title"))
	if post.Title == "" {
		this.showmsg("请输入标题!")
	}
	//获取标题颜色
	post.Color = strings.TrimSpace(this.GetString("color"))
	//文章是否置顶
	post.Istop, _ = this.GetInt("istop")
	//获取文章所属的标签
	tags := strings.TrimSpace(this.GetString("tags"))
	//获取文章的发布时间
	timestr := strings.TrimSpace(this.GetString("posttime"))
	//文章状态
	post.Status, _ = this.GetInt("status")
	//文章内容
	post.Content = this.GetString("content")
	post.Userid = this.userid
	post.Author = this.username
	//初始化随机数种子
	rand.Seed(time.Now().Unix())
	var r = rand.Intn(11)
	post.Cover = "/static/upload/blog" + fmt.Sprintf("%d", r) + ".jpg"
	//Mon Jan 2 15:04:05 -0700 MST 2006
	posttime, err := time.Parse("2006-01-02 15:04:05", timestr)
	if err == nil {
		post.Posttime = posttime
		post.Updated = posttime
	}
	//插入文章
	if err = post.Insert(); err != nil {
		this.showmsg("文章添加失败!")
	}

	//存储最终结果
	addtags := make([]string, 0)
	if tags != "" {
		tagarr := strings.Split(tags, ",")
		for _, v := range tagarr {
			if tag := strings.TrimSpace(v); tags != "" {
				exists := false
				for _, vv := range addtags {
					if vv == tag {
						exists = true
						break
					}
				}
				if !exists {
					addtags = append(addtags, tag)
				}
			}
		}
	}
	if len(addtags) > 0 {
		for _, v := range addtags {
			tag := &models.Tag{Name: v}
			if err := tag.Read("Name"); err == orm.ErrNoRows {
				tag.Count = 1
				tag.Insert()
			} else {
				tag.Count += 1
				tag.Update("count")
			}
			tp := &models.TagPost{Tagid: tag.Id, Postid: post.Id, Poststatus: post.Status, Posttime: post.Posttime}
			tp.Insert()
		}
		post.Tags = "," + strings.Join(addtags, ",") + ","
	}
	post.Updated = time.Now()
	post.Update("tags", "updated")
	this.Redirect("/admin", 302)
}

func (this *ArticleController) List() {
	status, _ := this.GetInt("status")
	searchtype := this.GetString("searchtype")
	keyword := this.GetString("keyword")

	//创建切片
	var list []*models.Post
	query := orm.NewOrm().QueryTable(new(models.Post)).Filter("status", status)
	if keyword != "" { //go
		switch searchtype {
		case "title":
			// select * from tb_post where title like '%keyword%'
			query = query.Filter("title__icontains", keyword)
		case "author":
			query = query.Filter("author__icontains", keyword)
		case "tag":
			query = query.Filter("tags__icontains", keyword)
		}
	}

	count, _ := query.Count()
	if count > 0 {
		offset := (this.pager.Page - 1) * this.pager.Pagesize
		query.Limit(this.pager.Pagesize, offset).All(&list)
	}
	this.pager.SetTotalnum(int(count))
	this.pager.SetUrlpath(fmt.Sprintf("/admin/article/list?status=%d&searchtype=%s&keyword=%s&page=%s", status, searchtype, keyword, "%d"))
	this.Data["pagebar"] = this.pager.ToString()
	this.Data["list"] = list
	this.Data["status"] = status
	this.Data["count_1"], _ = orm.NewOrm().QueryTable(&models.Post{}).Filter("status", 1).Count()
	this.Data["count_2"], _ = orm.NewOrm().QueryTable(&models.Post{}).Filter("status", 2).Count()
	this.Data["searchtype"] = searchtype
	this.Data["keyword"] = keyword
	this.display()

}

//删除文章
func (this *ArticleController) Delete() {
	//获取文章id
	id, _ := this.GetInt("id")
	//创建文章结构体，并初始化id
	post := &models.Post{Id: id}
	//判断查询文章是否出错
	if post.Read() == nil {
		//删除文章
		post.Delete()
	}
	this.Redirect("/admin/article/list", 302)
}

//编辑文章(跳转到文章编辑页面)
func (this *ArticleController) Edit() {
	//获取待编辑的文章的id
	id, _ := this.GetInt("id")
	//创建文章结构体并初始化id
	post := &models.Post{Id: id}
	//判断查询文章是否出现错误
	if post.Read() != nil {
		this.showmsg("未找到该文章!")
	}
	//去除标签两边的逗号
	post.Tags = strings.Trim(post.Tags, ",")
	this.Data["post"] = post
	this.Data["posttime"] = post.Posttime.Format("2006-01-02 15:04:05")
	this.display()
}

/*
思路：
第一种情况:需要判断用户有没有修改标签，如果没有直接更新
第二种情况:如果修改了标签，则需要对应的更新标签文章表和标签表，
首先需要判断修改之前该文章的标签是否为空，如果不为空，则将标签文章表中的相关记录删除，
更新标签表中count字段，然后处理用户输入的标签，得到合法的标签之后和添加逻辑就一样了。
*/
func (this *ArticleController) Update() {
	//创建文章结构体
	var post models.Post
	//id   title   color  istop  tags  posttime  status  content
	//获取文章id
	id, err := this.GetInt("id")
	if err != nil {
		this.showmsg("更新失败!")
	}
	post.Id = id
	//如果数据库中不存在对应的记录，则重定向到文章列表页面
	if post.Read() != nil {
		this.Redirect("/admin/article/list", 302)
	}
	//文章标题
	post.Title = strings.TrimSpace(this.GetString("title"))
	//文章颜色
	post.Color = strings.TrimSpace(this.GetString("color"))
	//是否置顶
	post.Istop, _ = this.GetInt("istop")
	//文章所属标签
	tags := strings.TrimSpace(this.GetString("tags"))
	//去除标签两边的逗号
	tags = strings.Trim(tags, ",")
	//文章发布时间
	timestr := strings.TrimSpace(this.GetString("posttime"))
	//将字符串时间转换为time对象
	if posttime, err := time.Parse("2006-01-02 15:04:05", timestr); err == nil {
		post.Posttime = posttime
	}
	//文章状态
	post.Status, _ = this.GetInt("status")
	//文章内容
	post.Content = this.GetString("content")
	if strings.Trim(post.Tags, ",") == tags {
		post.Update("title", "color", "istop", "posttime", "status", "content")
		this.Redirect("/admin/article/list", 302)
	}

	if post.Tags != "" {
		//创建标签文章结构体
		var tagpost models.TagPost
		query := orm.NewOrm().QueryTable(&tagpost).Filter("postid", post.Id)
		//创建切片，存储查询结果
		var tagpostarr []*models.TagPost
		//查询没有出错并且有对应记录
		if n, err := query.All(&tagpostarr); n > 0 && err == nil {
			for i := 0; i < len(tagpostarr); i++ {
				//创建标签结构体并初始化标签id
				var tag = &models.Tag{Id: tagpostarr[i].Tagid}
				//在标签中查询出对应记录，并且其count字段大于0
				if err = tag.Read(); err == nil && tag.Count > 0 {
					tag.Count--
					tag.Update("count")
				}
			}
		}
		//删除中间表中对应记录
		query.Delete()
	}
	//创建切片，用于存储用户输入的标签的处理结果
	addtags := make([]string, 0)
	//判断用户是否输入了标签
	if tags != "" {
		//通过逗号切割用户输入的标签
		tagarr := strings.Split(tags, ",")
		//遍历切割之后的结果
		for _, v := range tagarr {
			//去除每一个标签的前后空格并判断是否为空
			if tag := strings.TrimSpace(v); tag != "" {
				exists := false
				//遍历最终结果切片
				for _, vv := range addtags {
					if vv == tag {
						exists = true
						break
					}
				}

				if !exists {
					//不存在则添加
					addtags = append(addtags, tag)
				}
			}
		}
	}

	//判断处理之后的切片中是否还有值
	if len(addtags) > 0 {
		//遍历最终结果切片，取出每一个标签名称
		for _, v := range addtags {
			//创建标签并初始化标签名称
			tag := &models.Tag{Name: v}
			//根据标签名称查询，如果不存在则插入
			if err := tag.Read("Name"); err == orm.ErrNoRows {
				tag.Count = 1
				tag.Insert()
			} else {
				tag.Count += 1
				tag.Update("Count")
			}
			//创建标签文章对象并初始化个字段
			tp := &models.TagPost{Tagid: tag.Id, Postid: post.Id, Poststatus: post.Status, Posttime: post.Posttime}
			tp.Insert()
		}
		//用逗号拼接标签并在前后都拼接上逗号
		post.Tags = "," + strings.Join(addtags, ",") + ","
	}
	post.Update("title", "color", "istop", "posttime", "status", "content", "updated", "tags")
	this.Redirect("/admin/article/list", 302)
}

//批量操作
func (this *ArticleController) Batch() {
	//获取用户所选择的文章的id
	ids := this.GetStrings("ids[]")
	//获取用户所选择的操作
	op := this.GetString("op")
	//创建切片，用于存储转换后的文章id
	idarr := make([]int, 0)
	//遍历ids
	for _, v := range ids {
		//取出每一个id并转换为整形
		if id, _ := strconv.Atoi(v); id > 0 {
			idarr = append(idarr, id)
		}
	}
	//获得文章表的句柄
	query := orm.NewOrm().QueryTable(new(models.Post))
	switch op {
	//移至已发布(0)
	case "topub":
		query.Filter("id__in", idarr).Update(orm.Params{"status": 0})
	//移至草稿箱(1)
	case "todrafts":
		query.Filter("id__in", idarr).Update(orm.Params{"status": 1})
	//移至回收站(2)
	case "totrash":
		query.Filter("id__in", idarr).Update(orm.Params{"status": 2})
	//删除
	case "delete":
		for _, id := range idarr {
			//创建文章结构体，并初始化id
			obj := models.Post{Id: id}
			if obj.Read() == nil {
				obj.Delete()
			}
		}
	}

	this.Redirect(this.Ctx.Request.Referer(), 302)
}
