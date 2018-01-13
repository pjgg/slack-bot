package commands

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var (
	helloWorldCmd = &cobra.Command{
		Use: "hello",
		Run: helloWorldHandler,
	}
)

func init() {
	BaseCmd.AddCommand(helloWorldCmd)
}

func helloWorldHandler(cmd *cobra.Command, args []string) {
	// Run stuff as K8s command. 
	fmt.Println("Hello World Executed!")
	os.Exit(0)
}