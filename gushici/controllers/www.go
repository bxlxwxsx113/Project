package controllers

import (
	"fmt"
	"gushici/models"
	"math/rand"
)

type WwwController struct {
	BaseController
}

/*
1.明确展示页面需要显示哪些信息
2.分析对应的html
3.去数据库中查询出第二步分出的数据
4.遍历查询结果，将该结果出入模板数据中
*/
func (this *WwwController) Index() {
	//创建切片，用于存储过滤条件
	filter4 := make([]interface{}, 0)
	//将状态1追加到切片中
	filter4 = append(filter4, "status", 1)
	//将class_id=2(诗词古韵)加入过滤条件
	filter4 = append(filter4, "class_id", 2)
	//分页查询，注意：三个点
	result4, _ := models.NewGetList(1, 6, filter4...)
	//创建切片，用于存储一个分类的古诗
	//range $ind, $elem := .data.list4
	list4 := make([]map[string]interface{}, len(result4))
	for k, v := range result4 {
		//创建map，用于存储一首古诗
		row := make(map[string]interface{})
		// id   picurl   media   title   desc
		row["id"] = v.Id
		if v.Picurl == "" {
			//[1,16)
			var r = rand.Intn(16) + 1
			v.Picurl = "/upload/image/rand" + fmt.Sprintf("%d", r) + ".jpeg"
		}
		row["picurl"] = v.Picurl
		row["media"] = v.Media
		row["title"] = v.Title
		//描述部分超出30个汉字，只显示30个汉字
		if v.Desc != "" {
			nameRune := []rune(v.Desc)
			lth := len(nameRune)
			if lth > 30 {
				lth = 30
			}
			row["desc"] = string(nameRune[:lth])
		}
		list4[k] = row
	}

	//创建切片，用于存储过滤条件
	filter2 := make([]interface{}, 0)
	//将状态1追加到切片中
	filter2 = append(filter2, "status", 1)
	//将class_id=2(儿童古诗)加入过滤条件
	filter2 = append(filter2, "class_id", 3)
	//分页查询，注意：三个点
	result2, _ := models.NewGetList(1, 6, filter2...)
	//创建切片，用于存储一个分类的古诗
	list2 := make([]map[string]interface{}, len(result2))
	for k, v := range result2 {
		//创建map，用于存储一首古诗
		row := make(map[string]interface{})
		// id   picurl   media   title   desc
		row["id"] = v.Id
		if v.Picurl == "" {
			//[1,16)
			var r = rand.Intn(16) + 1
			v.Picurl = "/upload/image/rand" + fmt.Sprintf("%d", r) + ".jpeg"
		}
		row["picurl"] = v.Picurl
		row["media"] = v.Media
		row["title"] = v.Title
		//描述部分超出30个汉字，只显示30个汉字
		if v.Desc != "" {
			nameRune := []rune(v.Desc)
			lth := len(nameRune)
			if lth > 30 {
				lth = 30
			}
			row["desc"] = string(nameRune[:lth])
		}
		list2[k] = row
	}

	//创建切片，用于存储过滤条件
	filter := make([]interface{}, 0)
	//将状态1追加到切片中
	filter = append(filter, "status", 1)
	//将class_id=2(开心儿歌)加入过滤条件
	filter = append(filter, "class_id", 5)
	//分页查询，注意：三个点
	result, _ := models.NewGetList(1, 6, filter...)
	//创建切片，用于存储一个分类的古诗
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		//创建map，用于存储一首古诗
		row := make(map[string]interface{})
		// id   picurl   media   title   desc
		row["id"] = v.Id
		if v.Picurl == "" {
			//[1,16)
			var r = rand.Intn(16) + 1
			v.Picurl = "/upload/image/rand" + fmt.Sprintf("%d", r) + ".jpeg"
		}
		row["picurl"] = v.Picurl
		row["media"] = v.Media
		row["title"] = v.Title
		//描述部分超出30个汉字，只显示30个汉字
		if v.Desc != "" {
			nameRune := []rune(v.Desc)
			lth := len(nameRune)
			if lth > 30 {
				lth = 30
			}
			row["desc"] = string(nameRune[:lth])
		}
		list[k] = row
	}

	//创建切片，用于存储过滤条件
	filter3 := make([]interface{}, 0)
	//将状态1追加到切片中
	filter3 = append(filter3, "status", 1)
	//将class_id=2(诗词古韵)加入过滤条件
	filter3 = append(filter3, "class_id", 1)
	//分页查询，注意：三个点
	result3, _ := models.NewGetList(1, 6, filter3...)
	//创建切片，用于存储一个分类的古诗
	list3 := make([]map[string]interface{}, len(result3))
	for k, v := range result3 {
		//创建map，用于存储一首古诗
		row := make(map[string]interface{})
		// id   picurl   media   title   desc
		row["id"] = v.Id
		if v.Picurl == "" {
			//[1,16)
			var r = rand.Intn(16) + 1
			v.Picurl = "/upload/image/rand" + fmt.Sprintf("%d", r) + ".jpeg"
		}
		row["picurl"] = v.Picurl
		row["media"] = v.Media
		row["title"] = v.Title
		//描述部分超出30个汉字，只显示30个汉字
		if v.Desc != "" {
			nameRune := []rune(v.Desc)
			lth := len(nameRune)
			if lth > 30 {
				lth = 30
			}
			row["desc"] = string(nameRune[:lth])
		}
		list3[k] = row
	}

	out := make(map[string]interface{})
	out["list"] = list
	out["list2"] = list2
	out["list3"] = list3
	out["list4"] = list4
	out["class_id"] = 0
	this.Data["data"] = out
	//range $ind, $elem := .data.list4

	this.Layout = "public/www_layout.html"
	this.display()
}

