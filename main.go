package main

import (
	"Bluebell/model"
	"Bluebell/routers"
)

func main() {
	//应用数据库
	model.InitDb()
	routers.InitRouter()
}
