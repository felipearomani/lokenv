package examples_test

import (
	"context"
	"testing"

	"github.com/felipearomani/lokenv"
)

func TestProject(t *testing.T) {
	ctx := context.Background()
	project := lokenv.NewProject("my project")

	project.RegisterApp(ctx, lokenv.App{
		Name:    "users-service",
		WorkDir: "~/mockprojects/users-service",

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
