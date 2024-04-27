package main

import (
	"importMessengerToSignal/internal"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	internal.EnsureSignalCliBinary()
	dbus := internal.StartDbus()

	time.Sleep(2 * time.Second) // wait for dbus to start

	defer internal.StopDbus(dbus)

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
