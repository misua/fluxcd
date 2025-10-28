#!/bin/bash

API_URL="${API_URL:-http://localhost:8080}"

echo "🧪 Testing Go CRUD API at $API_URL"
echo "===================================="
echo ""

# Health check
echo "1️⃣ Health Check"
curl -s "$API_URL/health" | jq .
echo ""
echo ""

# Create item
echo "2️⃣ Create Item"
ITEM=$(curl -s -X POST "$API_URL/items" \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"MacBook Pro 16-inch"}')
echo "$ITEM" | jq .
ITEM_ID=$(echo "$ITEM" | jq -r '.id')
echo ""
echo ""

# Get all items
echo "3️⃣ Get All Items"
curl -s "$API_URL/items" | jq .
echo ""
echo ""

# Get specific item
echo "4️⃣ Get Item by ID ($ITEM_ID)"
curl -s "$API_URL/items/$ITEM_ID" | jq .
echo ""
echo ""

# Update item
echo "5️⃣ Update Item"
curl -s -X PUT "$API_URL/items/$ITEM_ID" \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"MacBook Pro 16-inch M3 Max"}' | jq .
echo ""
echo ""

# Get updated item
echo "6️⃣ Verify Update"
curl -s "$API_URL/items/$ITEM_ID" | jq .
echo ""
echo ""

# Create another item
echo "7️⃣ Create Another Item"
curl -s -X POST "$API_URL/items" \
  -H "Content-Type: application/json" \
  -d '{"name":"Mouse","description":"Logitech MX Master 3"}' | jq .
echo ""
echo ""

# Get all items again
echo "8️⃣ Get All Items (should have 2)"
curl -s "$API_URL/items" | jq .
echo ""
echo ""

# Delete item
echo "9️⃣ Delete Item ($ITEM_ID)"
curl -s -X DELETE "$API_URL/items/$ITEM_ID" -w "\nHTTP Status: %{http_code}\n"
echo ""
echo ""

# Verify deletion
echo "🔟 Get All Items (should have 1)"
curl -s "$API_URL/items" | jq .
echo ""
echo ""

echo "✅ All tests completed!"
