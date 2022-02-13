package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Word count
// You have multiple files with single words line by line
// Count each word and give a total summary of repeating words

// Master read the file
// word word2 word word3 word2
// word2 word2 word

// Workers
// {word: 2, word2: 2, word3: 1} - Worker 1
// {word: 1, word2: 2}           - Worker 2

// Reducer: Sums up all the results
// {word: 3, word2: 4, word3: 1}

// Map reduce has three stages
// Master -> Worker[] -> Reducer

// Concurrency model
// Golang - Channels (Inter Process Communication, message queue)

// Distributed computing model
// synchronization - hardware implements locks via instruction sets CMPXCHG
// waitGroup - sync package

// go command - goroutine
// wg.Add - semaphore (counting mutexes)

// Deadlock
// wg needs go routines to be running
// wg - 2 goroutines (watches them)
// wg - 2 goroutines but the go runtime has no go routines
// wg - Nobody called done but you did an Add
// wg panic that no one called done so we have a deadlock

// IPC - Inter Process Communication
// parent and children
// children are independent
// how do the children communicate with each other
// communicate with a message queue

// Multiple instances of main.go running
// How would we manage scheduling tasks?

// We implement a shared filesystem with volumes
// This is so we only care about delegating the tasks
// And not who processes them

var (
	dirPath              = "fixtures"
	channelName          = "send_stuff"
	channelNameProcessed = "processed_stuff"
	wg                   = &sync.WaitGroup{}
	redisClient          = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", getRedisHostOrDefault()),
	})
)

func getRedisHostOrDefault() string {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		return "localhost"
	}

	return redisHost
}

func main() {
	// Walk through the directory and gather the list of file names
	fileInfos, err := ioutil.ReadDir(dirPath)

	var files []string
	if err != nil {
		panic(err)
	}
	for _, fileInfo := range fileInfos {
		files = append(files, fileInfo.Name())
	}

	if os.Getenv("MASTER") != "" || os.Getenv("WORKER") != "" {
		benchmark(files, masterQueue, "asynchronous_with_external_queue")
	} else {
		benchmark(files, masterSync, "synchronous")
		benchmark(files, masterAsync, "asynchronous")
	}
}

func benchmark(files []string, method func([]string), tag string) {
	start := time.Now()
	method(files)
	end := time.Now()
	fmt.Println(fmt.Sprintf("Running a %s process\nThis started %s, ended %s,  %#+v seconds", tag, start.Format(time.RFC822), end.Format(time.RFC822), end.Unix()-start.Unix()))
	fmt.Println("---------------------------------------\n")
}

func processFile(file string) map[string]int {
	filePath := dirPath + "/" + file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	count := map[string]int{}
	for _, word := range strings.Split(string(content), "\n") {
		count[word] += 1
	}

	// Delay to show that it takes time
	time.Sleep(2 * time.Second)

	return count
}
