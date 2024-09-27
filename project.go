package lokenv

import (
	"context"
	"os/exec"
	"strings"
)

type (
	Variables map[string]string

	Project struct {
		apps []App
	}

	AppConfig struct {
		Name         string
		WorkDir      string
		PreCommands  []string
		PostCommands []string
		RunCommand   string
		Environment  Variables
	}

	App struct {
		preCmd  []*exec.Cmd
		postCmd []*exec.Cmd
		mainCmd *exec.Cmd
	}
)

func NewProject() *Project {
	return &Project{}
}

func (p *Project) RegisterApp(ctx context.Context, cfg AppConfig) {

	// main command
	run, args := getCommand(cfg.RunCommand)
	mainCmd := exec.CommandContext(ctx, run, args...)
	mainCmd.Dir = cfg.WorkDir
	//TODO add env
	//TODO add stdout

	app := App{
		mainCmd: mainCmd,
	}

	p.apps = append(p.apps, app)
}

func (p *Project) Start() error {
	return nil
}

func getCommand(c string) (string, []string) {
	slicedCommand := strings.Split(c, " ")
	runCommand := slicedCommand[0]

	var args []string
	copy(args, slicedCommand[1:])

	return runCommand, args
}
