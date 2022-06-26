package ran

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Definition struct {
	Env      EnvironmentVariables
	Commands map[string]Command
}

type Command struct {
	Name        string
	Description string
	Tasks       []Task
}
type Task struct {
	Name   string
	Script string
	When   []string
	Env    map[string]string
	Defer  string
	Call   CommandCall
}

type CommandCall struct {
	Command string
}

type EnvironmentVariables []string

func LoadDefinition(filename string) (Definition, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Definition{}, err
	}
	defer file.Close()
	return ParseDefinition(file)
}

func ParseDefinition(r io.Reader) (Definition, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return Definition{}, err
	}

	var raw struct {
		Env      map[string]string `yaml:"env"`
		Commands map[string]struct {
			Description string `yaml:"description"`
			Tasks       []struct {
				Name   string            `yaml:"name"`
				Script string            `yaml:"script"`
				When   []string          `yaml:"when"`
				Env    map[string]string `yaml:"env"`
				Defer  string            `yaml:"defer"`
				Call   struct {
					Command string `yaml:"command"`
				} `yaml:"call"`
			} `yaml:"tasks"`
		} `yaml:"commands"`
	}
	if err := yaml.Unmarshal(bs, &raw); err != nil {
		return Definition{}, err
	}

	def := Definition{
		Env:      appendEnv(os.Environ(), raw.Env),
		Commands: make(map[string]Command, len(raw.Commands)),
	}

	for name, c := range raw.Commands {
		tasks := make([]Task, len(c.Tasks))
		for i, t := range c.Tasks {
			tasks[i] = Task{
				Name:   t.Name,
				Script: t.Script,
				When:   t.When,
				Env:    t.Env,
				Defer:  t.Defer,
				Call: CommandCall{
					Command: t.Call.Command,
				}}
		}
		def.Commands[name] = Command{
			Name:        name,
			Description: c.Description,
			Tasks:       tasks,
		}
	}
	return def, nil
}

func appendEnv(env EnvironmentVariables, m map[string]string) EnvironmentVariables {
	for k, v := range m {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	// resize capacity to len(env) to prevent conflict when append values from multiple tasks.
	return env[:len(env):len(env)]
}
