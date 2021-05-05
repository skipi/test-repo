package main

import (
	"fmt"
	"io/ioutil"
	"path"
)

type Tasks []Task

type Task struct {
	Name string
	Jobs []Job
}

type Job struct {
	Name string
}

var baseYaml = `
version: v1.0
name: Test Results
agent:
  machine:
    type: e1-standard-2
    os_image: ubuntu2004
blocks:
  - name: Compile
    skip:
      when: "branch == 'master'"
    task:
      jobs:
        - name: compilation
          commands:
            - make run
`

func (me Task) toYAML() string {
	out := `
  - name: %s
    skip:
      when: "branch != 'master'"
    task:
      jobs:
  `
	out = fmt.Sprintf(out, me.Name)

	for _, job := range me.Jobs {
		out = fmt.Sprintf("%s\n%s", out, job.toYAML(me.Name))
	}
	return out
}

func (me Job) toYAML(dir string) string {
	out := `
        - name: %s
          commands:
            - checkout
            - test-results publish %s/%s
  `
	out = fmt.Sprintf(out, me.Name, dir, me.Name)
	return out
}

func (me Tasks) toYAML() string {
	out := ""
	for _, task := range me {
		out = fmt.Sprintf("%s\n%s", out, task.toYAML())
	}
	return out
}

func main() {
	baseDir := "./src"

	files, _ := ioutil.ReadDir(baseDir)
	tasks := Tasks{}
	for _, f := range files {
		task := Task{Name: f.Name()}
		switch f.IsDir() {
		case true:
			files, _ := ioutil.ReadDir(path.Join(baseDir, f.Name()))
			for _, sf := range files {
				switch sf.IsDir() {
				case false:
					task.Jobs = append(task.Jobs, Job{Name: sf.Name()})
				}
			}
		case false:
			job := Job{Name: f.Name()}
			task.Jobs = append(task.Jobs, job)
		}
		tasks = append(tasks, task)
	}
	out := fmt.Sprintf("%s%s", baseYaml, tasks.toYAML())
	fmt.Print(out)
}
