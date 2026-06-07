package main

import (
	"context"
	"fmt"
	"time"

	"github.com/istvzsig/ratelx"
)

func main() {

	// 3 tokens max, 1 token/sec refill
	limiter := ratelx.New(3, 1.0, false)

	fmt.Println("=== BURST MODE DEMO ===")

	// simulate burst traffic
	for i := 1; i <= 10; i++ {
		if limiter.Allow() {
			fmt.Printf("request %d → ALLOWED\n", i)
		} else {
			fmt.Printf("request %d → RATE LIMITED\n", i)
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("\n=== RECOVERY OVER TIME ===")

	// wait for refill
	time.Sleep(2 * time.Second)

	for i := 1; i <= 5; i++ {
		if limiter.Allow() {
			fmt.Printf("recovery request %d → ALLOWED\n", i)
		} else {
			fmt.Printf("recovery request %d → RATE LIMITED\n", i)
		}
	}

	fmt.Println("\n=== BLOCKING MODE DEMO ===")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 1; i <= 5; i++ {
		start := time.Now()
		err := limiter.Wait(ctx)

		elapsed := time.Since(start)

		fmt.Printf(
			"wait request %d → %v (waited %v)\n",
			i,
			err,
			elapsed,
		)

		fmt.Printf("wait request %d → %v (waited %v)\n", i, err, elapsed)
	}

}
