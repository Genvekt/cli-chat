package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Genvekt/cli-chat/services/chat-client/internal/app"
)

const (
	connectFlagChatID      = "chat_id"
	connectFlagChatIDShort = "i"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to chat",
	RunE: func(cmd *cobra.Command, args []string) error {
		chatID, err := cmd.Flags().GetInt64(connectFlagChatID)
		if err != nil {
			return fmt.Errorf("failed to read flag %s: %v", connectFlagChatID, err)
		}

		err = app.Connect(cmd.Context(), chatID)
		if err != nil {
			return fmt.Errorf("error occurred during connection: %v", err)
		}

		return nil
	},
}

func init() {
	connectCmd.Flags().Int64P(connectFlagChatID, connectFlagChatIDShort, 0, "Id of a chat")
	err := connectCmd.MarkFlagRequired(connectFlagChatID)
	if err != nil {
		os.Exit(1)
	}
}
