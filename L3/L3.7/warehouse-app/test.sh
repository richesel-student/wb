#!/bin/bash

BASE_URL="http://localhost:8080"

echo "🔐 LOGIN..."

TOKEN=$(curl -s -X POST "$BASE_URL/login" \
-H "Content-Type: application/json" \
-d '{"username":"admin","password":"admin"}' | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')

if [ -z "$TOKEN" ]; then
  echo "❌ LOGIN FAILED"
  exit 1
fi

echo "✅ TOKEN: $TOKEN"
echo "--------------------------------"

echo "📦 CREATE ITEM..."

curl -s -X POST "$BASE_URL/api/items" \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{"name":"apple","quantity":10}' | cat

echo -e "\n--------------------------------"

echo "📋 GET ITEMS..."

ITEMS=$(curl -s "$BASE_URL/api/items" \
-H "Authorization: Bearer $TOKEN")

echo "$ITEMS"

# вытаскиваем id первого элемента
ITEM_ID=$(echo "$ITEMS" | sed -n 's/.*"id":\([0-9]*\).*/\1/p' | head -n1)

if [ -z "$ITEM_ID" ]; then
  echo "❌ NO ITEMS FOUND"
  exit 1
fi

echo "✅ ITEM_ID: $ITEM_ID"
echo "--------------------------------"

echo "✏️ UPDATE ITEM..."

curl -s -X PUT "$BASE_URL/api/items/$ITEM_ID" \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{"name":"apple updated","quantity":20}' | cat

echo -e "\n--------------------------------"

echo "📜 GET HISTORY (after update)..."

curl -s "$BASE_URL/api/items/$ITEM_ID/history" \
-H "Authorization: Bearer $TOKEN" | cat

echo -e "\n--------------------------------"

echo "❌ DELETE ITEM..."

curl -s -X DELETE "$BASE_URL/api/items/$ITEM_ID" \
-H "Authorization: Bearer $TOKEN" | cat

echo -e "\n--------------------------------"

echo "📜 GET HISTORY (after delete)..."

curl -s "$BASE_URL/api/items/$ITEM_ID/history" \
-H "Authorization: Bearer $TOKEN" | cat

echo -e "\n--------------------------------"

echo "🏁 DONE"