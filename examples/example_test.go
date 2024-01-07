package examples_test

import (
	"testing"

	"github.com/felipearomani/lokenv"
)

func TestProject(t *testing.T) {

	project := lokenv.NewProject()

	project.RegisterApp(lokenv.App{
		Name:    "users-service",
		Workdir: "~/mockprojects/users-service",
		Command: []string{"go", "run", "cmd/main.go"},
		Env: lokenv.Variables{
			"APP_NAME": "GO APP API REST",
			"DB_URL":   "localhost:5432",
		},
	})

	if err := project.Start(); err != nil {
		t.Fatal(err)
	}
}
