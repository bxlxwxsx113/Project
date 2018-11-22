package models

import "github.com/astaxie/beego/orm"

//古诗
type InfoList struct {
	Id         int
	ClassId    int    //分类ID
	Title      string //标题
	Author     string //作者
	Keywords   string //关键字
	Desc       string //描述
	Content    string //内容
	Picurl     string //图片路径
	Media      string //诗词音频
	Posttime   int64  //提交时间
	Updatetime int64  //修改时间
	Status     int    //状态  1:正常状态   0：被删除的状态
	Orderid    int    //排序
}

func (infolist *InfoList) TableName() string {
	return TableName("info_list")
}

//page:当前页码
//pagesize:每页显示的记录的数量
//filters:过滤条件
//返回值一:查询结果
//返回值二：符合过滤条件总的记录的数量
func NewGetList(page, pageSize int, filters ...interface{}) ([]*InfoList, int64) {
	//获得句柄
	query := orm.NewOrm().QueryTable(TableName("info_list"))
	//判断过滤条件切片中是否有值
	if len(filters) > 0 {
		//获取过滤条件长度
		length := len(filters)
		/*
			filter4 = append(filter4, "status", 1)
			filter4 = append(filter4, "class_id", 2)
		*/
		//query.Filter("status", 1).Filter("class_id", 2)
		//拼接过滤条件
		//[status, 1, class_id, 2]
		for k := 0; k < length; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	//获取满足过滤条件的记录的数量
	total, _ := query.Count()
	//偏移量
	offset := (page - 1) * pageSize
	//存储查询结果
	list := make([]*InfoList, 0)
	//分页查询
	query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	return list, total
}

//根据id查询古诗
func NewsGetById(id int) (*InfoList, error) {
	//创建古诗对象
	r := new(InfoList)
	//根据指定的id查询出古诗
	err := orm.NewOrm().QueryTable(TableName("info_list")).Filter("id", id).One(r)
	//处理错误
	if err != nil {
		return nil, err
	}
	return r, nil
}
