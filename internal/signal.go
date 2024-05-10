package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const SignalCliBinary = "./signal-cli"

var dbusInstance *exec.Cmd
var sendMessagesOnCurrentDBusInstance = 0

func EnsureSignalCliBinary() {
	if _, err := os.Stat(SignalCliBinary); os.IsNotExist(err) {
		log.Fatalln("signal-cli binary not found, please download it from https://github.com/AsamK/signal-cli/releases/latest and place it in the current directory named signal-cli")
	}
}

func SendMessageToSignal(from string, to string, timestamp int64, message string, attachments []string, isGroup bool) {
	args := []string{"--dbus", "-u", from, "send"}

	if isGroup {
		args = append(args, "-g")
	}

	args = append(args, to)

	for _, attachment := range attachments {
		args = append(args, "--attachment", filepath.Join(GetConfig().Entrypoint, attachment))
	}

	args = append(args, "-m", fmt.Sprintf("[%s] %s", time.UnixMilli(timestamp).Format(time.DateTime), message))

	cmd := exec.Command(SignalCliBinary, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Couldn't send message to %s: %v, output: %s", to, err, string(output))
	}
}

func StartDbus() {
	cmd := exec.Command(SignalCliBinary, "daemon", "--dbus", "--receive-mode", "manual")
	err := cmd.Start()

	if err != nil {
		log.Fatalf("Couldn't start DBus: %v", err)
	}

	log.Printf("Started DBus with pid %d", cmd.Process.Pid)

	dbusInstance = cmd

	time.Sleep(10 * time.Second) // wait for dbus to boot up
}

func StopDbus() {
	if err := dbusInstance.Process.Signal(syscall.SIGTERM); err != nil {
		log.Printf("Failed to stop DBus: %v", err)
	}

	if _, err := dbusInstance.Process.Wait(); err != nil {
		log.Printf("Failed to wait for DBus to stop: %v", err)
	}

	log.Println("DBus has been stopped.")

	dbusInstance = nil
	sendMessagesOnCurrentDBusInstance = 0
}

func RestartDbus() {
	log.Println("Restarting DBus")
	StopDbus()
	StartDbus()
}

func IncreaseSendMessagesOnCurrentDBusInstance() {
	sendMessagesOnCurrentDBusInstance++
	if sendMessagesOnCurrentDBusInstance >= GetConfig().RestartDbusEvery {
		RestartDbus()
	}
}
