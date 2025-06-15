package stress_test

import (
	"context"
	"net/http"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
}

func MakeRequest(url string) Result {
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
