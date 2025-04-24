#!/usr/bin/env python3
import psycopg2

# Inline database configuration.
DB_HOST = "postgres"
DB_USER = "postgres"
DB_NAME = "postgres"
SSL_MODE = "disable"
SECRET_FILE = "/run/secrets/postgres_password"

# Read the password from the secret file.
try:
    with open(SECRET_FILE, "r") as f:
        DB_PASSWORD = f.read().strip()
except Exception as e:
    raise Exception("Error reading secret file: " + str(e))

# Build the connection string inline.
DB_CONN_STR = f"host={DB_HOST} user={DB_USER} password={DB_PASSWORD} dbname={DB_NAME} sslmode={SSL_MODE}"

def get_connection():
    return psycopg2.connect(DB_CONN_STR)


def cleanup_queue():
    try:
        conn = get_connection()
        cur = conn.cursor()
        cur.execute("DELETE FROM example_queue WHERE status = 'done';")
        cur.close()
        print("Queue cleaned", flush=True)
    except Exception as e:
        print(f"Error during dequeue: {e}", flush=True)


def cleanup_stream():
    try:
        conn = get_connection()
        cur = conn.cursor()
        cur.execute("SELECT last_event_id FROM consumer ORDER BY id DESC LIMIT 1;",)
        last_event_id = cur.fetchone()
        cur.execute("DELETE FROM example_stream WHERE id <= %s;", (last_event_id, ))
        cur.close()
        print("Stream cleaned", flush=True)
    except Exception as e:
        print(f"Error during dequeue: {e}", flush=True)


if __name__ == '__main__':
    cleanup_queue()
    cleanup_stream()
