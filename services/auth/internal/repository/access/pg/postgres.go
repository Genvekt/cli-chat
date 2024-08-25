package pg

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/Genvekt/cli-chat/libraries/db_client/pkg/db"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
	repoConverter "github.com/Genvekt/cli-chat/services/auth/internal/repository/access/pg/converter"
	repoModel "github.com/Genvekt/cli-chat/services/auth/internal/repository/access/pg/model"
)

const (
	roleAccessRuleTable = "role_access_rule" // user is a reserved word in postgres, needs double quotes in query
	roleColumn          = "role"
	endpointColumn      = "endpoint"
)

var _ repository.AccessRepository = (*accessRepositoryPostgres)(nil)

// accessRepositoryPostgres implements repository.AccessRepository for Postgres data source
type accessRepositoryPostgres struct {
	db db.Client
}

// NewAccessRepositoryPostgres provides access repository for postgres
func NewAccessRepositoryPostgres(db db.Client) *accessRepositoryPostgres {
	return &accessRepositoryPostgres{
		db: db,
	}
}

// GetEndpointAccessRule retrieves all access rules for endpoint
func (r *accessRepositoryPostgres) GetEndpointAccessRule(
	ctx context.Context,
	endpoint string,
) (*model.EndpointAccessRule, error) {
	builderSelectOne := sq.Select(roleColumn, endpointColumn).
		From(roleAccessRuleTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{endpointColumn: endpoint})

	query, args, err := builderSelectOne.ToSql()
	if err != nil {

		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "access_repository.GetRoleAccessRule",
		QueryRaw: query,
	}

	var rules []*repoModel.RoleAccessRule

	err = r.db.DB().ScanAllContext(ctx, &rules, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrRuleNotFound
		}
		return nil, err
	}

	if len(rules) == 0 {
		return nil, repository.ErrRuleNotFound
	}

	return repoConverter.ToEndpointAccessRule(endpoint, rules), nil
}
