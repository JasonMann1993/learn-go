package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接 非v8
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

// 初始化连接 v8
//func initClient() (err error){
//	rdb = redis.NewClient(&redis.Options{
//		Addr: "localhost:6379",
//		Password: "", // no password set
//		DB: 0, // use default DB
//		PoolSize: 100, // 连接池大小
//	})
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	_, err = rdb.Ping(ctx).Result()
//	return err
//}
//func V8Example() {
//	ctx := context.Background()
//	if err := initClient(); err != nil {
//		return
//	}
//
//	err := rdb.Set(ctx, "key", "value", time.Second * 60).Err()
//	if err != nil {
//		panic(err)
//	}
//
//	val, err := rdb.Get(ctx, "key").Result()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("key:", val)
//
//	val2, err := rdb.Get(ctx, "key2").Result()
//	if err == redis.Nil {
//		fmt.Println("key2 does not exist")
//	}else if err != nil {
//		panic(err)
//	}else {
//		fmt.Println("key2", val2)
//	}
//	// Output: key value
//	// key2 does not exist
//}

// 哨兵模式
//func initClient()(err error){
//	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
//		MasterName:    "master",
//		SentinelAddrs: []string{"x.x.x.x:26379", "xx.xx.xx.xx:26379", "xxx.xxx.xxx.xxx:26379"},
//	})
//	_, err = rdb.Ping().Result()
//	if err != nil {
//		return err
//	}
//	return nil
//}

func main() {
	err := initClient()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success")
	}
}
