package db

import (
	"context"
	"errors"
	"fmt"
	"learningflow/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubjectRepo struct {
	pool *pgxpool.Pool
	sq   sq.StatementBuilderType
}

func NewSubjectRepo(pool *pgxpool.Pool) *SubjectRepo {
	return &SubjectRepo{
		pool: pool,
		sq:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *SubjectRepo) GetAll(ctx context.Context) ([]models.Subject, error) {
	query, args, err := r.sq.Select("id", "name").
		From("subjects").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("db.subject.GetAll get all subjects: %w", err)
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db.subject.GetAll get all subjects: %w", err)
	}
	defer rows.Close()
	var subjects []models.Subject
	for rows.Next() {
		var subject models.Subject
		err := rows.Scan(&subject.ID, &subject.Name)
		if err != nil {
			return nil, fmt.Errorf("db.subject.GetAll get all subjects: %w", err)
		}
		subjects = append(subjects, subject)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("db.subject.GetAll get all subjects: %w", err)
	}
	return subjects, nil
}

func (r *SubjectRepo) GetByID(ctx context.Context, id int64) (*models.Subject, error) {
	query, args, err := r.sq.Select("id", "name").
		From("subjects").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("db.subject.GetByID build query: %w", err)
	}

	row := r.pool.QueryRow(ctx, query, args...)
	var subject models.Subject

	if err := row.Scan(&subject.ID, &subject.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("db.subject.GetByID scan: %w", err)
	}
	return &subject, nil
}
