package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	go func() {
		c := time.Tick(5 * time.Second)
		for range c {
			RunEverySecond()
		}
	}()

	for {
	}

}

func RunEverySecond() {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	listName := os.Getenv("LIST_QUEUE_NAME")
	listApproved := os.Getenv("LIST_APPROVED_NAME")
	queueLists, err := client.LRange(listName, 0, 0).Result()
	if err != nil {
		fmt.Println(err)
	}
	if len(queueLists) > 0 {
		var Obj interface{}
		err = json.Unmarshal([]byte(queueLists[0]), &Obj)
		if err != nil {
			fmt.Println(err)
		}
		if Obj.(map[string]interface{})["name"] != nil && Obj.(map[string]interface{})["email"] != nil {
			Obj.(map[string]interface{})["status"] = "approved"
		} else {
			Obj.(map[string]interface{})["status"] = "declined"
		}
		res, err := json.Marshal(Obj)
		if err != nil {
			fmt.Println(err)
		}

		Lists, err := client.RPush(listApproved, res).Result()
		if err != nil {
			fmt.Println(err)
		}
		_, err = client.LPop(listName).Result()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("pushed to queue lists: ", Lists, ", value: ", string(res))
	}
	fmt.Println("empty queue, ", time.Now())
}
