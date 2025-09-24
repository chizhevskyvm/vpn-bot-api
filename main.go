package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"vpn-bot-api/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		fmt.Println("🖥️ Service is running. Press CTRL+C to stop.")

		err := cmd.RunBot(ctx)
		if err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	fmt.Println("🛑 Shutting down gracefully...")
}
