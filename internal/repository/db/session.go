package db

import (
	"context"
	"fmt"
	"learningflow/internal/models"
	"learningflow/pkg/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	pool *pgxpool.Pool
	sq   sq.StatementBuilderType
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		pool: pool,
		sq:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	l := logger.FromContext(ctx)
	query, args, err := r.sq.Insert("sessions").
		Columns("token", "user_id", "expires_at").
		Values(session.Token, session.UserID, session.ExpiresAt).ToSql()
	if err != nil {
		l.Error("session.CreateSession failed to build query", "query", query, "args", args, "err", err)
		return fmt.Errorf("session.CreateSession: %w", err)
	}
	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("session.CreateSession: %w", err)
	}
	return nil
}

func (r *SessionRepository) GetByToken(ctx context.Context, token string) (*models.Session, error) {
	l := logger.FromContext(ctx)
	query, args, err := r.sq.Select("token", "user_id", "expires_at").
		From("sessions").
		Where(sq.Eq{"token": token}).ToSql()
	if err != nil {
		l.Error("session.GetByToken failed to build query", "token", token, "err", err)
		return nil, fmt.Errorf("session.GetByToken: %w", err)
	}
	row := r.pool.QueryRow(ctx, query, args...)
	session := &models.Session{}
	if err := row.Scan(&session.Token, &session.UserID, &session.ExpiresAt); err != nil {
		l.Error("session.GetByToken failed to scan row", "token", token, "err", err)
		return nil, fmt.Errorf("session.GetByToken: %w", err)
	}
	return session, nil
}

func (r *SessionRepository) Delete(ctx context.Context, token string) error {
	l := logger.FromContext(ctx)
	query, args, err := r.sq.Delete("sessions").Where(sq.Eq{"token": token}).ToSql()
	if err != nil {
		l.Error("session.Delete failed to build query", "token", token, "err", err)
		return fmt.Errorf("session.Delete: %w", err)
	}
	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("session.Delete: %w", err)
	}
	return nil
}
