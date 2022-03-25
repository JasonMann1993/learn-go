package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"golang.org/x/sync/errgroup"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// restful 示例
func restFul() {
	r := gin.Default()
	r.GET("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "GET",
		})
	})

	r.POST("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "POST",
		})
	})

	r.PUT("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "PUT",
		})
	})

	r.DELETE("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "DELETE",
		})
	})
}

// HTML 渲染
func htmlRender() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	// r.loadHTMLFiles("templates/posts/index.html", "templates/users/index.html")
	r.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.html", gin.H{
			"title": "posts/index",
		})
	})

	r.GET("users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.html", gin.H{
			"title": "users/index",
		})
	})

	r.GET("users/edit", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/edit.html", gin.H{
			"title": "users/edit",
		})
	})

	r.Run(":8080")
}

// 自定义模版
func selfTpl() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})

	router.LoadHTMLFiles("./index.tmpl")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", "<a href='https://liwenzhou.com'>李文周的博客</a>")
	})

	router.Run(":8080")
}

// 静态文件处理
func static() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/**/*")
	// ...
	r.Run(":8080")
}

// 模板继承 假设有以下文件
//templates
//├── includes
//│   ├── home.tmpl
//│   └── index.tmpl
//├── layouts
//│   └── base.tmpl
//└── scripts.tmpl

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.tmpl")
	if err != nil {
		panic(err.Error())
	}
	// 为layouts/和includes/目录生成 templates map
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include),files...)
	}
	return r
}
func indexFunc(c *gin.Context){
	c.HTML(http.StatusOK, "index.tmpl", nil)
}

func homeFunc(c *gin.Context){
	c.HTML(http.StatusOK, "home.tmpl", nil)
}

// 获取当前执行程序路径
func getCurrentPath() string {
	if ex, err := os.Executable(); err == nil {
		return filepath.Dir(ex)
	}
	return "./"
}

// xml、YMAL、protobuf渲染
func xml() {
	r := gin.Default()
	// gin.H 是map[string]interface{}的缩写
	r.GET("/someXML", func(c *gin.Context) {
		// 方式一：自己拼接JSON
		c.XML(http.StatusOK, gin.H{"message": "Hello world!"})
	})
	r.GET("/moreXML", func(c *gin.Context) {
		// 方法二：使用结构体
		type MessageRecord struct {
			Name    string
			Message string
			Age     int
		}
		var msg MessageRecord
		msg.Name = "小王子"
		msg.Message = "Hello world!"
		msg.Age = 18
		c.XML(http.StatusOK, msg)
	})

	r.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "ok", "status": http.StatusOK})
	})

	r.GET("/someProtoBuf", func(c *gin.Context) {
		reps := []int64{int64(1), int64(2)}
		label := "test"
		// protobuf 的具体定义写在 testdata/protoexample 文件中。
		data := &protoexample.Test{
			Label: &label,
			Reps:  reps,
		}
		// 请注意，数据在响应中变为二进制数据
		// 将输出被 protoexample.Test protobuf 序列化了的数据
		c.ProtoBuf(http.StatusOK, data)
	})
	r.Run(":8080")
}

// 获取querystring参数
func queryString() {
	// default 返回一个默认的路由引擎
	r := gin.Default()
	r.GET("/user/search", func(c *gin.Context) {
		username := c.DefaultQuery("username", "小王子")
		//username := c.Query("username")
		address := c.Query("address")
		// 输出json结果给调用方
		c.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			"address":  address,
		})
	})
	r.Run(":8080")
}

// 获取form参数
func formPara() {
	r := gin.Default()
	r.POST("/user/search", func(c *gin.Context) {
		// DefaultPostForm取不到值时会返回指定的默认值
		//username := c.DefaultPostForm("username", "小王子")
		username := c.PostForm("username")
		address := c.PostForm("address")
		//输出json结果给调用方
		c.JSON(http.StatusOK, gin.H{
			"message":  "ok",
			"username": username,
			"address":  address,
		})
	})
	r.Run(":8080")
}

// 获取json参数
func jsonPara() {
	r := gin.Default()
	r.POST("/json", func(c *gin.Context) {
		b, _ := c.GetRawData() // 从c.Request.Body读取请求数据
		// 定义map或结构体
		var m map[string]interface{}
		// 反序列化
		_ = json.Unmarshal(b, &m)

		c.JSON(http.StatusOK, m)
	})
}

// Login 参数绑定
type Login struct {
	User string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func paraBind() {
	router := gin.Default()
	// 绑定JSON的示例 （{"user": "json", "password": "123"}）
	router.POST("/loginJSON", func(c *gin.Context) {
		var login Login
		if err := c.ShouldBind(&login); err == nil {
			fmt.Printf("login info:%#v\n",login)
			c.JSON(http.StatusOK, gin.H{
				"user": login.User,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		}
	})
	// 绑定form表单示例 （user=jason&password=123)
	router.POST("/loginForm", func(c *gin.Context) {
		var login Login
		// ShouldBind()会根据请求的Content-Type自行选择绑定器
		if err := c.ShouldBind(&login); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"user":     login.User,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	// 绑定QueryString示例 (/loginQuery?user=jason&password=123)
	router.GET("/loginQuery", func(c *gin.Context) {
		var login Login
		// ShouldBind()会根据请求的Content-Type自行选择绑定器
		if err := c.ShouldBind(&login); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"user":     login.User,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	router.Run(":8080")
}

// 获取path参数
func pathPara() {
	r := gin.Default()
	r.GET("/user/search/:username/:address", func(c *gin.Context) {
		username := c.Param("username")
		address := c.Param("address")
		// 输出json结果给调用方
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"username": username,
			"address": address,
		})
	})

	r.Run(":8080")
}

// 文件上传
func fileUp () {
	router := gin.Default()
	router.LoadHTMLFiles("file_up.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK,"file_up.html",gin.H{})
	})
	// 处理 multipart forms提交文件时默认的内存限制时 32 MiB
	// 可以通过下面的方式修改
	// router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// 单个文件
		file, err := c.FormFile("f1")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		log.Println(file.Filename)
		dst := fmt.Sprintf("/Users/mj/GolandProjects/learn-go/gin/default/%s", file.Filename)
		// 上传文件到指定目录
		c.SaveUploadedFile(file, dst)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	})
	router.Run(":8080")
}

// 多文件上传
func multiFileUp() {
	router := gin.Default()
	// 处理multipart forms提交文件时默认的内存限制是32 MiB
	// 可以通过下面的方式修改
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["file"]

		for index, file := range files {
			log.Println(file.Filename)
			dst := fmt.Sprintf("/Users/mj/GolandProjects/learn-go/gin/default/%s_%d", file.Filename, index)
			// 上传文件到指定的目录
			c.SaveUploadedFile(file, dst)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%d files uploaded!", len(files)),
		})
	})
	router.Run(":8080")
}

