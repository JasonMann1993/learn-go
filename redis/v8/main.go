package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

// v8
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 100,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}

// 哨兵模式
//func initClient() (err error) {
//	rdb = redis.NewFailoverClient(&redis.FailoverOptions{
//		MasterName:    "mymaster",
//		SentinelAddrs: []string{":26379", ":26380", ":26381"},
//		Password: "",
//		DB: 0,
//	})
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	_, err = rdb.Ping(ctx).Result()
//	if err != nil {
//		return err
//	}
//	return nil
//}

// 集群
//func initClient()(err error){
//	rdb := redis.NewClusterClient(&redis.ClusterOptions{
//		Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
//	})
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//	_, err = rdb.Ping(ctx).Result()
//	if err != nil {
//		return err
//	}
//	return nil
//}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func V8Example() {
	ctx := context.Background()
	if err := initClient(); err != nil {
		return
	}

	if rdb == nil {
		fmt.Println("redis client connect is nil")
		return
	}

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

}

// zset 示例
func redisExample2() {
	zsetKey := "language_rank"
	languages := []*redis.Z{
		&redis.Z{Score: 90.0, Member: "Golang"},
		&redis.Z{Score: 98.0, Member: "Java"},
		&redis.Z{Score: 95.0, Member: "Python"},
		&redis.Z{Score: 97.0, Member: "JavaScript"},
		&redis.Z{Score: 99.0, Member: "C/C++"},
	}
	ctx := context.Background()
	// ZADD
	num, err := rdb.ZAdd(ctx, zsetKey, languages...).Result()

	if err != nil {
		fmt.Printf("zadd failed, err:%v\n", err)
		return
	}

	fmt.Printf("zadd %d succ.\n", num)

	// 把Golang的分数加10
	newScore, err := rdb.ZIncrBy(ctx, zsetKey, 10.0, "Golang").Result()

	if err != nil {
		fmt.Printf("zincrby failed, err:%v\n", err)
		return
	}
	fmt.Printf("Golang's score is %f now.\n", newScore)

	ret, err := rdb.ZRevRangeWithScores(ctx, zsetKey, 0, 2).Result()
	if err != nil {
		fmt.Printf("zrevrange failed, err:%v\n", err)
		return
	}
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}

	// 取95~100分的
	op := &redis.ZRangeBy{
		Min: "95",
		Max: "100",
	}
	ret, err = rdb.ZRangeByScoreWithScores(ctx, zsetKey, op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return
	}
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}
}

// watch 加 事务
func transactionDemo() {
	var (
		maxRetries = 1000
		routineCount = 10
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Increment 使用GET和SET命令以事务方式递增Key的值
	increment := func(key string) error {
		// 事务函数
		txf := func(tx *redis.Tx) error {
			// 获取可以的当前值或清零
			n, err := tx.Get(ctx, key).Int()
			if err != nil && err != redis.Nil {
				return err
			}

			// 实际的操作代码（乐观锁定中的本地操作）
			n++

			// 操作仅在 Watch 的 Key 没发生变化的情况下提交
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, key, n, 0)
				return nil
			})
			return err
		}

		// 最多重试 maxRetries 次
		for i := 0; i < maxRetries; i++ {
			err := rdb.Watch(ctx, txf, key)
			if err == nil {
				// 成功
				return nil
			}
			if err == redis.TxFailedErr {
				// 乐观锁丢失 重试
				continue
			}
			// 返回其他的错误
			return err
		}
		return errors.New("increment reached maximum number of retries")
	}
	// 模拟 routineCount 个并发同时去修改 counter3 的值
	var wg sync.WaitGroup
	wg.Add(routineCount)
	for i := 0; i < routineCount; i++ {
		go func() {
			defer wg.Done()
			if err := increment("counter3"); err != nil {
				fmt.Println("increment error:", err)
			}
		}()
	}
	wg.Wait()

	n, err := rdb.Get(context.TODO(), "counter3").Int()
	fmt.Println("ended with", n, err)
}

func main() {
	//V8Example()
	if err := initClient(); err != nil {
		return
	}
	//redisExample2()
	// 根据前缀获取key
	//vals, err := rdb.Keys(ctx,"prefix*").Result()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(vals)

	// 自定义执行命令
	//res, err := rdb.Do(ctx, "set", "key", "value").Result()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(res)

	// key比较多时，按通配符删除key
	//iter := rdb.Scan(ctx, 0, "*", 0).Iterator()
	//for iter.Next(ctx) {
	//	err := rdb.Del(ctx, iter.Val()).Err()
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(iter.Val())
	//}
	//if err := iter.Err(); err != nil {
	//	panic(err)
	//}

	// 缓冲命令，一次性发送
	//pipe := rdb.Pipeline()
	//incr := pipe.Incr(ctx,"pipeline_counter")
	//pipe.Expire(ctx,"pipeline_counter", time.Hour)
	//_, err := pipe.Exec(ctx)
	//fmt.Println(incr.Val(),err)

	// 另一种写法
	//var incr *redis.IntCmd
	//_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
	//	incr = pipe.Incr(ctx,"pipelined_counter")
	//	pipe.Expire(ctx,"pipelined_counter", time.Hour)
	//	return nil
	//})
	//fmt.Println(incr.Val(), err)

	// 事务
	//pipe := rdb.TxPipeline()
	//
	//incr := pipe.Incr(ctx,"tx_pipeline_counter")
	//pipe.Expire(ctx, "tx_pipeline_counter", time.Hour)
	//
	//_, err := pipe.Exec(ctx)
	//fmt.Println(incr.Val(), err)

	// 另一种写法
	//var incr *redis.IntCmd
	//_, err := rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
	//	incr = pipe.Incr(ctx, "tx_pipelined_counter")
	//	pipe.Expire(ctx, "tx_pipelined_counter", time.Hour)
	//	return nil
	//})
	//fmt.Println(incr.Val(), err)

	// watch
	//key := "watch_count"
	//err := rdb.Watch(ctx, func(tx *redis.Tx) error {
	//	n, err := tx.Get(ctx, key).Int()
	//	if err != nil && err != redis.Nil {
	//		return err
	//	}
	//	time.Sleep(time.Second * 4)// 此时在另外一个终端窗口中使用set命令修改watch_count的值
	//	_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
	//		pipe.Set(ctx, key, n+1, 0)
	//		return nil
	//	})
	//	return err
	//}, key)
	//if err != nil {
	//	fmt.Printf("key 发生变化，更新失败。err:%v\n", err)
	//} else {
	//	fmt.Println("更新成功")
	//}

	transactionDemo()
}
