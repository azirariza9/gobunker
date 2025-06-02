#!/bin/bash

# Log directories
LOG_DIR="./logs"
mkdir -p "$LOG_DIR"
touch "$LOG_DIR/tsc.log"
touch "$LOG_DIR/vite.log"

# Function to clean up processes
cleanup() {
    echo "Stopping application processes..."
    kill $GO_PID $TS_PID $TEST_PID $FRONTEND_PID 2>/dev/null
    echo "Stopping database containers..."
    (cd backend && podman compose -f compose.yaml down -v)
    exit 0
}

trap cleanup SIGINT SIGTERM

# Start database containers (detached)
(
    cd backend
    podman compose -f compose.yaml up -d
)

# Wait for database to be ready on host port 5433
echo "Waiting for database to be ready..."

while [ "$(podman inspect --format '{{.State.Health.Status}}' gobunker-db)" != "healthy" ]; do
    sleep 5
done
echo "Database is ready!"

# Start Go backend
(
    cd backend
    go run .
) &
GO_PID=$!

# Start frontend services
(
    cd frontend
    npm run ts > "../$LOG_DIR/tsc.log" 2>&1 &
    TS_PID=$!
    npm run test > "../$LOG_DIR/vite.log" 2>&1 &
    TEST_PID=$!
    wait
) &
FRONTEND_PID=$!

# Keep script running
wait