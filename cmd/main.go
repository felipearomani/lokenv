package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// LogStream is a interface to log the stdout and stderr
type LogStream interface {
	io.Writer
	Listen() <-chan []byte
}

// EnvironmentVars is a map of environment variables
type EnvironmentVars map[string]string

// ToSlice converts the environment variables to a slice of strings
func (e EnvironmentVars) ToSlice() []string {
	envSlice := make([]string, 0, len(e))
	for k, v := range e {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}
	return envSlice
}

type StdoutWriter struct {
	Stream chan []byte
}

func (s *StdoutWriter) Write(p []byte) (int, error) {
	s.Stream <- p
	return len(p), nil
}

func (s *StdoutWriter) Listen() <-chan []byte {
	return s.Stream
}

type App struct {
	Name      string
	PWD       string
	Command   []string
	Env       EnvironmentVars
	LogStream LogStream
	Process   *exec.Cmd
}

func (a *App) Run(ctx context.Context) error {
	if len(a.Command) == 0 {
		return fmt.Errorf("command not set")
	}

	a.Process = exec.CommandContext(ctx, a.Command[0], a.Command[1:]...)
	a.Process.Dir = a.PWD
	a.Process.Stdout = a.LogStream
	a.Process.Stderr = a.LogStream

	a.Process.Env = os.Environ()
	myEnvVars := a.Env.ToSlice()
	a.Process.Env = append(a.Process.Env, myEnvVars...)

	// Set the process group ID to the PID of the child process
	// I added that because when we run 'go run ...' code for example
	// it will throw the child process that is the binary, so it will add group id in process
	// than we can kill the entire group process furthermore
	a.Process.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := a.Process.Start(); err != nil {
		return fmt.Errorf("error on start: %w", err)
	}

	return nil
}

func (a *App) Stop() error {
	fmt.Printf("[%s] finishing app\n", a.Name)

	if err := syscall.Kill(-a.Process.Process.Pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("error on cancel the command: %w", err)
	}

	fmt.Printf("[%s] exited success\n", a.Name)

	return nil
}

func (a *App) StreamLogs(ctx context.Context) {
	logs := a.LogStream.Listen()
	for {
		select {
		case <-ctx.Done():
			return
		case log := <-logs:
			fmt.Printf("[%v]: %s\n", a.Name, string(log))
		}
	}
}

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
	)
	defer cancel()

	app := App{
		Name:    "go-hello",
		PWD:     "/Users/romani/studies/go-hello",
		Command: []string{"go", "run", "cmd/hello/main.go"},
		Env: EnvironmentVars{
			"DB_XPTO": "postgres://user@password:my-db-url:5432",
		},
		LogStream: &StdoutWriter{
			Stream: make(chan []byte),
		},
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("error on run: %v\n", err)
	}

	go app.StreamLogs(ctx)

	chQuit := make(chan os.Signal, 1)
	signal.Notify(chQuit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-chQuit

	if err := app.Stop(); err != nil {
		log.Fatalf("error on stop: %v\n", err)
	}
}

func findCommand(command string) {
	path, err := exec.LookPath(command)
	if err != nil {
		log.Fatalf("installing %v is in your future\n", command)
	}
	fmt.Printf("%v is available at %s\n", command, path)
}
