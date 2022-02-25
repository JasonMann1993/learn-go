package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile("./1.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("openFile failed,err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Ldate )
	log.Println("普通日志")
	log.SetPrefix("前缀")
	v := "很普通的"
	log.Printf("很普通%s\n", v)
	//log.Fatalln("fatal 日志")
	//log.Panicln("panic 日志")
	logger := log.New(os.Stdout, "<New>", log.Lshortfile|log.Ldate|log.Ltime)
	logger.Println("这是自定义的logger记录的日志。")
}
