package main

import (
	"context"
	"fmt"
	"time"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 20,
	})
}

func doCommand() {
	ctx := context.Background()
	val, err := rdb.Get(ctx, "h1").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
	//先获取到命令对象
	cmder := rdb.Get(ctx, "h1")
	fmt.Println(cmder.Val())
	fmt.Println(cmder.Err())
	a,err := rdb.Do(ctx, "get", "key").Result()
	fmt.Println(a,err)
	//	直接执行命令获取错误
	//err = rdb.Set(ctx,"key",10,time.Hour).Err()
	//fmt.Println(err)
	//fmt.Println(rdb.Get(ctx,"key").Val())
}

// doDemo rdb.Do 方法使用示例
func doDemo() {
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 直接执行命令获取错误
	err := rdb.Do(ctx, "set", "keyy", 10, "EX", 36).Err()
	fmt.Println(err)

	// 执行命令获取结果
	val, err := rdb.Do(ctx, "get", "key").Result()
	fmt.Println(val, err)
}

//zset 示例
func zsetDemo()  {
	key := "language_rank"
	languages:= []*redis.Z{
		{Score: 90,Member: "go"},
		{Score: 98,Member: "java"},
		{Score: 95,Member: "python"},
		{Score: 97,Member: "js"},
		{Score: 99,Member: "c/c++"},
	}
	ctx := context.Background()

//	zadd
	err := rdb.ZAdd(ctx,key,languages...).Err()
	fmt.Println(err)
//	取分数最高的前三个语言
	ret := rdb.ZRevRangeWithScores(ctx,key,0,2).Val()
	for _, z := range ret {
		fmt.Println(z.Member,z.Score)
	}

	//给go加10分
	newscore := rdb.ZIncrBy(ctx,key,10.0,"go").Val()
	fmt.Println("go new score ",newscore)
	//取95-100分的
	op := &redis.ZRangeBy{
		Min: "95",
		Max: "100",
	}

	ret,_ = rdb.ZRangeByScoreWithScores(ctx,key,op).Result()
	for _, z := range ret {
		fmt.Println(z.Member,z.Score)
	}


}
func main() {
	initRedis()
	//doDemo()
	//doCommand()
	zsetDemo()

}
