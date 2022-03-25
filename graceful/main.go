package main

import (
	"context"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func graceShutdown() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5* time.Second)
		c.String(http.StatusOK, "gin server 5s")
	})

	srv := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待终端信号来优雅的 关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1)// 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的 Ctrl + C 就是触发系统 SIGINT 信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify 把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给 quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Println("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("Server exiting")
}


// 优雅重启
// 我们可以使用 fvbock/endless 来替换默认的 ListenAndServe启动服务来实现
func graceRestart() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5*time.Second)
		c.String(http.StatusOK, "hello jason!")
	})
	//默认endless服务器会监听下列信号：
	// syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2，syscall.SIGINT，syscall.SIGTERM 和 syscall.SIGTSTP
	// 接收到 SIGHUP 信号将触发 `fork/restart` 实现优雅重启(kill -1 pid会发送SIGHUP信号)
	// 接收到 syscall.SIGINT或syscall.SIGTERM 信号将触发优雅关机
	// 接收到 SIGUSR2 信号将触发HammerTime
	// SIGUSR1 和 SIGTSTP 被用来触发一些用户自定义的hook函数
	if err := endless.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("listen: %s\n", err)
	}
	log.Println("server exiting")
}

func main() {
	//graceShutdown()
	graceRestart()
}