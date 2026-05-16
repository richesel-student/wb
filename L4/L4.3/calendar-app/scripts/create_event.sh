#!/bin/bash

curl -s -X POST "http://localhost:8080/create_event" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "date": "2025-03-28",
    "event": "Test event"
  }' | jq