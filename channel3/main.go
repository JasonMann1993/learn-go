package main

import (
	"fmt"
	"time"
)

func test(ch chan<- int) {
	for i := 0; i < 4; i++ {
		ch <- i
	}
	close(ch)
}

func main() {
	// 单向通道
	start := time.Now()
	ch := make(chan int, 20)
	go test(ch)
	for i := range ch {
		fmt.Println("接收的数据：", i)
	}
	end := time.Now()
	consume := end.Sub(start).Seconds()
	fmt.Println("程序执行耗时(s):", consume)
}