// 重定向
func redirect() {
	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		// 网址重定向
		//c.Redirect(http.StatusMovedPermanently, "https://www.sogo.com")
		// 路由重定向
		c.Request.URL.Path = "/test2"
		r.HandleContext(c)
	})
	r.GET("/test2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})
	// 匹配所有方法的路由
	r.Any("/test3",func(c *gin.Context) {

	})
	// 配置 404 页面
	r.LoadHTMLFiles("404.html")
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})
	// 路由组
	userGroup := r.Group("/user")
	{
		userGroup.GET("/index", func(c *gin.Context) {

		})
		userGroup.GET("/login", func(c *gin.Context) {

		})
		userGroup.POST("/login", func(c *gin.Context) {

		})

	}
	r.Run(":8080")
}


// 中间件
// Gin中的中间件必须是一个gin.HandlerFunc类型。例如我们像下面的代码一样定义一个统计请求耗时的中间件。
// StatCost 是一个统计耗时请求耗时的中间件
func StatCost() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Set("name", "jason")// 可以通过c.Set在请求上下文中设置值，后续的处理函数能够取到该值
		// 调用该请求的剩余处理程序
		c.Next()
		// 不调用该请求的剩余处理程序
		// c.Abort()
		// 计算耗时
		cost := time.Since(start)
		log.Println(cost)

		//gin中间件中使用goroutine
		//当在中间件或handler中启动新的goroutine时，不能使用原始的上下文（c *gin.Context），必须使用其只读副本（c.Copy()）。
	}
}
// 添加中间件
func middleware() {
	//中间件注意事项
	//gin默认中间件
	//gin.Default()默认使用了Logger和Recovery中间件，其中：
	//
	//Logger中间件将日志写入gin.DefaultWriter，即使配置了GIN_MODE=release。
	//Recovery中间件会recover任何panic。如果有panic的话，会写入500响应码。
	//如果不想使用上面两个默认的中间件，可以使用gin.New()新建一个没有任何默认中间件的路由。
	r := gin.New()
	// 注册一个全局中间件
	r.Use(StatCost())

	r.GET("/test", func(c *gin.Context) {
		name := c.MustGet("name").(string) // 从上下文取值
		log.Println(name)
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world!",
		})
	})

	// 为单个路由注册
	r.GET("/test2", StatCost(), func(c *gin.Context){
		name := c.MustGet("name").(string) // 从上下文取值
		log.Println(name)
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello world!",
		})
	})

	// 路由组注册
	// 写法1
	shopGroup := r.Group("/shop", StatCost())
	{
		shopGroup.GET("/index", func(c *gin.Context) {

		})
	}
	// 写法2
	shopGroup2 := r.Group("/shop2")
	shopGroup2.Use(StatCost())
	{
		shopGroup2.GET("/index", func(c *gin.Context) {

		})
	}
	r.Run()
}



// 在多个端口运行多个服务
var (
	g errgroup.Group
)
func router01() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 01",
		})
	})
	return e
}
func router02() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 02",
			},
		)
	})

	return e
}
func multiServer() {
	server01 := &http.Server{
		Addr:         ":8080",
		Handler:      router01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server02 := &http.Server{
		Addr:         ":8081",
		Handler:      router02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// 借助errgroup.Group或者自行开启两个goroutine分别启动两个服务
	g.Go(func() error {
		return server01.ListenAndServe()
	})
	g.Go(func() error {
		return server02.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	//// 创建一个默认的路由引擎
	//r := gin.Default()
	//// GET: 请求方式； /hello： 请求的路径
	//// 当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
	//r.GET("/hello", func(c *gin.Context) {
	//	// c.JSON:返回JSON格式数据， gin.H 是map[string]interface{}的缩写
	//	c.JSON(200, gin.H{
	//		"message": "Hello world!",
	//	})
	//})
	//// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	//r.Run()

	//htmlRender()

	//selfTpl()

	// 模板继承
	//r := gin.Default()
	//r.HTMLRender = loadTemplates("./templates")
	//r.GET("/index", indexFunc)
	//r.GET("/home", homeFunc)
	//r.Run()


	// xml、YMAL、protobuf渲染
	//xml()


	// 获取query参数
	//queryString()

	// 获取 form 参数
	//formPara()

	// 获取path参数
	//pathPara()

	// 参数绑定
	//paraBind()

	// 文件上传
	//fileUp()

	// 多文件上传
	//multiFileUp()

	// 重定向
	//redirect()

	// 中间件
	//middleware()

	multiServer()
}
