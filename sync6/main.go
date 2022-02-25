package main

import (
	"fmt"
	"sync"
)

func addNum( a, b int, deferFunc func()) {
	defer func(){
		deferFunc()
	}()
	c := a + b
	fmt.Printf("%d + %d = %d\n", a, b, c)
}

func main() {
	// sync.WaitGroup 分批
	total := 10
	step := 2
	fmt.Println("启动子协程")
	var wg sync.WaitGroup
	for i := 0; i < total; i = i + step {
		wg.Add(step)
		for j := 0; j < step; j++ {
			go addNum(i + j, 1, wg.Done)
		}
		wg.Wait()
	}
	fmt.Println("所有子协程执行完毕.")
}