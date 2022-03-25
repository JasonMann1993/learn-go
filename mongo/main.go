package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Student struct {
	Name string
	Age int
}

// 连接池模式
func ConnectToDB(uri, name string, timeout time.Duration, num uint64) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	o := options.Client().ApplyURI(uri)
	o.SetMaxPoolSize(num)

	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil, err
	}

	return client.Database(name), nil
}

func main() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}



	// 指定获取要操作的数据集
	collection := client.Database("jason").Collection("student")
	// 查询
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
		// do something with result....
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}



	// 插入
	//s1 := Student{"小红", 12}
	//s2 := Student{"小兰", 10}
	//s3 := Student{"小黄", 11}
	//// 一个
	//insertResult, err := collection.InsertOne(context.TODO(), s1)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	//// 多个
	//students := []interface{}{s2, s3}
	//insertManyResult, err := collection.InsertMany(context.TODO(), students)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("inserted multiple documents:", insertManyResult)



	// 更新
	filter := bson.D{{"name","小兰"}}
	//update := bson.D{
	//	{"$inc",bson.D{
	//		{"age",1},
	//	}},
	//}
	//// 给小兰增加一岁
	//updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)



	// 查找文档
	var result Student
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)


	// 删除名字是小黄的那个
	//deleteResult1, err := collection.DeleteOne(context.TODO(), bson.D{{"name","小黄"}})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult1.DeletedCount)
	//// 删除所有
	//deleteResult2, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult2.DeletedCount)


	// 断开连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
