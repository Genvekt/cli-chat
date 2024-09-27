package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Genvekt/cli-chat/services/chat-client/cmd/cli/create"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/app"
)

var (
	profileName string
	iniConfig   string
	envConfig   string
)

var (
	rootCmd = &cobra.Command{
		Use:   "cli-chat",
		Short: "CLI Chat application",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// We initialise application before execution any command
			var err error

			err = app.InitApp(cmd.Context(), profileName, iniConfig, envConfig)
			if err != nil {
				return fmt.Errorf("failed to initialise application: %v", err)
			}

			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&profileName, "profile", "", "profile name")
	rootCmd.PersistentFlags().StringVar(&iniConfig, "ini-config", "config.ini", "ini config path")
	rootCmd.PersistentFlags().StringVar(&envConfig, "env-config", ".env", "env config path")

	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(connectCmd)
}

// Execute acts as application entrypoint
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
