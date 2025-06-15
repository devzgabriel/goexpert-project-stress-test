package stress_test

import (
	"sync"
	"time"
)

type Params struct {
	URL         string
	Requests    int
	Concurrency int
}

func Run(config Params) *Report {

	startTime := time.Now()

	results := make(chan Result, config.Requests)

	concurrencyLimiter := make(chan struct{}, config.Concurrency)

	var wg sync.WaitGroup

	for i := 0; i < config.Requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Controla a concorrência
			concurrencyLimiter <- struct{}{}
			defer func() { <-concurrencyLimiter }()

			result := MakeRequest(config.URL)
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

	return report
}
