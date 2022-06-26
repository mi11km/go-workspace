package ran

import "io"

type RuntimeEnvironment struct {
	Stdin            io.Reader
	Stdout           io.Writer
	Stderr           io.Writer
	Env              EnvironmentVariables
	WorkingDirectory string
}

type CommandRunner interface {
	RunCommand(command string, renv RuntimeEnvironment) error
}

type StdCommandRunner struct {
	commands map[string]Command
	workDir  string

	logger Logger
}

func NewStdCommandRunner(commands map[string]Command, workDir string, logger Logger) CommandRunner {
	return &StdCommandRunner{commands: commands, workDir: workDir, logger: logger}
}

func (s StdCommandRunner) RunCommand(command string, renv RuntimeEnvironment) error {
	// TODO implement
	return nil
}
