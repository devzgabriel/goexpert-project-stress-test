package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
}

type Report struct {
	TotalRequests    int
	TotalTime        time.Duration
	StatusCodeCounts map[int]int
	SuccessfulReqs   int
}

var (
	url         string
	requests    int
	concurrency int
)

var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "Ferramenta de teste de carga para serviços web",
	Long:  "Uma ferramenta CLI em Go para realizar testes de carga em serviços web com relatórios detalhados",
	Run:   runStressTest,
}

func init() {
	rootCmd.Flags().StringVar(&url, "url", "", "URL do serviço a ser testado (obrigatório)")
	rootCmd.Flags().IntVar(&requests, "requests", 0, "Número total de requests (obrigatório)")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 1, "Número de chamadas simultâneas (obrigatório)")

	rootCmd.MarkFlagRequired("url")
	rootCmd.MarkFlagRequired("requests")
	rootCmd.MarkFlagRequired("concurrency")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runStressTest(cmd *cobra.Command, args []string) {
	fmt.Printf("Iniciando teste de carga...\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Total de requests: %d\n", requests)
	fmt.Printf("Concorrência: %d\n", concurrency)
	fmt.Println("==========================================")

	startTime := time.Now()

	results := make(chan Result, requests)

	// Controla o número de goroutines simultâneas
	sem := make(chan struct{}, concurrency)

	var wg sync.WaitGroup

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Controla a concorrência
			sem <- struct{}{}
			defer func() { <-sem }()

			result := makeRequest(url)
			results <- result
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	report := &Report{
		StatusCodeCounts: make(map[int]int),
	}

	for result := range results {
		report.TotalRequests++
		report.StatusCodeCounts[result.StatusCode]++

		if result.StatusCode == 200 {
			report.SuccessfulReqs++
		}
	}

	report.TotalTime = time.Since(startTime)

	printReport(report)
}

func makeRequest(url string) Result {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Result{
			StatusCode: 0,
			Duration:   time.Since(start),
		}
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return Result{
			StatusCode: 0,
			Duration:   time.Since(start),
		}
	}
	defer resp.Body.Close()

	return Result{
		StatusCode: resp.StatusCode,
		Duration:   time.Since(start),
	}
}

func printReport(report *Report) {
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
