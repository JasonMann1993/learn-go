package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f := bufio.NewReader(os.Stdin)
	fmt.Println("请输入内容：")
	a,_ := f.ReadString('\n')
	fmt.Println("输入的是",a)
}