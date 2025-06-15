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
	fmt.Printf("🚀 Iniciando teste de carga...\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Total de requests: %d\n", requests)
	fmt.Printf("Concorrência: %d\n", concurrency)
	fmt.Println("==========================================")
}

func PrintReport(report *Report) {
	fmt.Println("🎯 Resultado:")
	fmt.Printf("Tempo total de execução: %v\n", report.TotalTime)
	fmt.Printf("Total de requests realizados: %d\n", report.TotalRequests)
	fmt.Printf("Requests com status 200: %d\n", report.SuccessfulReqs)

	fmt.Println("\nDistribuição de códigos de status:")
	for statusCode, count := range report.StatusCodeCounts {
		percentage := float64(count) / float64(report.TotalRequests) * 100
		fmt.Printf("  Status %d: %d requests (%.2f%%)\n", statusCode, count, percentage)
	}

	if report.TotalRequests > 0 {
		successRate := float64(report.SuccessfulReqs) / float64(report.TotalRequests) * 100
		fmt.Printf("\nTaxa de sucesso: %.2f%%\n", successRate)

		avgTime := report.TotalTime / time.Duration(report.TotalRequests)
		fmt.Printf("Tempo médio por request: %v\n", avgTime)

		requestsPerSecond := float64(report.TotalRequests) / report.TotalTime.Seconds()
		fmt.Printf("Requests por segundo: %.2f\n", requestsPerSecond)
	}
}
