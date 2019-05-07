package main

import (
	"github.com/micro/go-log"
	"net/http"

	"github.com/micro/go-web"
	"github.com/julienschmidt/httprouter"
	_ "sss/IhomeWeb/models"
	"sss/IhomeWeb/handler"
)

func main() {
	// 创建一个新的web服务
	service := web.NewService(
		web.Name("go.micro.web.IhomeWeb"),
		web.Version("latest"),
		web.Address(":10086"),
	)

	// 服务初始化
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// 使用路由中间件来映射页面
	router := httprouter.New()
	router.NotFound = http.FileServer(http.Dir("html"))

	// 获取地区请求
	router.GET("/api/v1.0/areas",handler.GetArea)

	// 获取session
	router.GET("/api/v1.0/session",handler.GetSession)

	// 获取首页轮播图
	router.GET("/api/v1.0/house/index",handler.GetIndex)

	// 获取验证码图片
	router.GET("/api/v1.0/imagecode/:uuid",handler.GetImageCode)

	// 映射前端页面
	//service.Handle("/", http.FileServer(http.Dir("html")))
	service.Handle("/", router)


	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
