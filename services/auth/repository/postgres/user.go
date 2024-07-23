package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/Genvekt/cli-chat/services/auth/model"
)

const (
	userTable       = "\"user\""
	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type UserRepositoryPostgres struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewUserRepositoryPostgres(ctx context.Context, db *pgxpool.Pool) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		ctx: ctx,
		db:  db,
	}
}

func (r *UserRepositoryPostgres) Get(id int64) (*model.User, error) {
	user := &model.User{}

	// Building request for user by provided id
	builderSelectOne := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(userTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {

		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	err = r.db.QueryRow(r.ctx, query, args...).Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id %d: %v", id, err)
	}
	return user, nil
}

func (r *UserRepositoryPostgres) Create(user *model.User) (*model.User, error) {
	creationTime := time.Now().In(time.UTC)
	builderInsert := sq.Insert(userTable).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		Values(user.Name, user.Email, user.Role, creationTime, creationTime).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}
	err = r.db.QueryRow(r.ctx, query, args...).Scan(&user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}
	return user, nil
}
func (r *UserRepositoryPostgres) Update(user *model.User) error {
	updateTime := time.Now().In(time.UTC)
	builderUpdate := sq.Update(userTable).PlaceholderFormat(sq.Dollar).
		Set(nameColumn, user.Name).
		Set(emailColumn, user.Email).
		Set(roleColumn, user.Role).
		Set(updatedAtColumn, updateTime).
		Where(sq.Eq{idColumn: user.Id})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	_, err = r.db.Exec(r.ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update user with id %d: %v", user.Id, err)
	}
	return nil
}

func (r *UserRepositoryPostgres) Delete(id int64) error {
	builderDelete := sq.Delete(userTable).PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builderDelete.ToSql()
	_, err = r.db.Exec(r.ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to delete user with id %d: %v", id, err)
	}
	return nil
}
