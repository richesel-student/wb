#!/bin/bash

BASE_URL="http://localhost:8080"

echo "🚀 TEST START"

echo "-----------------------------"
echo "1. CREATE (POST)"
curl -s -X POST "$BASE_URL/items" \
-H "Content-Type: application/json" \
-d '{"asset_id":1,"type":"buy","amount_usd":1000}'
echo -e "\n"

echo "-----------------------------"
echo "2. GET ALL"
curl -s "$BASE_URL/items"
echo -e "\n"

echo "-----------------------------"
echo "3. UPDATE (PUT id=1)"
curl -s -X PUT "$BASE_URL/items/1" \
-H "Content-Type: application/json" \
-d '{"amount_usd":2000}'
echo -e "\n"

echo "-----------------------------"
echo "4. GET AFTER UPDATE"
curl -s "$BASE_URL/items"
echo -e "\n"

echo "-----------------------------"
echo "5. DELETE (id=1)"
curl -s -X DELETE "$BASE_URL/items/1"
echo -e "\n"

echo "-----------------------------"
echo "6. GET AFTER DELETE"
curl -s "$BASE_URL/items"
echo -e "\n"

echo "-----------------------------"
echo "✅ TEST DONE"