-- Создание расширений
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание таблицы администраторов
CREATE TABLE IF NOT EXISTS admins (
                                      id SERIAL PRIMARY KEY,
                                      name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    admin_status INTEGER NOT NULL CHECK (admin_status IN (-1, 0, 1)),
    created_chapters JSONB NOT NULL DEFAULT '[]'::jsonb,
    request_sent JSONB NOT NULL DEFAULT '[]'::jsonb,
    requests_received JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание таблицы глав
CREATE TABLE IF NOT EXISTS chapters (
                                        id SERIAL PRIMARY KEY,
                                        name VARCHAR(255) NOT NULL,
    start_node INTEGER NOT NULL,
    nodes JSONB NOT NULL DEFAULT '[]'::jsonb,
    characters JSONB NOT NULL DEFAULT '[]'::jsonb,
    status INTEGER NOT NULL CHECK (status IN (1, 2, 3)),
    updated_at JSONB NOT NULL DEFAULT '{}'::jsonb,
    author INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author) REFERENCES admins(id)
    );

-- Создание таблицы персонажей
CREATE TABLE IF NOT EXISTS characters (
                                          id SERIAL PRIMARY KEY,
                                          name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    color VARCHAR(255) NOT NULL,
    emotions JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание таблицы медиа
CREATE TABLE IF NOT EXISTS media (
                                     id SERIAL PRIMARY KEY,
                                     file_data BYTEA NOT NULL,
                                     content_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание таблицы узлов
CREATE TABLE IF NOT EXISTS nodes (
                                     id SERIAL PRIMARY KEY,
                                     slug VARCHAR(255) UNIQUE NOT NULL,
    chapter_id INTEGER NOT NULL,
    music_id INTEGER,
    background_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chapter_id) REFERENCES chapters(id),
    FOREIGN KEY (music_id) REFERENCES media(id),
    FOREIGN KEY (background_id) REFERENCES media(id)
    );

-- Создание таблицы запросов
CREATE TABLE IF NOT EXISTS requests (
                                        id SERIAL PRIMARY KEY,
                                        type INTEGER NOT NULL CHECK (type IN (0, 1, 2)),
    status INTEGER NOT NULL CHECK (status IN (0, 1)),
    requesting_admin INTEGER NOT NULL,
    requested_chapter_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (requesting_admin) REFERENCES admins(id),
    FOREIGN KEY (requested_chapter_id) REFERENCES chapters(id)
    );

-- Создание таблицы игроков
CREATE TABLE IF NOT EXISTS players (
                                       id SERIAL PRIMARY KEY,
                                       name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(255),
    password VARCHAR(255) NOT NULL,
    admin BOOLEAN NOT NULL DEFAULT FALSE,
    completed_chapters JSONB NOT NULL DEFAULT '[]'::jsonb,
    chapters_progress JSONB NOT NULL DEFAULT '{}'::jsonb,
    sound_settings INTEGER NOT NULL DEFAULT 100,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание индексов для оптимизации поиска
CREATE INDEX IF NOT EXISTS idx_chapters_author ON chapters(author);
CREATE INDEX IF NOT EXISTS idx_nodes_chapter ON nodes(chapter_id);
CREATE INDEX IF NOT EXISTS idx_requests_admin ON requests(requesting_admin);
CREATE INDEX IF NOT EXISTS idx_requests_chapter ON requests(requested_chapter_id);
CREATE INDEX IF NOT EXISTS idx_players_email ON players(email) USING GIST;