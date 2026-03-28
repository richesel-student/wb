#!/bin/bash

curl -s "http://localhost:8080/events_for_day?user_id=1&date=2025-03-28" | jq