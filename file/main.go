package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func CopyFile(dstName, srcName string)(written int64, err error) {
	// 以读方式打开源文件
	src, err := os.Open(srcName)
	if err != nil {
		fmt.Println("open src file error:", err)
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open dst file error:", err)
		return
	}
	defer dst.Close()

	return io.Copy(dst,src)
}

// cat命令实现
func cat(r *bufio.Reader) {
	for {
		buf, err := r.ReadBytes('\n') //注意是字符
		if err == io.EOF {
			// 退出之前将已读到的内容输出
			fmt.Fprintf(os.Stdout, "%s", buf)
			break
		}
		fmt.Fprintf(os.Stdout, "%s", buf)
	}
}

func main() {
	_, err := CopyFile("2.txt", "1.txt")
	if err != nil {
		fmt.Println("copy file failed, err:", err)
		return
	}
	fmt.Println("done")

	flag.Parse()
	if flag.NArg() == 0 {
		// 如果没有参数默认从标准输入读取内容
		cat(bufio.NewReader(os.Stdin))
	}

	// 依次读取每个指定文件的内容并打印到终端
	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if err != nil {
			fmt.Fprintf(os.Stdout, "reading from %s failed, err:%v\n", flag.Arg(i), err)
			continue
		}
		cat(bufio.NewReader(f))
	}
}