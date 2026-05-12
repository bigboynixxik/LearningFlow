package db

import (
	"context"
	"errors"
	"fmt"
	"learningflow/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *SlotRepo) Create(ctx context.Context, slot *models.Slot) error {
	query, args, err := r.sq.Insert("slots").
		Columns("tutor_id", "start_time", "end_time", "status").
		Values(slot.TutorID, slot.StartTime, slot.EndTime, "free").
		ToSql()

	if err != nil {
		return fmt.Errorf("db.slot.Create build: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23P04" {
			return models.ErrSlotOverlap
		}
		return fmt.Errorf("db.slot.Create exec: %w", err)
	}
	return nil
}

func (r *SlotRepo) GetFreeSlots(ctx context.Context, tutorID string) ([]models.Slot, error) {
	query, args, err := r.sq.Select("id", "tutor_id", "start_time", "end_time", "status").
		From("slots").
		Where(sq.Eq{"tutor_id": tutorID, "status": "free"}).
		Where(sq.Gt{"start_time": time.Now()}).
		OrderBy("start_time ASC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("db.slot.GetFreeSlots build: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db.slot.GetFreeSlots query: %w", err)
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		var s models.Slot
		if err := rows.Scan(&s.ID, &s.TutorID, &s.StartTime, &s.EndTime, &s.Status); err != nil {
			return nil, fmt.Errorf("db.slot.GetFreeSlots scan: %w", err)
		}
		slots = append(slots, s)
	}
	return slots, nil
}

func (r *SlotRepo) AddToCart(ctx context.Context, slotID string, studentID string) error {
	query, args, err := r.sq.Update("slots").
		Set("status", "in_cart").
		Set("student_id", studentID).
		Where(sq.Eq{"id": slotID, "status": "free"}).
		ToSql()

	if err != nil {
		return fmt.Errorf("db.slot.AddToCart build: %w", err)
	}

	ct, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("db.slot.AddToCart exec: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return models.ErrSlotUnavailable
	}
	return nil
}

func (r *SlotRepo) GetCart(ctx context.Context, studentID string) ([]models.Slot, error) {
	query, args, err := r.sq.Select(
		"s.id", "s.tutor_id", "s.start_time", "s.end_time", "s.status",
		"t.name", "t.hourly_rate",
	).
		From("slots s").
		Join("tutors t ON s.tutor_id = t.user_id").
		Where(sq.Eq{"s.student_id": studentID, "s.status": "in_cart"}).
		OrderBy("s.start_time ASC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("db.slot.GetCart build: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db.slot.GetCart query: %w", err)
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		var s models.Slot
		if err := rows.Scan(&s.ID, &s.TutorID, &s.StartTime, &s.EndTime, &s.Status, &s.TutorName, &s.HourlyRate); err != nil {
			return nil, fmt.Errorf("db.slot.GetCart scan: %w", err)
		}
		slots = append(slots, s)
	}
	return slots, nil
}

func (r *SlotRepo) RemoveFromCart(ctx context.Context, slotID string, studentID string) error {
	query, args, err := r.sq.Update("slots").
		Set("status", "free").
		Set("student_id", nil).
		Where(sq.Eq{"id": slotID, "student_id": studentID, "status": "in_cart"}).
		ToSql()

	if err != nil {
		return fmt.Errorf("db.slot.RemoveFromCart build: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("db.slot.RemoveFromCart exec: %w", err)
	}
	return nil
}

func (r *SlotRepo) CheckoutCart(ctx context.Context, studentID string) error {
	query, args, err := r.sq.Update("slots").
		Set("status", "booked").
		Where(sq.Eq{"student_id": studentID, "status": "in_cart"}).
		ToSql()

	if err != nil {
		return fmt.Errorf("db.slot.CheckoutCart build: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("db.slot.CheckoutCart exec: %w", err)
	}
	return nil
}

func (r *SlotRepo) GetCartCount(ctx context.Context, studentID string) (int, error) {
	query, args, err := r.sq.Select("COUNT(*)").
		From("slots").
		Where(sq.Eq{"student_id": studentID, "status": "in_cart"}).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("db.slot.GetCartCount build: %w", err)
	}

	var count int
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("db.slot.GetCartCount scan: %w", err)
	}
	return count, nil
}
