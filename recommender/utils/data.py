import numpy as np
from utils.db import get_connection


def load_interactions():
    """Load user-movie interactions from the database."""
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("""
        SELECT user_id, movie_id, relation_type
        FROM user_movies
    """)
    rows = cur.fetchall()
    cur.close()
    conn.close()
    return rows


def load_movie_genres():
    """Load movie-genre associations."""
    conn = get_connection()
    cur = conn.cursor()
    cur.execute("""
        SELECT movie_id, genre_id
        FROM movie_genres
    """)
    rows = cur.fetchall()
    cur.close()
    conn.close()
    return rows


def build_mappings(interactions):
    """Build contiguous ID mappings for users and movies."""
    user_ids = sorted(set(r[0] for r in interactions))
    movie_ids = sorted(set(r[1] for r in interactions))

    user_to_idx = {uid: idx + 1 for idx, uid in enumerate(user_ids)}
    movie_to_idx = {mid: idx + 1 for idx, mid in enumerate(movie_ids)}
    idx_to_user = {idx: uid for uid, idx in user_to_idx.items()}
    idx_to_movie = {idx: mid for mid, idx in movie_to_idx.items()}

    return user_to_idx, movie_to_idx, idx_to_user, idx_to_movie


def prepare_training_data(interactions, user_to_idx, movie_to_idx, movie_genres, num_genres, negative_ratio=4):
    """Prepare training data with negative sampling."""
    all_movie_idxs = set(movie_to_idx.values())

    # Build genre lookup per movie index
    genre_lookup = {}
    max_genres = 5
    for mid, gid in movie_genres:
        if mid in movie_to_idx:
            idx = movie_to_idx[mid]
            if idx not in genre_lookup:
                genre_lookup[idx] = []
            if len(genre_lookup[idx]) < max_genres:
                genre_lookup[idx].append(gid)

    # Build positive samples
    user_items = {}
    users, movies, genres_arr, labels = [], [], [], []

    for user_id, movie_id, relation_type in interactions:
        if user_id not in user_to_idx or movie_id not in movie_to_idx:
            continue
        u_idx = user_to_idx[user_id]
        m_idx = movie_to_idx[movie_id]

        if u_idx not in user_items:
            user_items[u_idx] = set()
        user_items[u_idx].add(m_idx)

        label = 1.0 if relation_type == "favorite" else 0.7
        users.append(u_idx)
        movies.append(m_idx)
        g = genre_lookup.get(m_idx, [0])
        genres_arr.append(g + [0] * (max_genres - len(g)))
        labels.append(label)

    # Generate negative samples
    rng = np.random.default_rng(42)
    num_positives = len(users)
    for _ in range(num_positives * negative_ratio):
        u_idx = users[rng.integers(num_positives)]
        neg_movie = rng.integers(1, len(all_movie_idxs) + 1)
        while neg_movie in user_items.get(u_idx, set()):
            neg_movie = rng.integers(1, len(all_movie_idxs) + 1)

        users.append(u_idx)
        movies.append(neg_movie)
        g = genre_lookup.get(neg_movie, [0])
        genres_arr.append(g + [0] * (max_genres - len(g)))
        labels.append(0.0)

    return (
        np.array(users, dtype=np.int32),
        np.array(movies, dtype=np.int32),
        np.array(genres_arr, dtype=np.int32),
        np.array(labels, dtype=np.float32),
    )
