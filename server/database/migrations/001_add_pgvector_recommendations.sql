CREATE EXTENSION IF NOT EXISTS vector;

-- 128-dim movie embeddings from TensorFlow NCF model
CREATE TABLE movie_embeddings (
    movie_id   int4 PRIMARY KEY REFERENCES movies(id) ON DELETE CASCADE,
    embedding  vector(128) NOT NULL,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 128-dim user preference vectors (aggregated from movie embeddings)
CREATE TABLE user_embeddings (
    user_id    int4 PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    embedding  vector(128) NOT NULL,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Precomputed recommendation cache
CREATE TABLE user_recommendations (
    user_id    int4 NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id   int4 NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    score      float4 NOT NULL,
    reason     text,
    computed_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, movie_id)
);

-- HNSW indexes for fast approximate nearest neighbor (cosine distance)
CREATE INDEX idx_user_embeddings_hnsw ON user_embeddings
    USING hnsw (embedding vector_cosine_ops) WITH (m = 16, ef_construction = 64);

CREATE INDEX idx_movie_embeddings_hnsw ON movie_embeddings
    USING hnsw (embedding vector_cosine_ops) WITH (m = 16, ef_construction = 64);

CREATE INDEX idx_user_recommendations_score ON user_recommendations (user_id, score DESC);
