# Use an official Python runtime as a parent image.
FROM python:3.11-alpine

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

# Set the working directory in the container.
WORKDIR /app

# Copy requirements.txt first to leverage Docker layer caching.
COPY requirements.txt .

# Install any dependencies specified in requirements.txt.
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of your application code.
COPY . .

# Expose the port on which the app runs.
EXPOSE 8080

# Run the Python script.
CMD ["python", "-u", "server.py"]
