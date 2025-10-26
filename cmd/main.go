package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfg Config
var rootCmd = &cobra.Command{
	Use:   "go-grep [target]",
	Short: "go-grep — grep like app written with go",
	Args:  cobra.MaximumNArgs(1),
	RunE:  onRootCommandRun,
}

type Config struct {
	target string

	isRecursive bool
	isVerbose   bool
}

func GetConfig() Config {
	return cfg
}

func (c Config) GetTarget() string {
	return c.target
}

func (c Config) GetIsRecursive() bool {
	return c.isRecursive
}

func (c Config) GetIsVerbose() bool {
	return c.isVerbose
}

func Main() {
	main()
}

func init() {
	rootCmd.Flags().StringVarP(&cfg.target, "target", "t", "./", "target")
	rootCmd.Flags().BoolVarP(&cfg.isRecursive, "recursive", "r", false, "is recursive")
	rootCmd.Flags().BoolVarP(&cfg.isVerbose, "verbose", "v", false, "verbosity")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func onRootCommandRun(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		cfg.target = args[0]
	}
	return nil
}
