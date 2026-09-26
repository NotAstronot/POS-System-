# POS Tahap 1 - curl test (PowerShell)
# Cara pakai:
#   1. Jalankan backend:  cd backend; go run ./cmd/server
#   2. Jalankan script:   powershell -ExecutionPolicy Bypass -File backend/postman/test-tahap1.ps1
#   Ganti $BASE ke http://localhost:8080 kalau jalan via docker-compose.

$BASE = "http://localhost:8082"

Write-Host "== 1. Register merchant =="
try {
  $reg = Invoke-RestMethod -Uri "$BASE/api/v1/auth/register-merchant" -Method Post -ContentType "application/json" -Body '{"merchant_name":"Toko Maju Jaya","slug":"toko-maju-jaya","username":"owner","owner_name":"Budi Owner","password":"secret123"}'
  $reg | ConvertTo-Json -Depth 5
} catch {
  Write-Host "(register mungkin sudah pernah jalan, lanjut login) $($_.Exception.Message)"
}

Write-Host "`n== 2. Login =="
$login = Invoke-RestMethod -Uri "$BASE/api/v1/auth/login" -Method Post -ContentType "application/json" -Body '{"username":"owner","password":"secret123"}'
$TOKEN = $login.data.token
Write-Host "token: $($TOKEN.Substring(0, [Math]::Min(30, $TOKEN.Length)))..."
$H = @{ Authorization = "Bearer $TOKEN" }

Write-Host "`n== 3. List products =="
Invoke-RestMethod -Uri "$BASE/api/v1/products" -Headers $H | ConvertTo-Json -Depth 5

Write-Host "`n== 4. Create product =="
$created = Invoke-RestMethod -Uri "$BASE/api/v1/products" -Method Post -Headers $H -ContentType "application/json" -Body '{"code":"PRD-001","name":"Kopi Susu 250ml","description":"Kopi susu gula aren","base_price":15000,"purchase_price":8000,"unit":"pcs","stock":100,"min_stock":10,"barcode":"8991234567890","has_variants":false,"tax_rate":0,"sort_order":1,"is_service":false}'
$created | ConvertTo-Json -Depth 5
$PID_ = $created.data.id
if (-not $PID_) { $PID_ = 1 }
Write-Host "productId=$PID_"

Write-Host "`n== 5. Update product =="
Invoke-RestMethod -Uri "$BASE/api/v1/products/$PID_" -Method Put -Headers $H -ContentType "application/json" -Body '{"code":"PRD-001","name":"Kopi Susu 250ml","base_price":18000,"unit":"pcs","stock":95,"min_stock":10,"barcode":"8991234567890","has_variants":false,"tax_rate":0,"sort_order":1,"is_service":false}' | ConvertTo-Json -Depth 5

Write-Host "`n== 6. Sync transactions (offline batch) =="
$syncBody = "{`"transactions`":[{`"client_order_id`":`"kasir01-20260926-001`",`"order_type`":`"dine_in`",`"table_number`":`"A1`",`"customer_name`":`"Andi`",`"items`":[{`"product_id`":$PID_,`"quantity`":2,`"note`":`"es sedikit`"}],`"payments`":[{`"method`":`"cash`",`"amount`":36000}]},{`"client_order_id`":`"kasir01-20260926-002`",`"order_type`":`"takeaway`",`"items`":[{`"product_id`":$PID_,`"quantity`":1}],`"payments`":[{`"method`":`"qris`",`"amount`":18000,`"reference_no`":`"QRIS123`"}]}]}"
Invoke-RestMethod -Uri "$BASE/api/v1/transactions/sync" -Method Post -Headers $H -ContentType "application/json" -Body $syncBody | ConvertTo-Json -Depth 8

Write-Host "`nSELESAI"
