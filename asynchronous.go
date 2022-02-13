package main

import "fmt"

func masterAsync(files []string) {
	ch := make(chan map[string]int)
	for _, file := range files {
		wg.Add(1)
		go workerAsync(file, ch) // goroutine
	}
	fmt.Println(reduceAsync(ch))
}

func workerAsync(file string, ch chan map[string]int) {
	ch <- processFile(file)
	wg.Done()
}

func reduceAsync(ch chan map[string]int) map[string]int {
	go func() {
		wg.Wait()
		close(ch)
	}()

	result := map[string]int{}

	for counts := range ch {
		for word, count := range counts {
			result[word] += count
		}
	}

	return result
}
