package main

import (
	"fmt"
	"sync"
	"time"
)

func dosomething(o *sync.Once) {
	fmt.Println("start:")
	o.Do(func(){
		fmt.Println("Do something...")
	})
	fmt.Println("finish.")
}

func main() {
	// sync once
	o := &sync.Once{}
	go dosomething(o)
	go dosomething(o)
	time.Sleep(time.Second * 1)
}