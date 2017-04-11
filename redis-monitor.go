package main

import (
	"fmt"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

func clearScreen() {
	// ANSI sequence
	fmt.Print("\033[H\033[2J")
}

func main() {
	address := ":6379"
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	r, err := redis.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Couldn't connect to Redis service at %s\n", address)
		os.Exit(1)
	}

	for ; ; time.Sleep(1 * time.Second) {
		clearScreen()

		keys, _ := redis.Strings(r.Do("KEYS", "*"))

		for _, key := range keys {

			valueType, _ := redis.String(r.Do("TYPE", key))

			fmt.Printf("%v (%v):\n", key, valueType)

			switch valueType {
			case "set":
				values, _ := redis.Strings(r.Do("SMEMBERS", key))
				for _, v := range values {
					fmt.Printf("  %v\n", v)
				}

			case "list":
				values, _ := redis.Strings(r.Do("LRANGE", key, 0, -1))
				for _, v := range values {
					fmt.Printf("  %v\n", v)
				}

			case "zset":
				values, _ := redis.Strings(r.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
				for i := 0; i < len(values)/2; i++ {
					fmt.Printf("  %v | %v\n", values[i*2+1], values[i*2])
				}

			case "string":
				value, _ := redis.String(r.Do("GET", key))
				fmt.Printf("  %v\n", value)

			case "hash":
				values, _ := redis.Strings(r.Do("HGETALL", key))
				for i := 0; i < len(values)/2; i++ {
					fmt.Printf("  %v | %v\n", values[i*2], values[i*2+1])
				}

			default:
				fmt.Printf("  (%v not supported)\n", valueType)
			}

			fmt.Println()
		}

	}
}
