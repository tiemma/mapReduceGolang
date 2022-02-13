package main

import "fmt"

func masterSync(files []string) {
	results := []map[string]int{}
	for _, file := range files {
		results = append(results, workerSync(file))
	}

	fmt.Println(reduceSync(results))
}

func workerSync(file string) map[string]int {
	return processFile(file)
}

func reduceSync(results []map[string]int) map[string]int {
	result := map[string]int{}

	for _, counts := range results {
		for word, count := range counts {
			result[word] += count
		}
	}

	return result
}
