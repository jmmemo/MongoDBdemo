package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	//	set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	//连接Mongodb	Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	//检查连接	check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//拿到handle
	collection := client.Database("test").Collection("trainers")

	//关闭连接
	//err = client.Disconnect(context.TODO())
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Connection to MongoDB closed")

	//创建几个结构体用于数据库插入
	ash := Trainer{
		Name: "Ash",
		Age:  10,
		City: "Pallet Town",
	}
	misty := Trainer{
		Name: "Misty",
		Age:  10,
		City: "Cerulean City",
	}
	brock := Trainer{
		Name: "Brock",
		Age:  15,
		City: "Pewter City",
	}

	//插入一个
	insertOneResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertOneResult.InsertedID)

	//插入多个
	trainers := []interface{}{misty, brock}
	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	//筛选器--用于匹配数据库中的文档
	filter := bson.D{{"name", "Ash"}}
	//更新文档--用于描述更新操作
	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	//updateOne--通过filter找到name字段的Ash，通过update操作对age加1
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched查找到了 %v documents and updated更新了 %v documents\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	//FindOne--通过filter查找单个文档，取新建的result的地址，暂存其内
	var result Trainer //声明暂存

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("找到单个文档Found a single documents:%+v\n", result)

	//Find--查找多个文档
	findOptions := options.Find()
	findOptions.SetLimit(2) //限制返回的文档数为2,即只查找前2个

	//同上声明暂存数组--results
	var results []*Trainer

	//拿到游标
	cursor, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	//每个游标允许我们一次解码一个文档
	for cursor.Next(context.TODO()) {
		//声明暂存元素
		var elem Trainer
		//对每个游标解码，暂存到elem
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		//把每个elem加入到数组results中
		results = append(results, &elem)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	//
	cursor.Close(context.TODO())

	fmt.Printf("找到的文档数组Found multiple documents(array of pointers): %+v\n", results)
	//
	for i, trainer := range results {
		fmt.Println(i, trainer)
	}

	//删除多个
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}
