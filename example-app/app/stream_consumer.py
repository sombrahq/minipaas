#!/usr/bin/env python3
import sys
import select
import psycopg2
import psycopg2.extensions
from datetime import datetime

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
    sys.exit("Error reading secret file: " + str(e))

# Build the connection string.
DB_CONN_STR = f"host={DB_HOST} user={DB_USER} password={DB_PASSWORD} dbname={DB_NAME} sslmode={SSL_MODE}"

# Consumer ID is now dictated by the consumer.
CONSUMER_ID = "3b469834-6551-4c97-9004-7aabeade4d49"

# The notification channel for the stream (and also the stream name for consumption)
CHANNEL = "example_stream"

def get_connection():
    conn = psycopg2.connect(DB_CONN_STR)
    # Autocommit is required for LISTEN/NOTIFY.
    conn.set_isolation_level(psycopg2.extensions.ISOLATION_LEVEL_AUTOCOMMIT)
    return conn

def ensure_consumer_exists(conn):
    """
    Ensure a consumer record exists for the given consumer id.
    If it does not exist, create one with last_event_id set to 0.
    """
    cur = conn.cursor()
    cur.execute("SELECT id, last_event_id FROM consumer WHERE id = %s;", (CONSUMER_ID,))
    row = cur.fetchone()
    if row:
        consumer_id, last_event_id = row
        print(f"Consumer '{CONSUMER_ID}' exists with last_event_id = {last_event_id}", flush=True)
    else:
        cur.execute("INSERT INTO consumer (id, last_event_id) VALUES (%s, 0) RETURNING id, last_event_id;",
                    (CONSUMER_ID,))
        row = cur.fetchone()
        consumer_id, last_event_id = row
        print(f"Consumer '{CONSUMER_ID}' created with id {consumer_id} and last_event_id = {last_event_id}", flush=True)
    cur.close()
    return consumer_id, last_event_id

def update_consumer_index(conn, consumer_id, new_event_id):
    """Update the consumer record with the new event id and current timestamp."""
    cur = conn.cursor()
    cur.execute("UPDATE consumer SET last_event_id = %s, updated_at = now() WHERE id = %s;",
                (new_event_id, consumer_id))
    cur.close()

def subscribe_stream_batch(last_id, batch_size=5):
    """
    Consume a batch of events from the stream by calling the
    minipaas_stream_consume function.
    Returns a list of rows, each containing (id, payload, created_at).
    """
    try:
        conn = get_connection()
        cur = conn.cursor()
        cur.execute("SELECT id, payload, created_at FROM minipaas_stream_consume(%s, %s, %s);",
                    (CHANNEL, last_id, batch_size))
        rows = cur.fetchall()
        cur.close()
        conn.close()
        return rows
    except Exception as e:
        print(f"Error subscribing to stream: {e}", flush=True)
        return []

def process_event(payload):
    """
    Process the event payload.
    Extend this function with your business logic.
    For now, simply print the payload.
    """
    print(f"Processing event: {payload}", flush=True)
    # Insert actual processing logic here.

def wait_for_notification(conn):
    """
    Wait for a notification on the connection.
    """
    print("Waiting for notification...", flush=True)
    while True:
        if select.select([conn], [], [], 5) != ([], [], []):
            conn.poll()
            if conn.notifies:
                # Clear notifications as trigger and exit.
                while conn.notifies:
                    _ = conn.notifies.pop(0)
                break

def run_consumer():
    """
    Main consumer loop: continuously consume events starting from the last
    processed event, and only wait for notifications if no more events are available.
    """
    conn = get_connection()
    consumer_id, last_event_id = ensure_consumer_exists(conn)
    print(f"Starting stream consumer with last_event_id: {last_event_id}", flush=True)

    cur = conn.cursor()
    cur.execute(f"LISTEN {CHANNEL};")
    print(f"Listening on channel '{CHANNEL}' as consumer '{CONSUMER_ID}'...", flush=True)
    cur.close()

    while True:
        # First, attempt to consume all available events.
        processed = False
        while True:
            events = subscribe_stream_batch(last_event_id, batch_size=5)
            if not events:
                break
            processed = True
            for event in events:
                new_id, payload, created_at = event
                process_event(payload)
                last_event_id = new_id
                update_consumer_index(conn, consumer_id, new_id)
                print(f"Consumer index updated to {new_id} at {datetime.now().isoformat()}", flush=True)
        if not processed:
            # No events processed in this iteration, wait for a notification.
            wait_for_notification(conn)

def main():
    run_consumer()

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print("\nExiting stream consumer.", flush=True)
        sys.exit(0)
