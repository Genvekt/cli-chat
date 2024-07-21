package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/Genvekt/cli-chat/services/auth/model"
	"github.com/Genvekt/cli-chat/services/auth/repository"
)

const (
	// user is reserved word in postgres, needs double quotes in query
	userTable       = "\"user\""
	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

var _ repository.UserRepository = (*UserRepositoryPostgres)(nil)

// UserRepositoryPostgres implements repository.UserRepository for Postgres data source
type UserRepositoryPostgres struct {
	db *pgxpool.Pool
}

// NewUserRepositoryPostgres creates UserRepositoryPostgres instance
func NewUserRepositoryPostgres(db *pgxpool.Pool) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		db: db,
	}
}

// Create adds new user to db and updates id in user on success
func (r *UserRepositoryPostgres) Create(ctx context.Context, user *model.User) error {
	// Create and update times are set to now() as default in database
	builderInsert := sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, roleColumn).
		Values(user.Name, user.Email, user.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %v", err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}

// Get retrieves user by id
func (r *UserRepositoryPostgres) Get(ctx context.Context, id int64) (*model.User, error) {
	builderSelectOne := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {

		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	user := &model.User{}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user by id %d: %v", id, err)
	}

	return user, nil
}

// Update updates user in db
func (r *UserRepositoryPostgres) Update(
	ctx context.Context,
	id int64,
	updateFunc func(user *model.User) error,
) error {
	oldUser, err := r.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot update user with id %d: %v", id, err)
	}

	// Leave old user untouched to detect updates
	updatedUser := *oldUser

	err = updateFunc(&updatedUser)
	if err != nil {
		return fmt.Errorf("cannot update user with id %d: %v", id, err)
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

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user with id %d: %v", id, err)
	}

	return nil
}

// Delete removes user by id
func (r *UserRepositoryPostgres) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(userTable).PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	if err != nil {

		return fmt.Errorf("failed to build query: %v", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user with id %d: %v", id, err)
	}

	return nil
}
