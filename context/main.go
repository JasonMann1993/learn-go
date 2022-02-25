package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 在 Go http包的Server中，每一个请求在都有一个对应的 goroutine 去处理。请求处理函数通常会启动额外的 goroutine 用来访问后端服务，比如数据库和RPC服务。用来处理一个请求的 goroutine 通常需要访问一些与请求特定的数据，比如终端用户的身份认证信息、验证相关的token、请求的截止时间。 当一个请求被取消或超时时，所有用来处理该请求的 goroutine 都应该迅速退出，然后系统才能释放这些 goroutine 占用的资源。

var (
	wg   sync.WaitGroup
	flag bool
)

// 全局变量模式
//func worker() {
//	for {
//		fmt.Println("i")
//		time.Sleep(time.Second)
//		if flag {
//			break
//		}
//	}
//	wg.Done()
//}

// 管道模式
//func worker(exitChan <-chan struct{}) {
//LOOP:
//	for {
//		fmt.Println("work")
//		time.Sleep(time.Second)
//		select {
//		case <-exitChan:
//			break LOOP
//		default:
//		}
//	}
//	wg.Done()
//}

// ctx
func worker(ctx context.Context) {
	go worker2(ctx)
LOOP:
	for {
		fmt.Println("work")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
	}
	wg.Done()
}
func worker2(ctx context.Context) {
LOOP:
	for {
		fmt.Println("worker2")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done(): // 等待上级通知
			break LOOP
		default:
		}
	}
}

// 当一个 goroutine 关闭时如何关闭该 goroutine 启动的其他 goroutine
func main() {
	// 全局变量模式
	// 全局变量方式存在的问题：
	// 1. 使用全局变量在跨包调用时不容易统一
	// 2. 如果worker中再启动goroutine，就不太好控制了。
	//wg.Add(1)
	//go worker()
	//time.Sleep(time.Second * 4)
	//flag = true
	//wg.Wait()
	//fmt.Println("finish")

	// 通道方式
	// 管道方式存在的问题：
	// 1. 使用全局变量在跨包调用时不容易实现规范和统一，需要维护一个共用的channel
	//ch := make(chan struct{})
	//wg.Add(1)
	//go worker(ch)
	//time.Sleep(time.Second * 4)
	//ch <- struct{}{}
	//close(ch)
	//wg.Wait()
	//fmt.Println("success")

	//  ctx
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 4)
	cancel()
	wg.Wait()
	fmt.Println("over")
}
