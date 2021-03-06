package ran

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type App struct{}

func NewApp() App {
	return App{}
}

func (app App) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	rootCmd := &cobra.Command{
		Use:           "ran",
		SilenceErrors: true,
	}
	rootCmd.SetArgs(args[1:])

	file := rootCmd.PersistentFlags().StringP("file", "f", "ran.yaml", "ran definition file.")
	logLevel := rootCmd.PersistentFlags().String("log-level", "info", "log level(debug, info, error, discard)")

	if err := rootCmd.PersistentFlags().Parse(args); err != nil && err != pflag.ErrHelp {
		rootCmd.Usage()
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	level, err := NewLogLevel(*logLevel)
	if err != nil {
		rootCmd.Usage()
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	logger := NewStdLogger(os.Stdout, level)

	def, err := LoadDefinition(*file)
	if err != nil {
		rootCmd.Usage()
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	abs, err := filepath.Abs(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	workingDir := filepath.Dir(abs)

	commandRunner := NewStdCommandRunner(def.Commands, workingDir, logger)

	for _, c := range def.Commands {
		rootCmd.AddCommand(&cobra.Command{
			Use:   c.Name,
			Short: c.Description,
			RunE: func(cmd *cobra.Command, args []string) error {
				return commandRunner.RunCommand(cmd.Use, RuntimeEnvironment{
					Stdin:            os.Stdin,
					Stdout:           os.Stdout,
					Stderr:           os.Stderr,
					Env:              def.Env,
					WorkingDirectory: "",
				})
			},
			SilenceErrors: true,
			SilenceUsage:  true,
		})
	}

	if err := rootCmd.Execute(); err != nil {
		logger.Error("%s", err.Error())
		return 1
	}
	return 0
}
