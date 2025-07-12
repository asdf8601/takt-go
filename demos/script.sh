#!/bin/bash

# Takt-go Demo Script - Enhanced version with realistic data
# This script demonstrates the main features of takt-go

# Function to simulate typing effect and execute command
typing_exec() {
    local command="$1"
    local delay=0.02
    local i=0

    # Type out the command
    printf "$ "
    while [ $i -lt ${#command} ]; do
        printf "%c" "${command:$i:1}"
        sleep $(echo "$delay + $(($RANDOM % 100)) * 0.001" | bc -l)
        i=$((i + 1))
    done

    echo ""

    # Execute the command
    eval "$command"
    sleep 1
    echo ""
}

# Function to simulate typing effect
typing() {
    local text="$1"
    local delay=0.03
    local i=0

    while [ $i -lt ${#text} ]; do
        printf "%c" "${text:$i:1}"
        sleep "$delay"
        i=$((i + 1))
    done

    echo ""
    sleep 1
}

create_csv() {
# Create a CSV file with realistic historical data
    cat > demo.csv << 'EOF'
timestamp,kind,notes
2025-07-10T17:30:00+02:00,out,End of day
2025-07-10T12:30:00+02:00,in,Back from lunch
2025-07-10T12:00:00+02:00,out,Lunch break
2025-07-10T09:00:00+02:00,in,Starting work on API refactor
2025-07-09T17:45:00+02:00,out,Late finish - debugging
2025-07-09T13:00:00+02:00,in,Back from lunch
2025-07-09T12:00:00+02:00,out,Lunch break
2025-07-09T08:45:00+02:00,in,Early start - fixing production bug
2025-07-08T18:00:00+02:00,out,End of day
2025-07-08T10:00:00+02:00,in,Starting work on new feature
2025-07-07T18:00:00+02:00,out,End of day
2025-07-07T10:00:00+02:00,in,Starting work on new feature
2025-07-06T18:00:00+02:00,out,End of day
2025-07-06T10:00:00+02:00,in,Starting work on new feature
2025-07-05T18:00:00+02:00,out,End of day
2025-07-05T10:00:00+02:00,in,Starting work on new feature
2025-07-04T18:00:00+02:00,out,End of day
2025-07-04T10:00:00+02:00,in,Starting work on new feature
2025-07-03T18:00:00+02:00,out,End of day
2025-07-03T10:00:00+02:00,in,Starting work on new feature
2025-07-02T18:00:00+02:00,out,End of day
2025-07-02T10:00:00+02:00,in,Starting work on new feature
2025-07-01T18:00:00+02:00,out,End of day
2025-07-01T10:00:00+02:00,in,Starting work on new feature
2025-06-30T18:00:00+02:00,out,End of day
2025-06-30T10:00:00+02:00,in,Starting work on new feature
2025-06-29T18:00:00+02:00,out,End of day
2025-06-29T10:00:00+02:00,in,Starting work on new feature
2025-06-28T18:00:00+02:00,out,End of day
2025-06-28T10:00:00+02:00,in,Starting work on new feature
2025-06-27T18:00:00+02:00,out,End of day
2025-06-27T10:00:00+02:00,in,Starting work on new feature
2025-06-26T18:00:00+02:00,out,End of day
2025-06-26T10:00:00+02:00,in,Starting work on new feature
2025-06-25T18:00:00+02:00,out,End of day
2025-06-25T10:00:00+02:00,in,Starting work on new feature
2025-06-24T18:00:00+02:00,out,End of day
2025-06-24T10:00:00+02:00,in,Starting work on new feature
2025-06-23T18:00:00+02:00,out,End of day
2025-06-23T10:00:00+02:00,in,Starting work on new feature
2025-06-22T18:00:00+02:00,out,End of day
2025-06-22T10:00:00+02:00,in,Starting work on new feature
2025-06-21T18:00:00+02:00,out,End of day
2025-06-21T10:00:00+02:00,in,Starting work on new feature
2025-06-20T18:00:00+02:00,out,End of day
2025-06-20T10:00:00+02:00,in,Starting work on new feature
2025-06-19T18:00:00+02:00,out,End of day
2025-06-19T10:00:00+02:00,in,Starting work on new feature
2025-06-18T18:00:00+02:00,out,End of day
2025-06-18T10:00:00+02:00,in,Starting work on new feature
2025-06-17T18:00:00+02:00,out,End of day
2025-06-17T10:00:00+02:00,in,Starting work on new feature
2025-06-16T18:00:00+02:00,out,End of day
2025-06-16T10:00:00+02:00,in,Starting work on new feature
2025-06-15T18:00:00+02:00,out,End of day
2025-06-15T10:00:00+02:00,in,Starting work on new feature
2025-06-14T18:00:00+02:00,out,End of day
2025-06-14T10:00:00+02:00,in,Starting work on new feature
2025-06-13T18:00:00+02:00,out,End of day
2025-06-13T10:00:00+02:00,in,Starting work on new feature
2025-06-12T18:00:00+02:00,out,End of day
2025-06-12T10:00:00+02:00,in,Starting work on new feature
2025-06-11T18:00:00+02:00,out,End of day
2025-06-11T10:00:00+02:00,in,Starting work on new feature
2025-06-10T18:00:00+02:00,out,End of day
2025-06-10T10:00:00+02:00,in,Starting work on new feature
2025-06-09T18:00:00+02:00,out,End of day
2025-06-09T10:00:00+02:00,in,Starting work on new feature
2025-06-08T18:00:00+02:00,out,End of day
2025-06-08T10:00:00+02:00,in,Starting work on new feature
2025-06-07T18:00:00+02:00,out,End of day
2025-06-07T10:00:00+02:00,in,Starting work on new feature
2025-06-06T18:00:00+02:00,out,End of day
2025-06-06T10:00:00+02:00,in,Starting work on new feature
2025-06-05T18:00:00+02:00,out,End of day
2025-06-05T10:00:00+02:00,in,Starting work on new feature
2025-06-04T18:00:00+02:00,out,End of day
2025-06-04T10:00:00+02:00,in,Starting work on new feature
2025-06-03T18:00:00+02:00,out,End of day
2025-06-03T10:00:00+02:00,in,Starting work on new feature
2025-06-02T18:00:00+02:00,out,End of day
2025-06-02T10:00:00+02:00,in,Starting work on new feature
2025-06-01T18:00:00+02:00,out,End of day
2025-06-01T10:00:00+02:00,in,Starting work on new feature
2025-05-31T18:00:00+02:00,out,End of day
2025-05-31T10:00:00+02:00,in,Starting work on new feature
2025-05-30T18:00:00+02:00,out,End of day
2025-05-30T10:00:00+02:00,in,Starting work on new feature
2025-05-29T18:00:00+02:00,out,End of day
2025-05-29T10:00:00+02:00,in,Starting work on new feature
2025-05-28T18:00:00+02:00,out,End of day
2025-05-28T10:00:00+02:00,in,Starting work on new feature
2025-05-27T18:00:00+02:00,out,End of day
2025-05-27T10:00:00+02:00,in,Starting work on new feature
2025-05-26T18:00:00+02:00,out,End of day
2025-05-26T10:00:00+02:00,in,Starting work on new feature
2025-05-25T18:00:00+02:00,out,End of day
2025-05-25T10:00:00+02:00,in,Starting work on new feature
2025-05-24T18:00:00+02:00,out,End of day
2025-05-24T10:00:00+02:00,in,Starting work on new feature
2025-05-23T18:00:00+02:00,out,End of day
2025-05-23T10:00:00+02:00,in,Starting work on new feature
2025-05-22T18:00:00+02:00,out,End of day
2025-05-22T10:00:00+02:00,in,Starting work on new feature
2025-05-21T18:00:00+02:00,out,End of day
2025-05-21T10:00:00+02:00,in,Starting work on new feature
2025-05-20T18:00:00+02:00,out,End of day
2025-05-20T10:00:00+02:00,in,Starting work on new feature
2025-05-19T18:00:00+02:00,out,End of day
2025-05-19T10:00:00+02:00,in,Starting work on new feature
2025-05-18T18:00:00+02:00,out,End of day
2025-05-18T10:00:00+02:00,in,Starting work on new feature
2025-05-17T18:00:00+02:00,out,End of day
2025-05-17T10:00:00+02:00,in,Starting work on new feature
2025-05-16T18:00:00+02:00,out,End of day
2025-05-16T10:00:00+02:00,in,Starting work on new feature
2025-05-15T18:00:00+02:00,out,End of day
2025-05-15T10:00:00+02:00,in,Starting work on new feature
2025-05-14T18:00:00+02:00,out,End of day
2025-05-14T10:00:00+02:00,in,Starting work on new feature
2025-05-13T18:00:00+02:00,out,End of day
2025-05-13T10:00:00+02:00,in,Starting work on new feature
2025-05-12T18:00:00+02:00,out,End of day
2025-05-12T10:00:00+02:00,in,Starting work on new feature
2025-05-11T18:00:00+02:00,out,End of day
2025-05-11T10:00:00+02:00,in,Starting work on new feature
2025-05-10T18:00:00+02:00,out,End of day
2025-05-10T10:00:00+02:00,in,Starting work on new feature
2025-05-09T18:00:00+02:00,out,End of day
2025-05-09T10:00:00+02:00,in,Starting work on new feature
2025-05-08T18:00:00+02:00,out,End of day
2025-05-08T10:00:00+02:00,in,Starting work on new feature
2025-05-07T18:00:00+02:00,out,End of day
2025-05-07T10:00:00+02:00,in,Starting work on new feature
2025-05-06T18:00:00+02:00,out,End of day
2025-05-06T10:00:00+02:00,in,Starting work on new feature
2025-05-05T18:00:00+02:00,out,End of day
2025-05-05T10:00:00+02:00,in,Starting work on new feature
2025-05-04T18:00:00+02:00,out,End of day
2025-05-04T10:00:00+02:00,in,Starting work on new feature
2025-05-03T18:00:00+02:00,out,End of day
2025-05-03T10:00:00+02:00,in,Starting work on new feature
2025-05-02T18:00:00+02:00,out,End of day
2025-05-02T10:00:00+02:00,in,Starting work on new feature
2025-05-01T18:00:00+02:00,out,End of day
2025-05-01T10:00:00+02:00,in,Starting work on new feature
2025-04-30T18:00:00+02:00,out,End of day
2025-04-30T10:00:00+02:00,in,Starting work on new feature
2025-04-29T18:00:00+02:00,out,End of day
2025-04-29T10:00:00+02:00,in,Starting work on new feature
2025-04-28T18:00:00+02:00,out,End of day
2025-04-28T10:00:00+02:00,in,Starting work on new feature
2025-04-27T18:00:00+02:00,out,End of day
2025-04-27T10:00:00+02:00,in,Starting work on new feature
2025-04-26T18:00:00+02:00,out,End of day
2025-04-26T10:00:00+02:00,in,Starting work on new feature
2025-04-25T18:00:00+02:00,out,End of day
2025-04-25T10:00:00+02:00,in,Starting work on new feature
2025-04-24T18:00:00+02:00,out,End of day
2025-04-24T10:00:00+02:00,in,Starting work on new feature
2025-04-23T18:00:00+02:00,out,End of day
2025-04-23T10:00:00+02:00,in,Starting work on new feature
2025-04-22T18:00:00+02:00,out,End of day
2025-04-22T10:00:00+02:00,in,Starting work on new feature
2025-04-21T18:00:00+02:00,out,End of day
2025-04-21T10:00:00+02:00,in,Starting work on new feature
2025-04-20T18:00:00+02:00,out,End of day
2025-04-20T10:00:00+02:00,in,Starting work on new feature
2025-04-19T18:00:00+02:00,out,End of day
2025-04-19T10:00:00+02:00,in,Starting work on new feature
2025-04-18T18:00:00+02:00,out,End of day
2025-04-18T10:00:00+02:00,in,Starting work on new feature
2025-04-17T18:00:00+02:00,out,End of day
2025-04-17T10:00:00+02:00,in,Starting work on new feature
2025-04-16T18:00:00+02:00,out,End of day
2025-04-16T10:00:00+02:00,in,Starting work on new feature
2025-04-15T18:00:00+02:00,out,End of day
2025-04-15T10:00:00+02:00,in,Starting work on new feature
2025-04-14T18:00:00+02:00,out,End of day
2025-04-14T10:00:00+02:00,in,Starting work on new feature
2025-04-13T18:00:00+02:00,out,End of day
2025-04-13T10:00:00+02:00,in,Starting work on new feature
2025-04-12T18:00:00+02:00,out,End of day
2025-04-12T10:00:00+02:00,in,Starting work on new feature
2025-04-11T18:00:00+02:00,out,End of day
2025-04-11T10:00:00+02:00,in,Starting work on new feature
2025-04-10T18:00:00+02:00,out,End of day
2025-04-10T10:00:00+02:00,in,Starting work on new feature
2025-04-09T18:00:00+02:00,out,End of day
2025-04-09T10:00:00+02:00,in,Starting work on new feature
2025-04-08T18:00:00+02:00,out,End of day
2025-04-08T10:00:00+02:00,in,Starting work on new feature
2025-04-07T18:00:00+02:00,out,End of day
2025-04-07T10:00:00+02:00,in,Starting work on new feature
2025-04-06T18:00:00+02:00,out,End of day
2025-04-06T10:00:00+02:00,in,Starting work on new feature
2025-04-05T18:00:00+02:00,out,End of day
2025-04-05T10:00:00+02:00,in,Starting work on new feature
2025-04-04T18:00:00+02:00,out,End of day
2025-04-04T10:00:00+02:00,in,Starting work on new feature
2025-04-03T18:00:00+02:00,out,End of day
2025-04-03T10:00:00+02:00,in,Starting work on new feature
2025-04-02T18:00:00+02:00,out,End of day
2025-04-02T10:00:00+02:00,in,Starting work on new feature
2025-04-01T18:00:00+02:00,out,End of day
2025-04-01T10:00:00+02:00,in,Starting work on new feature
2025-03-31T18:00:00+02:00,out,End of day
2025-03-31T10:00:00+02:00,in,Starting work on new feature
2025-03-30T18:00:00+02:00,out,End of day
2025-03-30T10:00:00+02:00,in,Starting work on new feature
2025-03-29T18:00:00+02:00,out,End of day
2025-03-29T10:00:00+02:00,in,Starting work on new feature
2025-03-28T18:00:00+02:00,out,End of day
2025-03-28T10:00:00+02:00,in,Starting work on new feature
2025-03-27T18:00:00+02:00,out,End of day
2025-03-27T10:00:00+02:00,in,Starting work on new feature
2025-03-26T18:00:00+02:00,out,End of day
2025-03-26T10:00:00+02:00,in,Starting work on new feature
2025-03-25T18:00:00+02:00,out,End of day
2025-03-25T10:00:00+02:00,in,Starting work on new feature
2025-03-24T18:00:00+02:00,out,End of day
2025-03-24T10:00:00+02:00,in,Starting work on new feature
2025-03-23T18:00:00+02:00,out,End of day
2025-03-23T10:00:00+02:00,in,Starting work on new feature
2025-03-22T18:00:00+02:00,out,End of day
2025-03-22T10:00:00+02:00,in,Starting work on new feature
2025-03-21T18:00:00+02:00,out,End of day
2025-03-21T10:00:00+02:00,in,Starting work on new feature
2025-03-20T18:00:00+02:00,out,End of day
2025-03-20T10:00:00+02:00,in,Starting work on new feature
2025-03-19T18:00:00+02:00,out,End of day
2025-03-19T10:00:00+02:00,in,Starting work on new feature
2025-03-18T18:00:00+02:00,out,End of day
2025-03-18T10:00:00+02:00,in,Starting work on new feature
2025-03-17T18:00:00+02:00,out,End of day
2025-03-17T10:00:00+02:00,in,Starting work on new feature
2025-03-16T18:00:00+02:00,out,End of day
2025-03-16T10:00:00+02:00,in,Starting work on new feature
2025-03-15T18:00:00+02:00,out,End of day
2025-03-15T10:00:00+02:00,in,Starting work on new feature
2025-03-14T18:00:00+02:00,out,End of day
2025-03-14T10:00:00+02:00,in,Starting work on new feature
2025-03-13T18:00:00+02:00,out,End of day
2025-03-13T10:00:00+02:00,in,Starting work on new feature
2025-03-12T18:00:00+02:00,out,End of day
2025-03-12T10:00:00+02:00,in,Starting work on new feature
2025-03-11T18:00:00+02:00,out,End of day
2025-03-11T10:00:00+02:00,in,Starting work on new feature
2025-03-10T18:00:00+02:00,out,End of day
2025-03-10T10:00:00+02:00,in,Starting work on new feature
2025-03-09T18:00:00+02:00,out,End of day
2025-03-09T10:00:00+02:00,in,Starting work on new feature
2025-03-08T18:00:00+02:00,out,End of day
2025-03-08T10:00:00+02:00,in,Starting work on new feature
2025-03-07T18:00:00+02:00,out,End of day
2025-03-07T10:00:00+02:00,in,Starting work on new feature
2025-03-06T18:00:00+02:00,out,End of day
2025-03-06T10:00:00+02:00,in,Starting work on new feature
2025-03-05T18:00:00+02:00,out,End of day
2025-03-05T10:00:00+02:00,in,Starting work on new feature
2025-03-04T18:00:00+02:00,out,End of day
2025-03-04T10:00:00+02:00,in,Starting work on new feature
2025-03-03T18:00:00+02:00,out,End of day
2025-03-03T10:00:00+02:00,in,Starting work on new feature
2025-03-02T18:00:00+02:00,out,End of day
2025-03-02T10:00:00+02:00,in,Starting work on new feature
2025-03-01T18:00:00+02:00,out,End of day
2025-03-01T10:00:00+02:00,in,Starting work on new feature
2025-02-28T18:00:00+02:00,out,End of day
2025-02-28T10:00:00+02:00,in,Starting work on new feature
2025-02-27T18:00:00+02:00,out,End of day
2025-02-27T10:00:00+02:00,in,Starting work on new feature
2025-02-26T18:00:00+02:00,out,End of day
2025-02-26T10:00:00+02:00,in,Starting work on new feature
2025-02-25T18:00:00+02:00,out,End of day
2025-02-25T10:00:00+02:00,in,Starting work on new feature
2025-02-24T18:00:00+02:00,out,End of day
2025-02-24T10:00:00+02:00,in,Starting work on new feature
2025-02-23T18:00:00+02:00,out,End of day
2025-02-23T10:00:00+02:00,in,Starting work on new feature
2025-02-22T18:00:00+02:00,out,End of day
2025-02-22T10:00:00+02:00,in,Starting work on new feature
2025-02-21T18:00:00+02:00,out,End of day
2025-02-21T10:00:00+02:00,in,Starting work on new feature
2025-02-20T18:00:00+02:00,out,End of day
2025-02-20T10:00:00+02:00,in,Starting work on new feature
2025-02-19T18:00:00+02:00,out,End of day
2025-02-19T10:00:00+02:00,in,Starting work on new feature
2025-02-18T18:00:00+02:00,out,End of day
2025-02-18T10:00:00+02:00,in,Starting work on new feature
2025-02-17T18:00:00+02:00,out,End of day
2025-02-17T10:00:00+02:00,in,Starting work on new feature
2025-02-16T18:00:00+02:00,out,End of day
2025-02-16T10:00:00+02:00,in,Starting work on new feature
2025-02-15T18:00:00+02:00,out,End of day
2025-02-15T10:00:00+02:00,in,Starting work on new feature
2025-02-14T18:00:00+02:00,out,End of day
2025-02-14T10:00:00+02:00,in,Starting work on new feature
2025-02-13T18:00:00+02:00,out,End of day
2025-02-13T10:00:00+02:00,in,Starting work on new feature
2025-02-12T18:00:00+02:00,out,End of day
2025-02-12T10:00:00+02:00,in,Starting work on new feature
2025-02-11T18:00:00+02:00,out,End of day
2025-02-11T10:00:00+02:00,in,Starting work on new feature
2025-02-10T18:00:00+02:00,out,End of day
2025-02-10T10:00:00+02:00,in,Starting work on new feature
2025-02-09T18:00:00+02:00,out,End of day
2025-02-09T10:00:00+02:00,in,Starting work on new feature
2025-02-08T18:00:00+02:00,out,End of day
2025-02-08T10:00:00+02:00,in,Starting work on new feature
2025-02-07T18:00:00+02:00,out,End of day
2025-02-07T10:00:00+02:00,in,Starting work on new feature
2025-02-06T18:00:00+02:00,out,End of day
2025-02-06T10:00:00+02:00,in,Starting work on new feature
2025-02-05T18:00:00+02:00,out,End of day
2025-02-05T10:00:00+02:00,in,Starting work on new feature
2025-02-04T18:00:00+02:00,out,End of day
2025-02-04T10:00:00+02:00,in,Starting work on new feature
2025-02-03T18:00:00+02:00,out,End of day
2025-02-03T10:00:00+02:00,in,Starting work on new feature
2025-02-02T18:00:00+02:00,out,End of day
2025-02-02T10:00:00+02:00,in,Starting work on new feature
2025-02-01T18:00:00+02:00,out,End of day
2025-02-01T10:00:00+02:00,in,Starting work on new feature
2025-01-31T18:00:00+02:00,out,End of day
2025-01-31T10:00:00+02:00,in,Starting work on new feature
2025-01-30T18:00:00+02:00,out,End of day
2025-01-30T10:00:00+02:00,in,Starting work on new feature
2025-01-29T18:00:00+02:00,out,End of day
2025-01-29T10:00:00+02:00,in,Starting work on new feature
2025-01-28T18:00:00+02:00,out,End of day
2025-01-28T10:00:00+02:00,in,Starting work on new feature
2025-01-27T18:00:00+02:00,out,End of day
2025-01-27T10:00:00+02:00,in,Starting work on new feature
2025-01-26T18:00:00+02:00,out,End of day
2025-01-26T10:00:00+02:00,in,Starting work on new feature
2025-01-25T18:00:00+02:00,out,End of day
2025-01-25T10:00:00+02:00,in,Starting work on new feature
2025-01-24T18:00:00+02:00,out,End of day
2025-01-24T10:00:00+02:00,in,Starting work on new feature
2025-01-23T18:00:00+02:00,out,End of day
2025-01-23T10:00:00+02:00,in,Starting work on new feature
2025-01-22T18:00:00+02:00,out,End of day
2025-01-22T10:00:00+02:00,in,Starting work on new feature
2025-01-21T18:00:00+02:00,out,End of day
2025-01-21T10:00:00+02:00,in,Starting work on new feature
2025-01-20T18:00:00+02:00,out,End of day
2025-01-20T10:00:00+02:00,in,Starting work on new feature
2025-01-19T18:00:00+02:00,out,End of day
2025-01-19T10:00:00+02:00,in,Starting work on new feature
2025-01-18T18:00:00+02:00,out,End of day
2025-01-18T10:00:00+02:00,in,Starting work on new feature
2025-01-17T18:00:00+02:00,out,End of day
2025-01-17T10:00:00+02:00,in,Starting work on new feature
2025-01-16T18:00:00+02:00,out,End of day
2025-01-16T10:00:00+02:00,in,Starting work on new feature
2025-01-15T18:00:00+02:00,out,End of day
2025-01-15T10:00:00+02:00,in,Starting work on new feature
2025-01-14T18:00:00+02:00,out,End of day
2025-01-14T10:00:00+02:00,in,Starting work on new feature
2025-01-13T18:00:00+02:00,out,End of day
2025-01-13T10:00:00+02:00,in,Starting work on new feature
2025-01-12T18:00:00+02:00,out,End of day
2025-01-12T10:00:00+02:00,in,Starting work on new feature
2025-01-11T18:00:00+02:00,out,End of day
2025-01-11T10:00:00+02:00,in,Starting work on new feature
2025-01-10T18:00:00+02:00,out,End of day
2025-01-10T10:00:00+02:00,in,Starting work on new feature
2025-01-09T18:00:00+02:00,out,End of day
2025-01-09T10:00:00+02:00,in,Starting work on new feature
2025-01-08T18:00:00+02:00,out,End of day
2025-01-08T10:00:00+02:00,in,Starting work on new feature
2025-01-07T18:00:00+02:00,out,End of day
2025-01-07T10:00:00+02:00,in,Starting work on new feature
2025-01-06T18:00:00+02:00,out,End of day
2025-01-06T10:00:00+02:00,in,Starting work on new feature
2025-01-05T18:00:00+02:00,out,End of day
2025-01-05T10:00:00+02:00,in,Starting work on new feature
2025-01-04T18:00:00+02:00,out,End of day
2025-01-04T10:00:00+02:00,in,Starting work on new feature
2025-01-03T18:00:00+02:00,out,End of day
2025-01-03T10:00:00+02:00,in,Starting work on new feature
2025-01-02T18:00:00+02:00,out,End of day
2025-01-02T10:00:00+02:00,in,Starting work on new feature
2025-01-01T18:00:00+02:00,out,End of day
2025-01-01T10:00:00+02:00,in,Starting work on new feature
EOF
}

clear

typing "=== Takt-go Demo ==="
typing "Time tracking made simple with Go"

# Set up demo environment
export TAKT_FILE=demo.csv
export TAKT_TARGET_HOURS=8


# Clean up any existing demo file
rm -f demo.csv demo.csv.bak

typing "Setting up demo with some historical data..."
create_csv

typing "Takt is governed by a csv file"
echo ""
typing_exec "cat demo.csv | head -10"
typing_exec "export TAKT_FILE=demo.csv"
sleep 3

typing "Let's set the target working hours per day:"
echo ""
typing_exec "export TAKT_TARGET_HOURS=8"
sleep 3
clear

typing "1. Let's check in to start tracking today:"
echo ""
typing_exec "takt check \"HEY! Working on takt-go demo\""
sleep 3

typing "View recent entries:"
echo ""
typing_exec "takt cat 5"
sleep 3

typing "Check daily summary (last 5 days):"
echo ""
typing_exec "takt day 5"
sleep 3

typing "View weekly summary (last 2 weeks):"
echo ""
typing_exec "takt week 2"
sleep 3

typing "Check monthly summary:"
echo ""
typing_exec "takt month 1"
sleep 3

clear

typing "Let's check out for a break:"
echo ""
typing_exec "takt c \"Break time\""
sleep 3

typing "View again the recent entries:"
echo ""
typing_exec "takt cat 5"
sleep 3

typing_exec "takt d 1"
sleep 3


typing "You can use the grid view for better visualization:"
echo ""
typing_exec "takt grid 2025"
sleep 5
clear


typing "Check available commands:"
echo ""
typing_exec "takt --help | head -20"

sleep 3
clear

typing "Demo completed! 🎉"
typing "Key features demonstrated:"
typing "- [x] Simple check in/out with notes"
typing "- [x] Historical data viewing"
typing "- [x] Daily/weekly/monthly summaries"
typing "- [x] Overtime/undertime balance tracking"
typing "- [x] Flexible time format support"
typing "Learn more at: https://github.com/asdf8601/takt-go"

# Clean up
rm -f demo.csv demo.csv.bak
