package create

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Genvekt/cli-chat/services/chat-client/internal/app"
)

const (
	chatFlagName      = "name"
	chatFlagNameShort = "n"

	chatFlagUsernames      = "usernames"
	chatFlagUsernamesShort = "u"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Create new chat",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := cmd.Flags().GetString(chatFlagName)
		if err != nil {
			return fmt.Errorf("failed to read flag %s: %v", chatFlagName, err)
		}

		usernames, err := cmd.Flags().GetStringSlice(chatFlagUsernames)
		if err != nil {

			return fmt.Errorf("failed to read flag %s: %v", chatFlagUsernames, err)
		}

		err = app.CreateChat(cmd.Context(), name, usernames)
		if err != nil {
			return fmt.Errorf("failed to create chat: %v", err)
		}

		return nil
	},
}

func init() {
	chatCmd.Flags().StringP(chatFlagName, chatFlagNameShort, "", "Name of a chat")
	err := chatCmd.MarkFlagRequired(chatFlagName)
	if err != nil {
		os.Exit(1)
	}

	chatCmd.Flags().StringSliceP(chatFlagUsernames, chatFlagUsernamesShort, []string{}, "Usernames of participants")
	err = chatCmd.MarkFlagRequired(chatFlagUsernames)
	if err != nil {
		os.Exit(1)
	}
}
