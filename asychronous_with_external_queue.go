package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

func masterQueue(files []string) {
	if os.Getenv("WORKER") != "" {
		go workerQueue()
	}

	for _, file := range files {
		if err := redisClient.Publish(context.Background(), channelName, file).Err(); err != nil {
			panic(err)
		}
	}

	fmt.Println(reduceQueue(files))
}

func workerQueue() {
	ctx := context.Background()
	subscriber := redisClient.Subscribe(ctx, channelName)

	fmt.Println(fmt.Sprintf("WORKER %s - Waiting for messages", os.Getenv("WORKER")))
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println("Received message from " + msg.Channel + " channel.")

		result := processFile(msg.Payload)
		fmt.Println(fmt.Sprintf("WORKER %s - Result of file %s is: #+%v", os.Getenv("WORKER"), msg.Payload, result))

		content, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		if err := redisClient.Publish(context.Background(), channelNameProcessed, content).Err(); err != nil {
			panic(err)
		}
	}
}

func reduceQueue(files []string) map[string]int {
	var results []map[string]int
	ctx := context.Background()
	subscriber := redisClient.Subscribe(ctx, channelNameProcessed)

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		var content map[string]int
		err = json.Unmarshal([]byte(msg.Payload), &content)
		if err != nil {
			panic(err)
		}

		results = append(results, content)
		if len(results) == len(files) {
			break
		}
	}

	result := map[string]int{}

	for _, counts := range results {
		for word, count := range counts {
			result[word] += count
		}
	}

	return result
}
