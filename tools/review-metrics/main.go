// Command review-metrics reports code review metrics for GitHub repositories.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/example/review-metrics/cmd"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := cmd.NewRootCmd().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "review-metrics:", err)
		return 1
	}
	return 0
}
