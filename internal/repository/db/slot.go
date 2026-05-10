package db

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SlotRepo struct {
	pool *pgxpool.Pool
	sq   sq.StatementBuilderType
}

func NewSlotRepo(pool *pgxpool.Pool) *SlotRepo {
	return &SlotRepo{
		pool: pool,
		sq:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// AddToCart создает тестовый слот для репетитора и сразу кладет его в корзину ученика
func (r *SlotRepo) AddToCart(ctx context.Context, tutorID string, studentID string) error {

	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(1 * time.Hour)

	query, args, err := r.sq.Insert("slots").
		Columns("tutor_id", "student_id", "start_time", "end_time", "status").
		Values(tutorID, studentID, startTime, endTime, "in_cart").
		ToSql()

	if err != nil {
		return fmt.Errorf("db.slot.AddToCart build query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("db.slot.AddToCart execute: %w", err)
	}

	return nil
}

// GetCartCount возвращает количество товаров в корзине для счетчика в шапке
func (r *SlotRepo) GetCartCount(ctx context.Context, studentID string) (int, error) {
	query, args, err := r.sq.Select("COUNT(*)").
		From("slots").
		Where(sq.Eq{"student_id": studentID, "status": "in_cart"}).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("db.slot.GetCartCount build query: %w", err)
	}

	var count int
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("db.slot.GetCartCount scan: %w", err)
	}

	return count, nil
}
