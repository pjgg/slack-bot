package commands

import (
	
	"os"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
)

var (
	// BaseCmd represents the base command when called without any subcommands
	BaseCmd = &cobra.Command{
		Use:   "basecmd",
		Short: "",
		Long:  "",
	}
	ConfPath string
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := BaseCmd.Execute(); err != nil {
		logrus.Error("BaseCmd " + err.Error())
		os.Exit(-1)
	}
}

func init() {
	BaseCmd.PersistentFlags().StringVarP(&ConfPath, "conf", "c", "", "path to config")
}
