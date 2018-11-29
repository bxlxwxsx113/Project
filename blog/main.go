package main

import (
	_ "blog/models"
	_ "blog/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
