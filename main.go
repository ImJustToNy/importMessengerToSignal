package main

import (
	"importMessengerToSignal/internal"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	internal.EnsureSignalCliBinary()
	internal.StartDbus()

	defer internal.StopDbus()

	go func() {
		for _, conversation := range internal.GetConfig().Conversations {
			internal.ProcessConversation(&conversation)
		}
	}()

	handleSignals()
	log.Println("Bye bye!")
}

func handleSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals // Block until a signal is received
}
