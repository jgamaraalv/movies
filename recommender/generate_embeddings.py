import numpy as np
from config import EMBEDDING_DIM, TOP_RECOMMENDATIONS, SIMILAR_USERS_LIMIT


def extract_and_store_embeddings(model, num_users, num_movies, idx_to_user, idx_to_movie, conn):
    """Extract embeddings from trained model and store in pgvector."""
    print("Extracting embeddings from model...")

    # Get embedding weights
    user_gmf_weights = model.get_layer("user_gmf_embedding").get_weights()[0]
    user_mlp_weights = model.get_layer("user_mlp_embedding").get_weights()[0]
    movie_gmf_weights = model.get_layer("movie_gmf_embedding").get_weights()[0]
    movie_mlp_weights = model.get_layer("movie_mlp_embedding").get_weights()[0]

    cur = conn.cursor()

    # Store movie embeddings (average of GMF + MLP, L2-normalized)
    print("Storing movie embeddings...")
    movie_count = 0
    for idx in range(1, num_movies + 1):
        if idx not in idx_to_movie:
            continue
        movie_id = idx_to_movie[idx]
        emb = (movie_gmf_weights[idx] + movie_mlp_weights[idx]) / 2.0
        norm = np.linalg.norm(emb)
        if norm > 0:
            emb = emb / norm
        emb_list = emb.tolist()
        cur.execute(
            """
            INSERT INTO movie_embeddings (movie_id, embedding, updated_at)
            VALUES (%s, %s::vector, CURRENT_TIMESTAMP)
            ON CONFLICT (movie_id)
            DO UPDATE SET embedding = EXCLUDED.embedding, updated_at = CURRENT_TIMESTAMP
            """,
            (movie_id, str(emb_list)),
        )
        movie_count += 1

    # Store user embeddings (average of GMF + MLP, L2-normalized)
    print("Storing user embeddings...")
    user_count = 0
    for idx in range(1, num_users + 1):
        if idx not in idx_to_user:
            continue
        user_id = idx_to_user[idx]
        emb = (user_gmf_weights[idx] + user_mlp_weights[idx]) / 2.0
        norm = np.linalg.norm(emb)
        if norm > 0:
            emb = emb / norm
        emb_list = emb.tolist()
        cur.execute(
            """
            INSERT INTO user_embeddings (user_id, embedding, updated_at)
            VALUES (%s, %s::vector, CURRENT_TIMESTAMP)
            ON CONFLICT (user_id)
            DO UPDATE SET embedding = EXCLUDED.embedding, updated_at = CURRENT_TIMESTAMP
            """,
            (user_id, str(emb_list)),
        )
        user_count += 1

    conn.commit()
    cur.close()
    print(f"Stored {movie_count} movie embeddings and {user_count} user embeddings.")


def compute_recommendations_for_user(conn, user_id):
    """Compute recommendations for a single user using collaborative filtering."""
    cur = conn.cursor()

    # Get user's embedding
    cur.execute("SELECT embedding FROM user_embeddings WHERE user_id = %s", (user_id,))
    row = cur.fetchone()
    if not row:
        cur.close()
        return

    user_embedding = row[0]

    # Find similar users via cosine distance
    cur.execute(
        """
        SELECT user_id
        FROM user_embeddings
        WHERE user_id != %s
        ORDER BY embedding <=> %s::vector
        LIMIT %s
        """,
        (user_id, str(user_embedding.tolist()), SIMILAR_USERS_LIMIT),
    )
    similar_users = [r[0] for r in cur.fetchall()]

    if not similar_users:
        cur.close()
        return

    # Get current user's movies to exclude
    cur.execute("SELECT movie_id FROM user_movies WHERE user_id = %s", (user_id,))
    user_movie_ids = set(r[0] for r in cur.fetchall())

    # Collect movies from similar users
    placeholders = ",".join(["%s"] * len(similar_users))
    cur.execute(
        f"""
        SELECT um.movie_id, um.relation_type, m.popularity, m.score
        FROM user_movies um
        JOIN movies m ON m.id = um.movie_id
        WHERE um.user_id IN ({placeholders})
        """,
        similar_users,
    )

    movie_scores = {}
    for movie_id, relation_type, popularity, score in cur.fetchall():
        if movie_id in user_movie_ids:
            continue
        if movie_id not in movie_scores:
            movie_scores[movie_id] = {
                "weighted_count": 0.0,
                "scores": [],
                "popularity": popularity or 0.0,
            }
        weight = 1.0 if relation_type == "favorite" else 0.5
        movie_scores[movie_id]["weighted_count"] += weight
        if score is not None:
            movie_scores[movie_id]["scores"].append(score)

    # Score and rank
    scored_movies = []
    for movie_id, data in movie_scores.items():
        avg_score = sum(data["scores"]) / len(data["scores"]) if data["scores"] else 0.0
        pop = min(data["popularity"], 100.0) / 100.0
        final_score = data["weighted_count"] * 0.6 + (avg_score / 10.0) * 0.3 + pop * 0.1
        scored_movies.append((movie_id, final_score))

    scored_movies.sort(key=lambda x: x[1], reverse=True)
    top = scored_movies[:TOP_RECOMMENDATIONS]

    # Store recommendations
    cur.execute("DELETE FROM user_recommendations WHERE user_id = %s", (user_id,))
    for movie_id, score in top:
        cur.execute(
            """
            INSERT INTO user_recommendations (user_id, movie_id, score, computed_at)
            VALUES (%s, %s, %s, CURRENT_TIMESTAMP)
            ON CONFLICT (user_id, movie_id)
            DO UPDATE SET score = EXCLUDED.score, computed_at = CURRENT_TIMESTAMP
            """,
            (user_id, movie_id, float(score)),
        )

    conn.commit()
    cur.close()


def compute_all_recommendations(conn):
    """Compute recommendations for all users with embeddings."""
    cur = conn.cursor()
    cur.execute("SELECT user_id FROM user_embeddings")
    user_ids = [r[0] for r in cur.fetchall()]
    cur.close()

    print(f"Computing recommendations for {len(user_ids)} users...")
    for i, user_id in enumerate(user_ids):
        compute_recommendations_for_user(conn, user_id)
        if (i + 1) % 10 == 0:
            print(f"  Processed {i + 1}/{len(user_ids)} users")

    print("Done computing recommendations.")
