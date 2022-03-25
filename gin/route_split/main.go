package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/routers"
	"net/http"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello q1mi!",
	})
}

// 下面最基础的gin路由注册方式，适用于路由条目比较少的简单项目或者项目demo。
func base() {
	r := gin.Default()
	r.GET("/hello", helloHandler)
	if err := r.Run(); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}


// 拆分文件
func split() {
	//当项目的规模增大后就不太适合继续在项目的main.go文件中去实现路由注册相关逻辑了，我们会倾向于把路由部分的代码都拆分出来，形成一个单独的文件或包：
	//我们在routers.go文件中定义并注册路由信息
	// 调用 routers.go 中定义的setupRouter
	r := routers.SetupRouter()
	if err := r.Run(); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}


// 拆分多个文件
func multiSplit() {
	r := gin.Default()
	routers.LoadShop(r)
	routers.LoadBlog(r)
	if err := r.Run(); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}


func main() {
	// 基础
	//base()



}

//路由拆分到不同的APP
//有时候项目规模实在太大，那么我们就更倾向于把业务拆分的更详细一些，例如把不同的业务代码拆分成不同的APP。
//
//因此我们在项目目录下单独定义一个app目录，用来存放我们不同业务线的代码文件，这样就很容易进行横向扩展。大致目录结构如下：
//
//gin_demo
//├── app
//│   ├── blog
//│   │   ├── handler.go
//│   │   └── router.go
//│   └── shop
//│       ├── handler.go
//│       └── router.go
//├── go.mod
//├── go.sum
//├── main.go
//└── routers
//└── routers.go
//其中app/blog/router.go用来定义blog相关的路由信息，具体内容如下：
//
//func Routers(e *gin.Engine) {
//	e.GET("/post", postHandler)
//	e.GET("/comment", commentHandler)
//}
//app/shop/router.go用来定义shop相关路由信息，具体内容如下：
//
//func Routers(e *gin.Engine) {
//	e.GET("/goods", goodsHandler)
//	e.GET("/checkout", checkoutHandler)
//}
//routers/routers.go中根据需要定义Include函数用来注册子app中定义的路由，Init函数用来进行路由的初始化操作：
//
//type Option func(*gin.Engine)
//
//var options = []Option{}
//
//// 注册app的路由配置
//func Include(opts ...Option) {
//	options = append(options, opts...)
//}
//
//// 初始化
//func Init() *gin.Engine {
//	r := gin.Default()
//	for _, opt := range options {
//		opt(r)
//	}
//	return r
//}
//main.go中按如下方式先注册子app中的路由，然后再进行路由的初始化：
//
//func main() {
//	// 加载多个APP的路由配置
//	routers.Include(shop.Routers, blog.Routers)
//	// 初始化路由
//	r := routers.Init()
//	if err := r.Run(); err != nil {
//		fmt.Println("startup service failed, err:%v\n", err)
//	}
//}