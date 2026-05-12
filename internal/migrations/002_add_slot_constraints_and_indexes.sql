-- +goose Up

-- Включаем расширение для работы с диапазонами и GiST-индексами
CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE slots
    ADD CONSTRAINT no_overlapping_slots
        EXCLUDE USING gist (
        tutor_id WITH =,
        tstzrange(start_time, end_time) WITH &&
        );

CREATE INDEX idx_slots_tutor_free ON slots (tutor_id, start_time) WHERE status = 'free';

CREATE INDEX idx_slots_student_cart ON slots (student_id) WHERE status = 'in_cart';

CREATE INDEX idx_tutor_subjects_subject_id ON tutor_subjects (subject_id);


-- +goose Down

DROP INDEX IF EXISTS idx_tutor_subjects_subject_id;
DROP INDEX IF EXISTS idx_slots_student_cart;
DROP INDEX IF EXISTS idx_slots_tutor_free;

ALTER TABLE slots DROP CONSTRAINT IF EXISTS no_overlapping_slots;
