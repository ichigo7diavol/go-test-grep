package cmd

import (
	"fmt"
	"os"

	"github.com/ichigo7diavol/go-test-grep/pkg/grep"
	"github.com/spf13/cobra"
)

const (
	MaximumNArgs = 2
)

var cfg grep.Config
var rootCmd = &cobra.Command{
	Use:   "go-grep [target]",
	Short: "go-grep — grep like app written with go",
	Args:  cobra.MaximumNArgs(MaximumNArgs),
	RunE:  onRootCommandRun,
}

func Main() {
	main()
}

func init() {
	rootCmd.Flags().StringVarP(&cfg.Expression, "expression", "e", "./", "regex expression")
	rootCmd.Flags().StringVarP(&cfg.Target, "target", "t", "./", "target")
	rootCmd.Flags().BoolVarP(&cfg.IsRecursive, "recursive", "r", false, "is recursive")
	rootCmd.Flags().BoolVarP(&cfg.IsVerbose, "verbose", "v", false, "is verbose")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseArgs(args []string) {
	if len(args) >= 1 {
		cfg.Expression = args[0]
	}
	if len(args) >= 2 {
		cfg.Target = args[1]
	}
}

func onRootCommandRun(cmd *cobra.Command, args []string) error {
	parseArgs(args)
	return grep.Execute(cfg)
}
