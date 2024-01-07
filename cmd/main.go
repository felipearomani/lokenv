package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	// findCommand("git")

	var (
		ctx, cancel = context.WithCancel(context.Background())
	)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", "cmd/hello/main.go")
	cmd.Dir = "/Users/romani/studies/go-hello"
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	cmd.Env = os.Environ()
	myEnvVars := []string{"DB_XPTO=postgres://user@password:my-db-url:5432"}
	cmd.Env = append(cmd.Env, myEnvVars...)

	// Set the process group ID to the PID of the child process
	// I added that becaus when we run 'go run ...' code for example
	// it will throw the a child process that is the binary, so it will add group id in process
	// than we can kill the entire group process furthermore
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		log.Fatalf("error on start: %v\n", err)
	}
	fmt.Printf("process created: %v\n", cmd.Process.Pid)

	<-time.After(10 * time.Second)

	fmt.Println("send kill process")
	// Send SIGTERM to the entire process group
	// Here, I sent a negative PID that will inform the SO to kill all groups, including the children.
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
		fmt.Printf("error on cancel the command: %v\n", err)
	}

	// this for is to mock a a infinity loop aplication withoud deadlock the main goroutine
	for {
		time.Sleep(5 * time.Second)
	}
}

// func findCommand(command string) {
// 	path, err := exec.LookPath(command)
// 	if err != nil {
// 		log.Fatalf("installing %v is in your future\n", command)
// 	}
// 	fmt.Printf("%v is available at %s\n", command, path)
// }
