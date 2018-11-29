package models

import (
	"bytes"
	"fmt"
)

type Pager struct {
	Page     int    //页码
	Pagesize int    //每页显示几条记录
	Totalnum int    //总的记录数
	Urlpath  string //每页对应的url
}

//创建pager对象
func NewPager(page, pagesize, totalnum int, urlpath string) *Pager {
	pager := new(Pager)
	pager.Page = page
	pager.Pagesize = pagesize
	pager.Totalnum = totalnum
	pager.Urlpath = urlpath
	return pager
}

//设置page
func (this *Pager) SetPage(page int) {
	this.Page = page
}

//设置pagesize
func (this *Pager) SetPagesize(pagesize int) {
	this.Pagesize = pagesize
}

//设置总数量
func (this *Pager) SetTotalnum(totalnum int) {
	this.Totalnum = totalnum
}

//设置urlpath
func (this *Pager) SetUrlpath(urlpath string) {
	this.Urlpath = urlpath
}

func (this *Pager) url(page int) string {
	//   "/index%d.html"
	//fmt.Sprintf("/index%d.html", 3)
	return fmt.Sprintf(this.Urlpath, page)
}

func (this *Pager) ToString() string {
	//文章总数量小于等于每页显示的文章的数量
	if this.Totalnum <= this.Pagesize {
		return ""
	}
	//偏移量
	offset := 5
	//显示10个页码
	linknum := 10
	//总的页码
	var totalpage int
	//计算总的页码
	if this.Totalnum%this.Pagesize != 0 {
		totalpage = this.Totalnum/this.Pagesize + 1
	} else {
		totalpage = this.Totalnum / this.Pagesize
	}

	var from int //开始页码
	var to int   //显示到哪一页
	//总页数小于10，直接从第一页显示到最后一页
	if totalpage < linknum {
		from = 1
		to = totalpage
	} else {
		from = this.Page - offset
		to = from + linknum
		if from < 1 {
			from = 1
			to = from + linknum - 1 //  1 + 10 - 1
		} else if to > totalpage { //结束页大于总页数
			to = totalpage
			from = to - linknum + 1 //  20 - 10 = 10 + 1 = 11
		}
	}
	var buf bytes.Buffer
	buf.WriteString("<div class='page'>")
	//上一页
	if this.Page > 1 {
		/*
			4  3
			func (this *Pager) url(page int) string {
				//   "/index%d.html"
				//fmt.Sprintf("/index%d.html", 3)  /index3.html
				return fmt.Sprintf(this.Urlpath, page)
			}

		*/
		buf.WriteString(fmt.Sprintf("<a href='%s'>&laquo;</a>", this.url(this.Page-1))) //   <<
	}
	for i := from; i <= to; i++ {
		if i == this.Page {
			buf.WriteString(fmt.Sprintf("<b>%d</b>", i))
		} else {
			buf.WriteString(fmt.Sprintf("<a href='%s'>%d</a>", this.url(i), i))
		}
	}

	//下一页
	if this.Page < totalpage {
		buf.WriteString(fmt.Sprintf("<a href='%s'>&raquo;</a>", this.url(this.Page+1)))
	}

	buf.WriteString("</div>")
	return buf.String()
}
