CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE movies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) UNIQUE NOT NULL,
    release_date DATE NOT NULL,
    genre VARCHAR(100) NOT NULL,
    distributor VARCHAR(255),
    budget BIGINT,
    mpa_rating VARCHAR(10),
    box_office JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_movies_title ON movies(title);
CREATE INDEX idx_movies_genre ON movies(genre);
CREATE INDEX idx_movies_release_date ON movies(release_date);