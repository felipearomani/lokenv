package lokenv

type (
	Variables map[string]string

	Project struct {
		apps []App
	}

	App struct {
		Name    string
		Workdir string
		Command []string
		Env     Variables
	}
)

func NewProject() *Project {
	return &Project{}
}

func (p *Project) RegisterApp(app App) {
	p.apps = append(p.apps, app)
}

func (p *Project) Start() error {
	return nil
}
