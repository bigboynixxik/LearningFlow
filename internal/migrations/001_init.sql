-- +goose Up

-- Таблица пользователей (Базовая для авторизации)
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       role VARCHAR(50) NOT NULL, -- 'student' или 'tutor'
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
                          token VARCHAR(255) PRIMARY KEY, -- Сам токен, который мы отдадим в Cookie
                          user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          expires_at TIMESTAMP WITH TIME ZONE NOT NULL, -- Время "протухания" сессии
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица предметов (Справочник)
CREATE TABLE subjects (
                          id BIGSERIAL PRIMARY KEY,
                          name VARCHAR(255) UNIQUE NOT NULL
);

-- Профиль репетитора (Связан с users)
CREATE TABLE tutors (
                        user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                        name VARCHAR(255) NOT NULL,
                        hourly_rate INT NOT NULL DEFAULT 0, -- В рублях
                        description TEXT,
                        photo_path VARCHAR(255)
);

-- Связь Многие-ко-Многим (Репетитор <-> Предмет)
CREATE TABLE tutor_subjects (
                                tutor_id UUID REFERENCES tutors(user_id) ON DELETE CASCADE,
                                subject_id BIGINT REFERENCES subjects(id) ON DELETE CASCADE,
                                PRIMARY KEY (tutor_id, subject_id)
);

-- Слоты времени (Расписание и Корзина в одном флаконе)
CREATE TABLE slots (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       tutor_id UUID NOT NULL REFERENCES tutors(user_id) ON DELETE CASCADE,
                       student_id UUID REFERENCES users(id) ON DELETE SET NULL, -- ID ученика (если добавлено в корзину/куплено)
                       start_time TIMESTAMP WITH TIME ZONE NOT NULL,
                       end_time TIMESTAMP WITH TIME ZONE NOT NULL,
                       status VARCHAR(50) NOT NULL DEFAULT 'free', -- 'free', 'in_cart', 'booked'
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Наполняем справочник предметов
INSERT INTO subjects (name) VALUES ('Математика'), ('Информатика'), ('Английский язык'), ('Физика');

-- +goose Down
DROP TABLE slots;
DROP TABLE tutor_subjects;
DROP TABLE tutors;
DROP TABLE subjects;
DROP TABLE users;