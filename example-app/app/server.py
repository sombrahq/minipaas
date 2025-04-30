#!/usr/bin/env python3
import http.server
import json
import logging

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

class RequestHandler(http.server.BaseHTTPRequestHandler):

    def do_GET(self):
        if self.path.startswith("/records"):
            self.handle_records_get()
        elif self.path.startswith("/queue"):
            self.handle_queue_get()
        elif self.path.startswith("/stream"):
            self.handle_stream_get()
        elif self.path.startswith("/consumers"):
            self.handle_consumers_get()
        elif self.path.startswith("/"):
            self.handle_health_get()
        else:
            self.send_error(404, "Resource not found")

    def do_POST(self):
        if self.path.startswith("/records"):
            self.handle_records_post()
        elif self.path.startswith("/queue"):
            self.handle_queue_post()
        elif self.path.startswith("/stream"):
            self.handle_stream_post()
        elif self.path.startswith("/error"):
            self.handle_error_post()
        else:
            self.send_error(404, "Resource not found")

    # ----- /records Handlers -----
    def handle_records_get(self):
        try:
            conn = get_connection()
            cur = conn.cursor()
            cur.execute("SELECT id, data FROM records;")
            rows = cur.fetchall()
            records = [{"id": r[0], "data": r[1]} for r in rows]
            cur.close()
            conn.close()
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            response_json = json.dumps({"records": records})
            self.send_header("Content-Length", str(len(response_json)))
            self.end_headers()
            self.wfile.write(response_json.encode("utf-8"))
        except Exception as e:
            self.send_error(500, str(e))

    def handle_records_post(self):
        try:
            content_length = int(self.headers.get("Content-Length", 0))
        except ValueError:
            content_length = 0
        body = self.rfile.read(content_length).decode("utf-8")
        try:
            payload = json.loads(body)
        except Exception:
            self.send_error(400, "Invalid JSON")
            return
        if "data" not in payload:
            self.send_error(400, "Missing 'data' field")
            return
        try:
            conn = get_connection()
            cur = conn.cursor()
            data = payload["data"]
            cur.execute("INSERT INTO records (data) VALUES (%s) RETURNING id;", (data,))
            new_id = cur.fetchone()[0]
            conn.commit()
            cur.close()
            conn.close()
            response = {"id": new_id, "data": data}
            response_json = json.dumps(response)
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(response_json)))
            self.end_headers()
            self.wfile.write(response_json.encode("utf-8"))
        except Exception as e:
            self.send_error(500, str(e))

    # ----- /requests Handlers (Queue) -----
    def handle_queue_get(self):
        try:
            conn = get_connection()
            cur = conn.cursor()
            cur.execute("SELECT id, payload, status, created_at, processed_at FROM example_queue ORDER BY id;")
            rows = cur.fetchall()
            queue = []
            for row in rows:
                queue.append({
                    "id": row[0],
                    "payload": row[1],
                    "status": row[2],
                    "created_at": row[3].isoformat() if row[3] else None,
                    "processed_at": row[4].isoformat() if row[4] else None,
                })
            cur.close()
            conn.close()
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            response_json = json.dumps({"queue": queue})
            self.send_header("Content-Length", str(len(response_json)))
            self.end_headers()
            self.wfile.write(response_json.encode("utf-8"))
        except Exception as e:
            self.send_error(500, str(e))

    # ----- / Handlers (Health) -----
    def handle_health_get(self):
        try:
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
        except Exception as e:
            self.send_error(500, str(e))

    def handle_queue_post(self):
        try:
            content_length = int(self.headers.get("Content-Length", 0))
        except ValueError:
            content_length = 0
        body = self.rfile.read(content_length).decode("utf-8")
        try:
            payload = json.loads(body)
        except Exception:
            self.send_error(400, "Invalid JSON")
            return
        if "payload" not in payload:
            self.send_error(400, "Missing 'payload' field")
            return
        try:
            conn = get_connection()
            cur = conn.cursor()
            data = payload["payload"]
            cur.execute("SELECT * FROM minipaas_queue_enqueue('example_queue', %s);", (json.dumps(data),))
            row = cur.fetchone()
            conn.commit()
            cur.close()
            conn.close()
            response = {
                "id": row[0],
                "payload": row[1],
                "status": row[2],
                "created_at": row[3].isoformat() if row[3] else None,
            }
            response_json = json.dumps(response)
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(response_json)))
            self.end_headers()
            self.wfile.write(response_json.encode("utf-8"))
        except Exception as e:
            self.send_error(500, str(e))

    # ----- /streams Handlers -----
    def handle_stream_post(self):
        # Add a new event into the stream (minipaas_stream_event table)
        try:
            content_length = int(self.headers.get("Content-Length", 0))
        except ValueError:
            content_length = 0
        body = self.rfile.read(content_length).decode("utf-8")
        try:
            payload = json.loads(body)
        except Exception:
            self.send_error(400, "Invalid JSON")
            return
        if "event" not in payload:
            self.send_error(400, "Missing 'event' field")
            return
        try:
            conn = get_connection()
            cur = conn.cursor()
            cur.execute(
                "SELECT * FROM minipaas_stream_publish('example_stream', %s);",
                (json.dumps(payload["event"]),)
            )
            row = cur.fetchone()
            conn.commit()
            cur.close()
            conn.close()
            response = {
                "id": row[0],
                "payload": row[1],
                "created_at": row[2].isoformat() if row[2] else None,
            }
            response_json = json.dumps(response)
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(response_json)))
            self.end_headers()
            self.wfile.write(response_json.encode("utf-8"))
        except Exception as e:
            self.send_error(500, str(e))

    # ----- /streams Handlers -----
    def handle_error_post(self):
        try:
            code = int(self.path.split("/")[-1])
            logging.error(f"Raising {code}!")
        except (ValueError, IndexError):
            code = 200

        try:
            self.send_response(code)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", "0")
            self.end_headers()
        except Exception as e:
            self.send_error(500, str(e))

    def handle_stream_get(self):
        try:
            conn = get_connection()
            cur = conn.cursor()
            cur.execute("SELECT id, payload, created_at FROM example_stream ORDER BY id;")
            rows = cur.fetchall()
            queue = []
            for row in rows:
                queue.append({
                    "id": row[0],
                    "payload": row[1],
                    "created_at": row[2].isoformat() if row[2] else None,
                })
            cur.close()
            conn.close()
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            response_json = json.dumps({"stream": queue})
            self.send_header("Content-Length", str(len(response_json)))
            self.end_headers()
            self.wfile.write(response_json.encode("utf-8"))
        except Exception as e:
            self.send_error(500, str(e))

    def handle_consumers_get(self):
        # Check consumer status from the consumer index table.
        try:
            conn = get_connection()
            cur = conn.cursor()
            cur.execute("SELECT * FROM consumer;")
            rows = cur.fetchall()
            consumers = []
            for row in rows:
                consumers.append({
                    "id": row[0],
                    "last_event_id": row[1],
                    "updated_at": row[2].isoformat() if row[2] else None,
                })
            cur.close()
            conn.close()
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            response_json = json.dumps({"consumers": consumers})
            self.send_header("Content-Length", str(len(response_json)))
            self.end_headers()
            self.wfile.write(response_json.encode("utf-8"))
        except Exception as e:
            self.send_error(500, str(e))

if __name__ == '__main__':
    port = 8080
    server_address = ("", port)
    try:
        conn = get_connection()
        conn.close()
        print("✅ Database connection successful.")
    except Exception as e:
        print("❌ Database connection failed:", e)
        exit(1)
    print(f"Starting server on port {port}...")
    httpd = http.server.HTTPServer(server_address, RequestHandler)
    httpd.serve_forever()
