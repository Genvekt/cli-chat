package create

import "github.com/spf13/cobra"

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create something",
}

func init() {
	CreateCmd.AddCommand(chatCmd)
	CreateCmd.AddCommand(messageCmd)
}
