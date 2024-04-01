#!/bin/bash

# Logger function to write messages to a log file
log() {
    local log_file="./logs/run_integration_test.log"
    local status="$1"
    local message="$2"
    touch $log_file
    echo "$(date +'%Y-%m-%d %H:%M:%S') [${status}] ${message}" >> "$log_file"
}

# Function to echo messages with a specified status and date
echo_with_status() {
    local status="$1"
    local message="$2"
    local date_time=$(date +'%Y-%m-%d %H:%M:%S')
    echo "$date_time - status: $status - $message"
}

log_and_echo() {
    local status="$1"
    local message="$2"
    
    # Log the message
    log "$status" "$message"
    
    # Echo the message with status and date
    echo_with_status "$status" "$message"
}


# Function to check if a command is available
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if Go is installed
if ! command_exists "go"; then
    log "error" "Go is not installed. Please install Go before running this script."
    exit 1
fi

# Build the executable
go build -o ./builds/ .

# Check if the build was successful
if [ $? -ne 0 ]; then
    log "error" "Failed to build the executable."
    exit 1
fi

# Check if Screen is installed
if ! command_exists "screen"; then
    log "error" "Screen is not installed. Please install Screen before running this script, eg sudo apt install screen."
    exit 1
fi

log_and_echo "info" "secret-site built under ./builds."

log_and_echo "info" "secret-site is testing api"

test_output=$( go test ./test/api/api_test.go -v -parallel 2)

echo $test_output > test_results.tmp

# Check if the test completed successfully
if [[ $test_output != *"ok"* ]]; then
    log_and_echo "error" "Test failed. Test output: $test_output"
    exit 1
fi

# Check if the test completed successfully
if [ $? -ne 0 ]; then
    log_and_echo "error" "secret-site test failed, see logs."
    exit 1
fi

log_and_echo "info" "test completed successfully."

# Remove the temporary file
rm test_results.tmp
