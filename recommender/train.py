import numpy as np
import tensorflow as tf
from sklearn.model_selection import train_test_split

from config import EPOCHS, BATCH_SIZE, LEARNING_RATE, EARLY_STOPPING_PATIENCE
from models.ncf import build_ncf_model
from utils.data import load_interactions, load_movie_genres, build_mappings, prepare_training_data
from utils.db import get_connection
from generate_embeddings import extract_and_store_embeddings, compute_all_recommendations


def main():
    print("=" * 60)
    print("Moovies Recommendation Training Pipeline")
    print("=" * 60)

    # Load data
    print("\n[1/5] Loading interactions from database...")
    interactions = load_interactions()
    if not interactions:
        print("No interactions found. Exiting.")
        return

    print(f"  Found {len(interactions)} interactions")

    print("\n[2/5] Preparing training data...")
    movie_genres = load_movie_genres()
    user_to_idx, movie_to_idx, idx_to_user, idx_to_movie = build_mappings(interactions)

    num_users = len(user_to_idx)
    num_movies = len(movie_to_idx)

    # Determine number of genres
    genre_ids = set(g[1] for g in movie_genres)
    num_genres = max(genre_ids) if genre_ids else 0

    print(f"  Users: {num_users}, Movies: {num_movies}, Genres: {num_genres}")

    users, movies, genres, labels = prepare_training_data(
        interactions, user_to_idx, movie_to_idx, movie_genres, num_genres
    )
    print(f"  Training samples: {len(labels)} (positive: {np.sum(labels > 0)}, negative: {np.sum(labels == 0)})")

    # Split data
    indices = np.arange(len(labels))
    train_idx, val_idx = train_test_split(indices, test_size=0.1, random_state=42)

    train_data = {
        "user_input": users[train_idx],
        "movie_input": movies[train_idx],
        "genre_input": genres[train_idx],
    }
    val_data = {
        "user_input": users[val_idx],
        "movie_input": movies[val_idx],
        "genre_input": genres[val_idx],
    }

    # Build and compile model
    print("\n[3/5] Building NCF model...")
    model = build_ncf_model(num_users, num_movies, num_genres)
    model.compile(
        optimizer=tf.keras.optimizers.Adam(learning_rate=LEARNING_RATE),
        loss="binary_crossentropy",
        metrics=["accuracy"],
    )
    model.summary()

    # Train
    print("\n[4/5] Training model...")
    callbacks = [
        tf.keras.callbacks.EarlyStopping(
            monitor="val_loss",
            patience=EARLY_STOPPING_PATIENCE,
            restore_best_weights=True,
        ),
    ]

    model.fit(
        train_data,
        labels[train_idx],
        validation_data=(val_data, labels[val_idx]),
        epochs=EPOCHS,
        batch_size=BATCH_SIZE,
        callbacks=callbacks,
        verbose=1,
    )

    # Extract embeddings and compute recommendations
    print("\n[5/5] Generating embeddings and recommendations...")
    conn = get_connection()
    extract_and_store_embeddings(model, num_users, num_movies, idx_to_user, idx_to_movie, conn)
    compute_all_recommendations(conn)
    conn.close()

    print("\nTraining pipeline complete!")


if __name__ == "__main__":
    main()
