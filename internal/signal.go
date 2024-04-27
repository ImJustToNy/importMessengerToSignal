package internal

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const SignalCliBinary = "./signal-cli"

func EnsureSignalCliBinary() {
	if _, err := os.Stat(SignalCliBinary); os.IsNotExist(err) {
		log.Fatalln("signal-cli binary not found, please download it from https://github.com/AsamK/signal-cli/releases/latest and place it in the current directory named signal-cli")
	}
}

func SendMessageToSignal(from string, to string, message string, attachments []string) {
	args := []string{"--dbus", "-u", from, "send", "-g", to}

	if len(attachments) > 0 {
		for _, attachment := range attachments {
			args = append(args, "--attachment", filepath.Join(GetConfig().Entrypoint, attachment))
		}
	}

	if message != "" {
		args = append(args, "-m", message)
	}

	cmd := exec.Command(SignalCliBinary, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Couldn't send message to %s: %v, output: %s", to, err, string(output))
	}
}

func StartDbus() *exec.Cmd {
	cmd := exec.Command(SignalCliBinary, "daemon", "--dbus")
	err := cmd.Start()

	if err != nil {
		log.Fatalf("Couldn't start dbus: %v", err)
	}

	log.Printf("Started dbus with pid %d", cmd.Process.Pid)

	return cmd
}

func StopDbus(dbus *exec.Cmd) {
	if err := dbus.Process.Signal(syscall.SIGTERM); err != nil {
		log.Printf("Failed to stop DBus: %v", err)
	}

	if _, err := dbus.Process.Wait(); err != nil {
		log.Printf("Failed to wait for DBus to stop: %v", err)
	}

	log.Println("DBus has been stopped.")
}
