#!/bin/bash

# Function to get all network interfaces and their IPs
get_interfaces() {
    interfaces=($(ip -o -4 addr show | awk '{print $2 " " $4}'))
    echo "Available network interfaces with their IPs:"
    for i in "${!interfaces[@]}"; do
        if (( i % 2 == 0 )); then
            ip=${interfaces[i + 1]%/*}  # Remove the subnet mask
            echo "$((i / 2)). ${interfaces[i]} - $ip"
        fi
    done
}

# Function to run ./dra command
run_dra() {
    local remote_address_port=$1
    local interface=$2
    local ip=$3
    local local_port=$4
    echo "Running ./dra on interface $interface ($ip:$local_port)"
    ./dra -raddr=$remote_address_port -laddr=$ip:$local_port -transport=sctp
    exit_code=$?
    if [ $exit_code -eq 0 ]; then
        echo "Connection successful on interface $interface ($ip:$local_port)"
    else
        echo "Connection failed on interface $interface ($ip:$local_port) with exit code $exit_code"
    fi
}

if [ $# -ne 2 ]; then
    echo "Usage: $0 <remote_address>:<port> <local_port>"
    exit 1
fi

remote_address_port=$1
local_port=$2

get_interfaces

read -p "Enter the interface numbers you want to test (space-separated): " -a selected_interfaces

for index in "${selected_interfaces[@]}"; do
    if (( index * 2 < ${#interfaces[@]} )); then
        interface=${interfaces[index * 2]}
        ip=${interfaces[index * 2 + 1]%/*}  # Remove the subnet mask
        run_dra $remote_address_port $interface $ip $local_port
    else
        echo "Invalid interface number: $index"
    fi
done

echo "All tests completed."
