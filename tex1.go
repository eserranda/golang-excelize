// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/redis/go-redis/v9"
// )

// var client = redis.NewClient(&redis.Options{
// 	Addr: "localhost:6379",
// 	DB:   0,
// })

// var ctx = context.Background()

// func main() {
// 	result, err := client.Ping(ctx).Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(result)

// 	client.RPush(ctx, "eser", "Eliaser randa", "Ganteng", "Ganteng Banget")

// 	resultData, err := client.LRange(ctx, "eser", 0, -1).Result()
// 	if err != nil {
// 		log.Fatal("Error getting data from Redis:", err)
// 	}

// 	fmt.Println("Hasil Data : ", resultData)
// }
// }