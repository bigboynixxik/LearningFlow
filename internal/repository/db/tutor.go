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

type TutorRepo struct {
	pool *pgxpool.Pool
	sq   sq.StatementBuilderType
}

func NewTutorRepo(pool *pgxpool.Pool) *TutorRepo {
	return &TutorRepo{
		pool: pool,
		sq:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *TutorRepo) GetAll(ctx context.Context) ([]models.Tutor, error) {
	query, args, err := r.sq.Select(
		"t.user_id", "t.name", "t.hourly_rate",
		"COALESCE(t.description, '')", "COALESCE(t.photo_path, '')",
		"COALESCE(array_agg(ts.subject_id) FILTER (WHERE ts.subject_id IS NOT NULL), '{}')",
	).
		From("tutors t").
		LeftJoin("tutor_subjects ts ON t.user_id = ts.tutor_id").
		GroupBy("t.user_id", "t.name", "t.hourly_rate", "t.description", "t.photo_path").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("db.tutor.GetAll failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db.tutor.GetAll failed to execute query: %w", err)
	}
	defer rows.Close()

	var tutors []models.Tutor
	for rows.Next() {
		var t models.Tutor
		if err := rows.Scan(&t.UserID, &t.Name, &t.HourlyRate, &t.Description, &t.PhotoPath, &t.SubjectIDs); err != nil {
			return nil, fmt.Errorf("db.tutor.GetAll failed to scan row: %w", err)
		}
		tutors = append(tutors, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db.tutor.GetAll rows.Err: %w", err)
	}
	return tutors, nil
}

func (r *TutorRepo) GetByID(ctx context.Context, id string) (*models.Tutor, error) {
	query, args, err := r.sq.Select(
		"t.user_id", "t.name", "t.hourly_rate",
		"COALESCE(t.description, '')", "COALESCE(t.photo_path, '')",
		"COALESCE(array_agg(ts.subject_id) FILTER (WHERE ts.subject_id IS NOT NULL), '{}')",
	).
		From("tutors t").
		LeftJoin("tutor_subjects ts ON t.user_id = ts.tutor_id").
		Where(sq.Eq{"t.user_id": id}).
		GroupBy("t.user_id", "t.name", "t.hourly_rate", "t.description", "t.photo_path").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("db.tutor.GetByID failed to build query: %w", err)
	}

	var t models.Tutor
	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&t.UserID, &t.Name, &t.HourlyRate, &t.Description, &t.PhotoPath, &t.SubjectIDs); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("db.tutor.GetByID failed to scan: %w", err)
	}
	return &t, nil
}

func (r *TutorRepo) GetBySubjectID(ctx context.Context, subjectID int64) ([]models.Tutor, error) {
	query, args, err := r.sq.Select(
		"t.user_id", "t.name", "t.hourly_rate",
		"COALESCE(t.description, '')", "COALESCE(t.photo_path, '')",
		"COALESCE(array_agg(ts.subject_id) FILTER (WHERE ts.subject_id IS NOT NULL), '{}')",
	).
		From("tutors t").
		LeftJoin("tutor_subjects ts ON t.user_id = ts.tutor_id").
		// Обернули в sq.Expr для безопасного маппинга плейсхолдера
		Where(sq.Expr("t.user_id IN (SELECT tutor_id FROM tutor_subjects WHERE subject_id = ?)", subjectID)).
		GroupBy("t.user_id", "t.name", "t.hourly_rate", "t.description", "t.photo_path").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("db.tutor.GetBySubjectID failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db.tutor.GetBySubjectID failed to execute: %w", err)
	}
	defer rows.Close()

	var tutors []models.Tutor
	for rows.Next() {
		var t models.Tutor
		if err := rows.Scan(&t.UserID, &t.Name, &t.HourlyRate, &t.Description, &t.PhotoPath, &t.SubjectIDs); err != nil {
			return nil, fmt.Errorf("db.tutor.GetBySubjectID failed to scan row: %w", err)
		}
		tutors = append(tutors, t)
	}
	return tutors, nil
}
