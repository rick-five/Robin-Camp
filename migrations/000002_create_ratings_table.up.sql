CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    movie_title VARCHAR(255) NOT NULL,
    rater_id VARCHAR(255) NOT NULL,
    rating DECIMAL(2,1) NOT NULL CHECK (rating >= 0.5 AND rating <= 5.0 AND rating * 2 = ROUND(rating * 2)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(movie_title, rater_id),
    FOREIGN KEY (movie_title) REFERENCES movies(title) ON DELETE CASCADE
);

CREATE INDEX idx_ratings_movie_title ON ratings(movie_title);
CREATE INDEX idx_ratings_rater_id ON ratings(rater_id);