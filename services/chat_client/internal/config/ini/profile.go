package ini

import (
	"fmt"

	"gopkg.in/ini.v1"

	"github.com/Genvekt/cli-chat/services/chat-client/internal/config"
)

const (
	defaultProfile = "default"
	profilePrefix  = "profile "
	usernameKey    = "username"
	passwordKey    = "password"
)

var _ config.ProfileConfig = (*profileConfigIni)(nil)

type profileConfigIni struct {
	username string
	password string
}

func NewProfileConfigIni(configPath string, profileName string) (*profileConfigIni, error) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %v", err)
	}

	profile, err := cfg.GetSection(fmt.Sprintf("%s%s", profilePrefix, defaultProfile))
	if profileName == "" && err != nil {
		return nil, fmt.Errorf("failed to get default profile: %v", err)
	}

	if profileName != "" {
		profile, err = cfg.GetSection(fmt.Sprintf("%s%s", profilePrefix, profileName))
		if err != nil {
			return nil, fmt.Errorf("failed to get profile %s: %v", profileName, err)
		}
	}

	username, err := profile.GetKey(usernameKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get username: %v", err)
	}

	password, err := profile.GetKey(passwordKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get password: %v", err)
	}

	return &profileConfigIni{
		username: username.String(),
		password: password.String(),
	}, nil
}

func (c *profileConfigIni) Username() string {
	return c.username
}

func (c *profileConfigIni) Password() string {
	return c.password
}
