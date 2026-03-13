#!/bin/bash

# Configuration
BACKEND_DIR="./backend"
FRONTEND_DIR="./frontend"
BACKEND_LOG="backend.log"
FRONTEND_LOG="frontend.log"
PID_FILE=".pid_tracker"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

start() {
    echo -e "${BLUE}Starting services...${NC}"

    # Start Backend
    if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
        echo -e "${RED}Backend port 8080 is already in use.${NC}"
    else
        echo "Starting Backend..."
        (cd "$BACKEND_DIR" && go run cmd/server/main.go > "../$BACKEND_LOG" 2>&1) &
        echo $! >> "$PID_FILE"
    fi

    # Start Frontend
    if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null ; then
        echo -e "${RED}Frontend port 5173 is already in use.${NC}"
    else
        echo "Starting Frontend..."
        (cd "$FRONTEND_DIR" && npm run dev > "../$FRONTEND_LOG" 2>&1) &
        echo $! >> "$PID_FILE"
    fi

    echo -e "${GREEN}Services started in background.${NC}"
    echo -e "Logs: $BACKEND_LOG, $FRONTEND_LOG"
}

stop() {
    echo -e "${BLUE}Stopping services...${NC}"
    
    # Kill by port for reliability
    BE_PID=$(lsof -t -i:8080)
    if [ ! -z "$BE_PID" ]; then
        echo "Stopping Backend (PID $BE_PID)..."
        kill $BE_PID
    fi

    FE_PID=$(lsof -t -i:5173)
    if [ ! -z "$FE_PID" ]; then
        echo "Stopping Frontend (PID $FE_PID)..."
        kill $FE_PID
    fi

    # Also kill by PIDs in file if they exist
    if [ -f "$PID_FILE" ]; then
        while read p; do
            kill -9 $p 2>/dev/null
        done < "$PID_FILE"
        rm "$PID_FILE"
    fi

    echo -e "${GREEN}Services stopped.${NC}"
}

status() {
    echo -e "${BLUE}Service Status:${NC}"
    
    if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
        echo -e "Backend:  ${GREEN}RUNNING${NC} (Port 8080)"
    else
        echo -e "Backend:  ${RED}STOPPED${NC}"
    fi

    if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null ; then
        echo -e "Frontend: ${GREEN}RUNNING${NC} (Port 5173)"
    else
        echo -e "Frontend: ${RED}STOPPED${NC}"
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        stop
        sleep 2
        start
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        exit 1
esac
