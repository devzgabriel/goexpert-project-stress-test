package stress_test

import (
	"fmt"
	"time"
)

type Report struct {
	TotalRequests    int
	TotalTime        time.Duration
	StatusCodeCounts map[int]int
	SuccessfulReqs   int
}

func PrintHeader(url string, requests int, concurrency int) {
	fmt.Printf("🚀 Starting load test...\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Total requests: %d\n", requests)
	fmt.Printf("Concurrency: %d\n", concurrency)
	fmt.Println("==========================================")
}

func PrintReport(report *Report) {
	fmt.Println("🎯 Result:")
	fmt.Printf("Total execution time: %v\n", report.TotalTime)
	fmt.Printf("Total requests made: %d\n", report.TotalRequests)

	fmt.Println("\nStatus code distribution:")
	for statusCode, count := range report.StatusCodeCounts {
		percentage := float64(count) / float64(report.TotalRequests) * 100
		fmt.Printf("  Status %d: %d requests (%.2f%%)\n", statusCode, count, percentage)
	}

	if report.TotalRequests > 0 {
		successRate := float64(report.SuccessfulReqs) / float64(report.TotalRequests) * 100
		fmt.Printf("\nSuccess rate: %.2f%%\n", successRate)

		avgTime := report.TotalTime / time.Duration(report.TotalRequests)
		fmt.Printf("Average time per request: %v\n", avgTime)

		requestsPerSecond := float64(report.TotalRequests) / report.TotalTime.Seconds()
		fmt.Printf("Requests per second: %.2f\n", requestsPerSecond)
	}
}
