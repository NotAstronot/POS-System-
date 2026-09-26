"""Test isolasi tenant: Tenant A tidak boleh melihat produk Tenant B, dan sebaliknya.

Jalankan:  python3 backend/postman/test-tenant-isolation.py [--base http://127.0.0.1:8080]
Keluar dengan kode 0 jika semua asersi lolos, 1 jika ada yang gagal.
Hanya memakai stdlib (urllib) agar tanpa dependensi.
"""
import json
import sys
import time
import urllib.request
import urllib.error

BASE = "http://127.0.0.1:8080"
for i, a in enumerate(sys.argv):
    if a == "--base" and i + 1 < len(sys.argv):
        BASE = sys.argv[i + 1]

TS = str(int(time.time()))
RESULTS = []


def req(method, path, body=None, token=None):
    data = json.dumps(body).encode() if body is not None else None
    r = urllib.request.Request(BASE + path, data=data, method=method)
    r.add_header("Content-Type", "application/json")
    if token:
        r.add_header("Authorization", "Bearer " + token)
    try:
        with urllib.request.urlopen(r, timeout=15) as resp:
            return resp.status, json.loads(resp.read().decode() or "{}")
    except urllib.error.HTTPError as e:
        try:
            return e.code, json.loads(e.read().decode() or "{}")
        except Exception:
            return e.code, {}


def check(name, cond, detail=""):
    RESULTS.append((name, bool(cond), detail))
    print(("PASS" if cond else "FAIL"), "-", name, detail)


# 1. Register dua tenant terpisah
s, ra = req("POST", "/api/v1/auth/register-merchant", {
    "merchant_name": "Tenant A Isolasi", "slug": f"tenant-a-{TS}",
    "username": f"owner_a_{TS}", "owner_name": "Owner A", "password": "secret123"})
check("register Tenant A -> 201", s == 201, f"status={s}")
tok_a = (ra.get("data") or {}).get("token", "")
tid_a = ((ra.get("data") or {}).get("claims") or {}).get("tenant_id")

s, rb = req("POST", "/api/v1/auth/register-merchant", {
    "merchant_name": "Tenant B Isolasi", "slug": f"tenant-b-{TS}",
    "username": f"owner_b_{TS}", "owner_name": "Owner B", "password": "secret123"})
check("register Tenant B -> 201", s == 201, f"status={s}")
tok_b = (rb.get("data") or {}).get("token", "")
tid_b = ((rb.get("data") or {}).get("claims") or {}).get("tenant_id")
check("token A & B diperoleh dan tenant_id berbeda",
      bool(tok_a) and bool(tok_b) and tid_a != tid_b,
      f"tenant_a={tid_a} tenant_b={tid_b}")

# 2. Masing-masing buat 1 produk (kode/barcode unik per run)
pa = {"code": f"A-{TS}", "name": f"Produk Milik A {TS}", "base_price": 10000,
      "unit": "pcs", "stock": 10, "barcode": f"111{TS}"}
s, ca = req("POST", "/api/v1/products", pa, tok_a)
check("Tenant A buat produk -> 201", s == 201, f"status={s}")
pid_a = (ca.get("data") or {}).get("id")

pb = {"code": f"B-{TS}", "name": f"Produk Milik B {TS}", "base_price": 20000,
      "unit": "pcs", "stock": 10, "barcode": f"222{TS}"}
s, cb = req("POST", "/api/v1/products", pb, tok_b)
check("Tenant B buat produk -> 201", s == 201, f"status={s}")
pid_b = (cb.get("data") or {}).get("id")

# 3. GET /products dengan JWT A: harus lihat produk A, TIDAK boleh lihat produk B
s, la = req("GET", "/api/v1/products", token=tok_a)
names_a = [p.get("name", "") for p in (la.get("data") or [])]
check("GET products (JWT A) -> 200", s == 200, f"status={s} jumlah={len(names_a)}")
check("JWT A melihat produk A", any(f"Produk Milik A {TS}" in n for n in names_a))
check("JWT A TIDAK melihat produk B", not any(f"Produk Milik B {TS}" in n for n in names_a),
      f"names={names_a}")

# 4. Kebalikannya dengan JWT B
s, lb = req("GET", "/api/v1/products", token=tok_b)
names_b = [p.get("name", "") for p in (lb.get("data") or [])]
check("GET products (JWT B) -> 200", s == 200, f"status={s} jumlah={len(names_b)}")
check("JWT B melihat produk B", any(f"Produk Milik B {TS}" in n for n in names_b))
check("JWT B TIDAK melihat produk A", not any(f"Produk Milik A {TS}" in n for n in names_b),
      f"names={names_b}")

# 5. Akses langsung lintas tenant via ID harus ditolak (404, bukan bocor 200)
s, _ = req("GET", f"/api/v1/products/{pid_b}", token=tok_a)
check("JWT A GET produk B by ID -> 404", s == 404, f"status={s}")
s, _ = req("GET", f"/api/v1/products/{pid_a}", token=tok_b)
check("JWT B GET produk A by ID -> 404", s == 404, f"status={s}")

# 6. Tanpa token tetap ditolak
s, _ = req("GET", "/api/v1/products")
check("GET products tanpa token -> 401", s == 401, f"status={s}")

failed = [r for r in RESULTS if not r[1]]
print(f"\n{len(RESULTS) - len(failed)}/{len(RESULTS)} lolos")
sys.exit(1 if failed else 0)
