package main

//Average time/run with 32 workers: 5.23 ms

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"
	"unicode"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument.
func WordCount(text string, numWorkers int) map[string]int {

	//Split the text into chunks
	chunks := make([]string, numWorkers)
	words := strings.Fields(text)
	chunkSize := (len(words) + numWorkers - 1) / numWorkers
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunks[i] = strings.Join(words[start:end], " ")
	}

	//Map stage: count word frequencies in each chunk
	results := make(chan map[string]int, numWorkers)
	var wg sync.WaitGroup
	for _, chunk := range chunks {
		wg.Add(1)
		go func(chunk string) {
			defer wg.Done()
			freqs := make(map[string]int)
			words := strings.Fields(chunk)
			for _, word := range words {
				word = strings.TrimFunc(word, func(r rune) bool {
					return unicode.IsPunct(r) || unicode.IsSymbol(r)
				})
				word = strings.ToLower(word)
				freqs[word]++
			}
			results <- freqs
		}(chunk)
	}
	wg.Wait()
	close(results)

	//Reduce stage: combine word frequencies from all chunks
	freqs := make(map[string]int)
	for partial := range results {
		for word, count := range partial {
			freqs[word] += count
		}
	}

	return freqs
}

// Benchmark how long it takes to count word frequencies in text numRuns times,
// using numWorkers goroutines.
//
// Return the total time elapsed.
func benchmark(text string, numRuns, numWorkers int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text, numWorkers)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {
	// read in DataFile as a string called data
	data, err := ioutil.ReadFile(DataFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Benchmark with different numbers of workers
	for _, numWorkers := range []int{1, 2, 4, 8, 16, 32, 64, 128} {
		fmt.Printf("Running benchmark with %d workers...\n", numWorkers)
		numRuns := 100
		runtimeMillis := benchmark(string(data), numRuns, numWorkers)
		printResults(runtimeMillis, numRuns)
	}
}
