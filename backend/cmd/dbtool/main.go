package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=pos_system sslmode=disable connect_timeout=5"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("ping: %v", err)
	}

	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "apply":
		files, err := filepath.Glob("migrations/*.sql")
		if err != nil {
			log.Fatal(err)
		}
		sort.Strings(files)
		for _, f := range files {
			base := filepath.Base(f)
			if base <= "002_admin_permissions.sql" || base >= "016_add_user_is_active.sql" {
				continue // 001+002 sudah; 016 sudah
			}
			b, err := os.ReadFile(f)
			if err != nil {
				log.Fatalf("read %s: %v", base, err)
			}
			if _, err := db.Exec(string(b)); err != nil {
				log.Fatalf("apply %s: %v", base, err)
			}
			fmt.Println("OK", base)
		}
		fmt.Println("semua migrasi 003-015 diterapkan")

	case "check":
		markers := []struct {
			label, query string
		}{
			{"003 branches", `SELECT to_regclass('public.branches') IS NOT NULL`},
			{"004 employees", `SELECT to_regclass('public.employees') IS NOT NULL`},
			{"005 cash_transactions", `SELECT to_regclass('public.cash_transactions') IS NOT NULL`},
			{"006 purchase_orders", `SELECT to_regclass('public.purchase_orders') IS NOT NULL`},
			{"007 sales_quotations", `SELECT to_regclass('public.sales_quotations') IS NOT NULL`},
			{"008 warehouses", `SELECT to_regclass('public.warehouses') IS NOT NULL`},
			{"009 bank_transfers", `SELECT to_regclass('public.bank_transfers') IS NOT NULL`},
			{"010 journal_entries", `SELECT to_regclass('public.journal_entries') IS NOT NULL`},
			{"011 fixed_assets", `SELECT to_regclass('public.fixed_assets') IS NOT NULL`},
			{"012 orders.order_number", `SELECT COUNT(*) FROM information_schema.columns WHERE table_name='orders' AND column_name='order_number'`},
			{"013 marketplace_connections", `SELECT to_regclass('public.marketplace_connections') IS NOT NULL`},
			{"014 orders.client_order_id", `SELECT COUNT(*) FROM information_schema.columns WHERE table_name='orders' AND column_name='client_order_id'`},
			{"015 tenants", `SELECT to_regclass('public.tenants') IS NOT NULL`},
			{"015 products.tenant_id", `SELECT COUNT(*) FROM information_schema.columns WHERE table_name='products' AND column_name='tenant_id'`},
			{"015 RLS products", `SELECT relrowsecurity FROM pg_class WHERE relname='products'`},
			{"016 users.is_active", `SELECT COUNT(*) FROM information_schema.columns WHERE table_name='users' AND column_name='is_active'`},
		}
		for _, m := range markers {
			var v bool
			if err := db.QueryRow(m.query).Scan(&v); err != nil {
				var n int
				if err2 := db.QueryRow(m.query).Scan(&n); err2 == nil {
					fmt.Printf("%-35s = %d\n", m.label, n)
					continue
				}
				fmt.Printf("%-35s = ERROR %v\n", m.label, err)
				continue
			}
			fmt.Printf("%-35s = %v\n", m.label, v)
		}
		var tenants, users, warehouses int
		db.QueryRow("SELECT COUNT(*) FROM tenants").Scan(&tenants)
		db.QueryRow("SELECT COUNT(*) FROM users").Scan(&users)
		db.QueryRow("SELECT COUNT(*) FROM warehouses").Scan(&warehouses)
		fmt.Printf("state: tenants=%d users=%d warehouses=%d\n", tenants, users, warehouses)

	case "verify":
		dump := func(label, query string) {
			fmt.Println("---", label)
			rows, err := db.Query(query)
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}
			defer rows.Close()
			cols, _ := rows.Columns()
			fmt.Println(cols)
			for rows.Next() {
				vals := make([]interface{}, len(cols))
				ptrs := make([]interface{}, len(cols))
				for i := range vals {
					ptrs[i] = &vals[i]
				}
				if err := rows.Scan(ptrs...); err != nil {
					fmt.Println("scan err:", err)
					return
				}
				out := make([]string, len(cols))
				for i, v := range vals {
					switch t := v.(type) {
					case nil:
						out[i] = "NULL"
					case []byte:
						out[i] = string(t)
					default:
						out[i] = fmt.Sprint(t)
					}
				}
				fmt.Println(out)
			}
		}
		dump("tenants", "SELECT id, name, slug, status FROM tenants ORDER BY id")
		dump("users", "SELECT id, username, name, role, tenant_id, is_active FROM users ORDER BY id")
		dump("products", "SELECT id, code, name, base_price, stock, tenant_id FROM products ORDER BY id")
		dump("orders", "SELECT id, order_number, user_id, total, discount_amount, payment_method, status, client_order_id, tenant_id FROM orders ORDER BY id")
		dump("order_items", "SELECT id, order_id, product_id, quantity, price, subtotal, tenant_id FROM order_items ORDER BY id")
		dump("transactions", "SELECT id, order_id, amount, payment_method, status, tenant_id FROM transactions ORDER BY id")
		dump("warehouses", "SELECT id, code, name, tenant_id FROM warehouses ORDER BY id")

	default:
		fmt.Println("usage: dbtool [apply|check|verify]")
	}
}
