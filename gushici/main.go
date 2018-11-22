package main

import (
	"github.com/astaxie/beego"
	_ "gushici/models"
	_ "gushici/routers"
)

func main() {
	beego.Run()
}
