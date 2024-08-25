package converter

import (
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	repoModel "github.com/Genvekt/cli-chat/services/auth/internal/repository/access/pg/model"
)

// ToEndpointAccessRule converts repository rules to service level rules
func ToEndpointAccessRule(endpoint string, rules []*repoModel.RoleAccessRule) *model.EndpointAccessRule {
	if len(rules) == 0 {
		return &model.EndpointAccessRule{
			Endpoint: endpoint,
			Roles:    map[int]struct{}{},
		}
	}

	rule := &model.EndpointAccessRule{
		Endpoint: endpoint,
		Roles:    make(map[int]struct{}, len(rules)),
	}

	for _, roleRule := range rules {
		rule.Roles[roleRule.Role] = struct{}{}
	}

	return rule
}
