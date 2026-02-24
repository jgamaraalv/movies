import tensorflow as tf
from config import EMBEDDING_DIM, MLP_LAYERS, DROPOUT_RATE, GENRE_EMBEDDING_DIM


def build_ncf_model(num_users, num_movies, num_genres, max_genres=5):
    """Build Neural Collaborative Filtering model (GMF + MLP)."""
    # Inputs
    user_input = tf.keras.layers.Input(shape=(1,), name="user_input")
    movie_input = tf.keras.layers.Input(shape=(1,), name="movie_input")
    genre_input = tf.keras.layers.Input(shape=(max_genres,), name="genre_input")

    # GMF embeddings
    user_gmf_emb = tf.keras.layers.Embedding(
        num_users + 1, EMBEDDING_DIM, name="user_gmf_embedding"
    )(user_input)
    user_gmf_emb = tf.keras.layers.Flatten()(user_gmf_emb)

    movie_gmf_emb = tf.keras.layers.Embedding(
        num_movies + 1, EMBEDDING_DIM, name="movie_gmf_embedding"
    )(movie_input)
    movie_gmf_emb = tf.keras.layers.Flatten()(movie_gmf_emb)

    # MLP embeddings
    user_mlp_emb = tf.keras.layers.Embedding(
        num_users + 1, EMBEDDING_DIM, name="user_mlp_embedding"
    )(user_input)
    user_mlp_emb = tf.keras.layers.Flatten()(user_mlp_emb)

    movie_mlp_emb = tf.keras.layers.Embedding(
        num_movies + 1, EMBEDDING_DIM, name="movie_mlp_embedding"
    )(movie_input)
    movie_mlp_emb = tf.keras.layers.Flatten()(movie_mlp_emb)

    # Genre embedding
    genre_emb = tf.keras.layers.Embedding(
        num_genres + 1, GENRE_EMBEDDING_DIM, name="genre_embedding"
    )(genre_input)
    genre_emb = tf.keras.layers.GlobalAveragePooling1D()(genre_emb)

    # GMF path: element-wise product
    gmf_output = tf.keras.layers.Multiply()([user_gmf_emb, movie_gmf_emb])

    # MLP path
    mlp_input = tf.keras.layers.Concatenate()([user_mlp_emb, movie_mlp_emb, genre_emb])
    mlp_output = mlp_input
    for units in MLP_LAYERS:
        mlp_output = tf.keras.layers.Dense(units, activation="relu")(mlp_output)
        mlp_output = tf.keras.layers.Dropout(DROPOUT_RATE)(mlp_output)

    # NeuMF: concatenate GMF + MLP
    neumf = tf.keras.layers.Concatenate()([gmf_output, mlp_output])
    output = tf.keras.layers.Dense(1, activation="sigmoid", name="prediction")(neumf)

    model = tf.keras.Model(
        inputs=[user_input, movie_input, genre_input],
        outputs=output,
        name="NCF",
    )

    return model