//通过古诗id查询古诗
func (this *WwwController) Show() {
	//获取前台页面传递过来的id
	id, _ := this.GetInt(":id")
	//根据id查询古诗
	News, _ := models.NewsGetById(id)
	//创建map，用于存储页面上需要显示的信息
	row := make(map[string]interface{})
	row["class_id"] = News.ClassId
	// title  picurl  media  content
	//判断是否查询到了古诗
	if News != nil {
		//标题
		row["title"] = News.Title
		if News.Picurl == "" {
			//[1,16)
			var r = rand.Intn(16) + 1
			News.Picurl = "/upload/image/rand" + fmt.Sprintf("%d", r) + ".jpeg"
		}
		//图片路径
		row["picurl"] = News.Picurl
		//音频
		row["media"] = News.Media
		//内容
		row["content"] = News.Content
	}

	this.Data["data"] = row
	this.Layout = "public/www_layout.html"
	this.display()
}

//根据分类id查询古诗
func (this *WwwController) List() {
	//获取分类id
	catId, cerr := this.GetInt(":class_id")
	//创建切片，用于存储过滤条件
	filters := make([]interface{}, 0)
	//将正常状态这个条件放到切片中
	filters = append(filters, "status", 1)
	if cerr == nil {
		//将class_id追加到切片中
		filters = append(filters, "class_id", catId)
	}
	this.pageSize = 16
	//分页查询
	result, count := models.NewGetList(1, this.pageSize, filters...)
	//分类切片
	list := make([]map[string]interface{}, len(result))
	//id   picurl   media  title  desc
	for k, v := range result {
		//创建map，用于存储每一首古诗
		row := make(map[string]interface{})
		row["id"] = v.Id
		if v.Picurl == "" {
			//[1,16)
			var r = rand.Intn(16) + 1
			v.Picurl = "/upload/image/rand" + fmt.Sprintf("%d", r) + ".jpeg"
		}
		//图片路径
		row["picurl"] = v.Picurl
		//音频
		row["media"] = v.Media
		//标题
		row["title"] = v.Title
		//描述部分超出30个汉字，只显示30个汉字
		if v.Desc != "" {
			nameRune := []rune(v.Desc)
			lth := len(nameRune)
			if lth > 30 {
				lth = 30
			}
			row["desc"] = string(nameRune[:lth])
		}
		list[k] = row
	}
	classArr := make(map[int]string)
	classArr[1] = "国学经典"
	classArr[2] = "诗词古韵"
	classArr[3] = "儿童古诗"
	classArr[5] = "开心儿歌"
	out := make(map[string]interface{})
	//分类id
	out["class_id"] = catId
	//分类名称
	out["class_name"] = classArr[catId]
	//满足过滤条件的记录的数量
	out["count"] = count
	//当前页码
	out["page"] = 1
	out["list"] = list
	this.Data["data"] = out

	this.Layout = "public/www_layout.html"
	this.display()
}
