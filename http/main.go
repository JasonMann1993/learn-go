package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello jason！")
}

func main() {

	//resp, err := http.Get("https://www.liwenzhou.com/")
	//if err != nil {
	//	fmt.Printf("get failed, err:%v\n", err)
	//}
	//defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Printf("read from resp.Body failed, err:%v\n", err)
	//	return
	//}
	//fmt.Print(string(body))

	// --------- param -----------
	//apiUrl := "http://localhost"
	//// URL param
	//data := url.Values{}
	//data.Set("name", "小王子")
	//data.Set("age", "18")
	//u, err := url.ParseRequestURI(apiUrl)
	//if err != nil {
	//	fmt.Printf("parse url requestUrl failed, err:%v\n", err)
	//}
	//u.RawQuery = data.Encode() // URL encode
	//resp, err := http.Get(u.String())
	//if err != nil {
	//	fmt.Printf("post failed, err:%v\n", err)
	//	return
	//}
	//defer resp.Body.Close()
	//b, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Printf("get resp failed, err:%v\n", err)
	//	return
	//}
	//fmt.Println(string(b))


	// -----------POST----------
	//url := "http://localhost"
	//// 表单数据
	//contentType := "application/x-www-form-urlencode"
	//data := "name=小王子&age=18"
	//// json
	////contentType := "application/json"
	////data := `{"name":"小王子", "age":18}`
	//res, err := http.Post(url,contentType,strings.NewReader(data))
	//if err != nil {
	//	fmt.Println("err:",err)
	//	return
	//}
	//defer res.Body.Close()
	//b, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println("err:",err)
	//	return
	//}
	//fmt.Println(string(b))

	// 自定义 client
	//client := &http.Client{
	//}
	//req, _ := http.NewRequest("GET", "http://localhost", nil)
	//req.Header.Add("If-None-Match", `W/"wyzzy"`)
	//resp, _ := client.Do(req)
	//defer resp.Body.Close()
	//b, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(b))


	// server
	//http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	//})
	//log.Fatal(http.ListenAndServe(":8080", nil))

	//http.HandleFunc("/",sayHello)
	//err := http.ListenAndServe(":9090", nil)
	//if err != nil {
	//	fmt.Printf("http server failed, err:%v\n", err)
	//	return
	//}

	// 自定义 server
	s := &http.Server{
		Addr: ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}