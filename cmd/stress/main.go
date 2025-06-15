package main

import (
	"fmt"
	"os"

	"github.com/devzgabriel/stress-test/internal/stress_test"
	"github.com/spf13/cobra"
)

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
	config := stress_test.Params{
		URL:         url,
		Requests:    requests,
		Concurrency: concurrency,
	}
	stress_test.PrintHeader(config.URL, config.Requests, config.Concurrency)
	report := stress_test.Run(config)
	stress_test.PrintReport(report)
}
