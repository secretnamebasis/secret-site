#!/bin/bash

source ./bin/config.sh

# Function to perform action and retrieve status
perform_action() {
    local method="$1"
    local url="$2"
    local data="$3"
    
    curl -s \
        -X "$method" \
        -H 'Content-Type: application/json' \
        -d "$data" \
        "$url" \
        | jq -r '.status'
}


# Define functions to perform actions
check_users() {
    perform_action \
        "GET" \
        "$ENDPOINT"
}

retrieve_user() {
    local id="$1"
    perform_action \
        "GET" \
        "$ENDPOINT/$id"
}

create_fail() {
    perform_action \
        "POST" \
        "$ENDPOINT" \
        '{
            "user":"secret",
            "wallet":"'"$CREATE_ADDRESS_FAIL"'"
        }'
}

create_success() {
    perform_action \
        "POST" \
        "$ENDPOINT" \
        '{
            "user":"secret",
            "wallet":"'"$CREATE_ADDRESS_SUCCESS"'"
        }'
}

update_fail() {
    local id="$1"
    perform_action \
        "PUT" \
        "$ENDPOINT/$id" \
        '{
            "user":"secret",
            "wallet":"'"$UPDATE_ADDRESS_FAIL"'"
        }'
}

update_success() {
    local id="$1"
    perform_action \
        "PUT" \
        "$ENDPOINT/$id" \
        '{
            "user":"secret",
            "wallet":"'"$UPDATE_ADDRESS_SUCCESS"'"
        }'
}

delete_user() {
    local id="$1"
    perform_action \
        "DELETE" \
        "$ENDPOINT/$id"
}

# Define actions and expected statuses
actions=(
    # create 
    "create_fail" "error"
    "retrieve_user 1" "error"
    "create_success" "success"
    "retrieve_user 1" "success"
    # update
    "update_fail 1" "error"
    "retrieve_user 1" "success"
    "update_success 1" "success"
    "retrieve_user 1" "success"
    #delete
    "delete_user 1" "success"
    "retrieve_user 1" "error"
)

# Loop through actions and execute them with sleep intervals
for ((i = 0; i < ${#actions[@]}; i+=2)); do
    action=${actions[i]}
    expected_status=${actions[i+1]}
    
    if [[ "$action" =~ ([a-zA-Z_]+)[[:space:]]*([0-9]*) ]]; then
        func_name="${BASH_REMATCH[1]}"
        id="${BASH_REMATCH[2]}"
        
        actual_status=$("$func_name" "$id")
    else
        actual_status=$("$action")
    fi
    
    echo "$action"
    if [[ "$actual_status" != "$expected_status" ]]; then
        echo "Failed $action : Expected $expected_status - Actual $actual_status"
    fi
    
    sleep 1
done
