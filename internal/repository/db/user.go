package db

import (
	"context"
	"fmt"
	"learningflow/internal/models"
	"learningflow/pkg/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
	sq   sq.StatementBuilderType
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
		sq:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	l := logger.FromContext(ctx)

	user.ID = uuid.NewString()

	query, args, err := r.sq.Insert("users").
		Columns("id", "email", "password_hash", "role").
		Values(user.ID, user.Email, user.PasswordHash, user.Role).
		ToSql()
	if err != nil {
		l.Error("db.Create failed to build query", "query", query, "args", args, "err", err)
		return fmt.Errorf("UserRepository.Create build query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		l.Error("user.Create failed to execute", "query", query, "args", args, "err", err)
		return fmt.Errorf("UserRepository.Create exec: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	l := logger.FromContext(ctx)

	query, args, err := r.sq.Select("id", "email", "password_hash", "role").
		From("users").
		Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		l.Error("user.GetByEmail failed to build query", "query", query, "args", args, "err", err)
		return nil, fmt.Errorf("user.GetByEmail build query: %w", err)
	}
	row := r.pool.QueryRow(ctx, query, args...)
	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user.GetByEmail no rows found")
		}
		return nil, fmt.Errorf("user.GetByEmail scan: %w", err)
	}
	return &user, nil
}
