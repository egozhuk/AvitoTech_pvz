-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       email TEXT NOT NULL UNIQUE,
                       password_hash TEXT NOT NULL,
                       role TEXT NOT NULL CHECK (role IN ('employee', 'moderator')),
                       created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE pvz (
                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                     city TEXT NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань')),
                     registration_date TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE reception (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           pvz_id UUID NOT NULL REFERENCES pvz(id),
                           date_time TIMESTAMP NOT NULL DEFAULT now(),
                           status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
);

CREATE TABLE product (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
                         reception_id UUID NOT NULL REFERENCES reception(id) ON DELETE CASCADE,
                         date_time TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS product;
DROP TABLE IF EXISTS reception;
DROP TABLE IF EXISTS pvz;
DROP TABLE IF EXISTS users;
DROP TABLE users CASCADE;