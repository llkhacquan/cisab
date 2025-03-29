-- Enable the citext extension
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE user_role AS ENUM ('employer', 'employee');

CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    email         CITEXT UNIQUE UNIQUE NOT NULL,
    password_hash VARCHAR(255)         NOT NULL,
    name          VARCHAR(200)         NOT NULL,
    role          user_role            NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
