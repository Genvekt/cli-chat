package create

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Genvekt/cli-chat/services/chat-client/internal/app"
)

const (
	msgFlagChatID      = "chat_id"
	msgFlagChatIDShort = "i"

	msgFlagMessage      = "message"
	msgFlagMessageShort = "m"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Create new message",
	RunE: func(cmd *cobra.Command, args []string) error {
		chatID, err := cmd.Flags().GetInt64(msgFlagChatID)
		if err != nil {
			return fmt.Errorf("failed to read flag %s: %v", msgFlagChatID, err)
		}

		message, err := cmd.Flags().GetString(msgFlagMessage)
		if err != nil {
			return fmt.Errorf("failed to read flag %s: %v", msgFlagMessage, err)
		}

		if message == "" {
			return fmt.Errorf("message is too short")
		}

		err = app.SendMessage(cmd.Context(), chatID, message)
		if err != nil {
			return fmt.Errorf("failed to send message: %v", err)
		}

		return nil
	},
}

func init() {
	messageCmd.Flags().Int64P(msgFlagChatID, msgFlagChatIDShort, 0, "Id of a chat")
	err := messageCmd.MarkFlagRequired(msgFlagChatID)
	if err != nil {
		os.Exit(1)
	}
	messageCmd.Flags().StringP(msgFlagMessage, msgFlagMessageShort, "", "Text of a message")
	err = messageCmd.MarkFlagRequired(msgFlagMessage)
	if err != nil {
		os.Exit(1)
	}
}
