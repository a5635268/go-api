package cmd

import (
	"errors"
	"fmt"
	"go-api/cmd/api"
	"go-api/cmd/test"
	"os"

	"github.com/spf13/cobra"
	"go-api/cmd/config"
	"go-api/common/global"
	"go-api/tools"
)

var rootCmd = &cobra.Command{
	Use:          "go-api",
	Short:        "go-api",
	// SilenceUsage: true,
	Long:         `go-api`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New(tools.Red("requires at least one arg"))
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := `欢迎使用 ` + tools.Green(`go-api `+global.Version) + ` 可以使用 ` + tools.Red(`-h`) + ` 查看命令`
	fmt.Printf("%s\n", usageStr)
}

func init() {
	rootCmd.AddCommand(test.StartCmd)
	rootCmd.AddCommand(api.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
