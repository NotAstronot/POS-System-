# POS Tahap 1 - curl test (bash / Git Bash / WSL / Linux)
# Cara pakai:
#   cd backend && go run ./cmd/server   # terminal 1
#   bash backend/postman/test-tahap1.sh # terminal 2
# Ganti BASE ke http://localhost:8080 kalau jalan via docker-compose.

set -e
BASE="${BASE:-http://localhost:8082}"

echo "== 1. Register merchant =="
curl -s -X POST "$BASE/api/v1/auth/register-merchant" \
  -H "Content-Type: application/json" \
  -d '{"merchant_name":"Toko Maju Jaya","slug":"toko-maju-jaya","username":"owner","owner_name":"Budi Owner","password":"secret123"}' || echo "(register mungkin sudah pernah jalan, lanjut login)"

echo; echo "== 2. Login =="
TOKEN=$(curl -s -X POST "$BASE/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"owner","password":"secret123"}' | jq -r '.data.token')
echo "token: ${TOKEN:0:30}..."

echo; echo "== 3. List products =="
curl -s "$BASE/api/v1/products" -H "Authorization: Bearer $TOKEN" | jq .

echo; echo "== 4. Create product =="
PID=$(curl -s -X POST "$BASE/api/v1/products" -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"code":"PRD-001","name":"Kopi Susu 250ml","description":"Kopi susu gula aren","base_price":15000,"purchase_price":8000,"unit":"pcs","stock":100,"min_stock":10,"barcode":"8991234567890","has_variants":false,"tax_rate":0,"sort_order":1,"is_service":false}' | tee /dev/stderr | jq -r '.data.id // 1')
echo "productId=$PID"

echo; echo "== 5. Update product =="
curl -s -X PUT "$BASE/api/v1/products/$PID" -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"code":"PRD-001","name":"Kopi Susu 250ml","base_price":18000,"unit":"pcs","stock":95,"min_stock":10,"barcode":"8991234567890","has_variants":false,"tax_rate":0,"sort_order":1,"is_service":false}' | jq .

echo; echo "== 6. Sync transactions (offline batch) =="
curl -s -X POST "$BASE/api/v1/transactions/sync" -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"transactions\":[{\"client_order_id\":\"kasir01-20260926-001\",\"order_type\":\"dine_in\",\"table_number\":\"A1\",\"customer_name\":\"Andi\",\"items\":[{\"product_id\":$PID,\"quantity\":2,\"note\":\"es sedikit\"}],\"payments\":[{\"method\":\"cash\",\"amount\":36000}]},{\"client_order_id\":\"kasir01-20260926-002\",\"order_type\":\"takeaway\",\"items\":[{\"product_id\":$PID,\"quantity\":1}],\"payments\":[{\"method\":\"qris\",\"amount\":18000,\"reference_no\":\"QRIS123\"}]}]}" | jq .

echo "SELESAI"
