package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/Genvekt/cli-chat/services/auth/internal/client/db"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
	repoConverter "github.com/Genvekt/cli-chat/services/auth/internal/repository/user/converter"
	repoModel "github.com/Genvekt/cli-chat/services/auth/internal/repository/user/model"
)

const (
	userTable       = "\"user\"" // user is a reserved word in postgres, needs double quotes in query
	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

var _ repository.UserRepository = (*userRepositoryPostgres)(nil)

// userRepositoryPostgres implements repository.UserRepository for Postgres data source
type userRepositoryPostgres struct {
	db db.Client
}

// NewUserRepositoryPostgres creates UserRepositoryPostgres instance
func NewUserRepositoryPostgres(db db.Client) *userRepositoryPostgres {
	return &userRepositoryPostgres{
		db: db,
	}
}

// Create adds new user to db and updates id in user on success
func (r *userRepositoryPostgres) Create(ctx context.Context, user *model.User) (int64, error) {
	// Create and update times are set to now() as default in database
	builderInsert := sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, roleColumn).
		Values(user.Name, user.Email, user.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var newUserID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&newUserID)
	if err != nil {
		return 0, err
	}

	return newUserID, nil
}

// Get retrieves user by id
func (r *userRepositoryPostgres) Get(ctx context.Context, id int64) (*model.User, error) {
	builderSelectOne := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {

		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	dbUser := &repoModel.User{}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	err = r.db.DB().ScanOneContext(ctx, dbUser, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return repoConverter.ToUserFromRepo(dbUser), nil
}

// Update updates user in db
func (r *userRepositoryPostgres) Update(
	ctx context.Context,
	id int64,
	updateFunc func(user *model.User) error,
) error {
	oldUser, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	// Leave old user untouched to detect updates
	updatedUser := *oldUser

	err = updateFunc(&updatedUser)
	if err != nil {
		return err
	}

	if updatedUser.ID != id {
		return fmt.Errorf("id must not change after update: changed from %d to %d", id, updatedUser.ID)
	}

	builderUpdate := sq.Update(userTable).PlaceholderFormat(sq.Dollar)

	if oldUser.Name != updatedUser.Name {
		builderUpdate = builderUpdate.Set(nameColumn, updatedUser.Name)
	}

	if oldUser.Email != updatedUser.Email {
		builderUpdate = builderUpdate.Set(emailColumn, updatedUser.Email)
	}

	if oldUser.Role != updatedUser.Role {
		builderUpdate = builderUpdate.Set(roleColumn, updatedUser.Role)
	}

	builderUpdate = builderUpdate.Set(updatedAtColumn, time.Now()).Where(sq.Eq{idColumn: oldUser.ID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	if len(args) <= 2 {
		// nothing to update, there are only updatedAt and ID filter filled
		return nil
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes user by id
func (r *userRepositoryPostgres) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(userTable).PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {

		return fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

// GetList retrieves users by their names
func (r *userRepositoryPostgres) GetList(ctx context.Context, names []string) ([]*model.User, error) {
	builderQuery := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(userTable).
		Where(sq.Eq{nameColumn: names})

	query, args, err := builderQuery.ToSql()
	if err != nil {

		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "user_repository.Query",
		QueryRaw: query,
	}

	var usersDB []*repoModel.User
	err = r.db.DB().ScanAllContext(ctx, &usersDB, q, args...)
	if err != nil {
		return nil, err
	}

	return repoConverter.ToUsersFromRepo(usersDB), nil
}
