package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a new instance of our bot.
	bot, err := NewBot()
    // defer bot.client.Disconnect()
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	// Start the bot service in a separate goroutine.
	go func() {
		if err := bot.Start(); err != nil {
			log.Fatalf("Bot encountered an error: %v", err)
		}
	}()

	// Wait for termination signal (Ctrl+C, SIGTERM, etc.).
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down bot...")
}
