from urllib.parse import urlparse

import psycopg2
from pgvector.psycopg2 import register_vector

from config import DATABASE_URL


def get_connection():
    parsed = urlparse(DATABASE_URL)
    conn = psycopg2.connect(
        host=parsed.hostname,
        port=parsed.port or 5432,
        dbname=parsed.path.lstrip("/"),
        user=parsed.username,
        password=parsed.password,
        sslmode="disable",
    )
    register_vector(conn)
    return conn
