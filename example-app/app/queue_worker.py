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

# The notification channel for the queue "example_queue"
CHANNEL = "example_queue"

def get_connection():
    conn = psycopg2.connect(DB_CONN_STR)
    # Autocommit is required for LISTEN/NOTIFY.
    conn.set_isolation_level(psycopg2.extensions.ISOLATION_LEVEL_AUTOCOMMIT)
    return conn

def dequeue_messages(queue_name, batch_size=5):
    """
    Dequeue a batch of messages from the specified queue by calling the
    minipaas_queue_dequeue function.
    Returns a list of rows, each containing (id, payload, status, created_at).
    """
    try:
        conn = get_connection()
        cur = conn.cursor()
        cur.execute("SELECT id, payload, status, created_at FROM minipaas_queue_dequeue(%s, %s);",
                    (queue_name, batch_size))
        rows = cur.fetchall()
        cur.close()
        conn.close()
        return rows
    except Exception as e:
        print(f"Error during dequeue: {e}", flush=True)
        return []

def process_message(job_id, payload):
    """
    Process the message payload. Extend this function with your business logic.
    For now, simply print the job id and payload.
    """
    print(f"Processing job id {job_id} with payload: {payload}", flush=True)
    # Insert your processing logic here.

def ack_message(queue_name, job_id):
    """
    Acknowledge a message by calling the minipaas_queue_ack function.
    Opens a new connection so as not to interfere with the listener.
    """
    try:
        conn = get_connection()
        cur = conn.cursor()
        cur.execute("SELECT minipaas_queue_ack(%s, ARRAY[%s]::int[]);", (queue_name, job_id))
        conn.commit()
        cur.close()
        conn.close()
        print(f"Acked message with id {job_id}", flush=True)
    except Exception as e:
        print(f"Error acknowledging message {job_id}: {e}", flush=True)

def consume_messages(queue_name, batch_size=5):
    """
    Consume messages from the queue in batches until no more messages remain.
    """
    total = 0
    while True:
        messages = dequeue_messages(queue_name, batch_size)
        if not messages:
            break
        for row in messages:
            job_id, payload, status, created_at = row
            process_message(job_id, payload)
            ack_message(queue_name, job_id)
            print(f"Processed job {job_id} at {datetime.now().isoformat()}", flush=True)
            total += 1
    return total

def wait_for_notification(conn):
    """
    Wait for a notification on the connection.
    """
    print("Waiting for notification...", flush=True)
    # Wait with a timeout of 5 seconds; if nothing comes, loop until a notification arrives.
    while True:
        if select.select([conn], [], [], 5) != ([], [], []):
            conn.poll()
            if conn.notifies:
                # Clear all notifications (we use them as a trigger only).
                while conn.notifies:
                    _ = conn.notifies.pop(0)
                break

def run_worker():
    """
    Main consumer loop: first attempt to consume messages; if none,
    wait for a notification and then process available messages.
    """
    conn = get_connection()
    cur = conn.cursor()
    # Start listening on the channel.
    cur.execute(f"LISTEN {CHANNEL};")
    print(f"Listening on channel '{CHANNEL}' for queue notifications...", flush=True)
    cur.close()

    while True:
        # First, try to consume all pending messages.
        processed = consume_messages(CHANNEL, batch_size=5)
        if processed == 0:
            # No messages left, wait for a notification.
            wait_for_notification(conn)

def main():
    run_worker()

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print("\nExiting queue consumer.", flush=True)
        sys.exit(0)
