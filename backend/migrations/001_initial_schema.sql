CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL DEFAULT '',
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    base_price NUMERIC(15,2) NOT NULL DEFAULT 0,
    unit VARCHAR(20) NOT NULL DEFAULT 'pcs',
    image_url TEXT NOT NULL DEFAULT '',
    category_id BIGINT REFERENCES categories(id),
    has_variants BOOLEAN NOT NULL DEFAULT false,
    tax_rate NUMERIC(5,2) NOT NULL DEFAULT 11.00,
    sort_order INT NOT NULL DEFAULT 0,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS product_availability (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(product_id, date)
);

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL DEFAULT '',
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'cashier',
    outlet_id VARCHAR(100) NOT NULL DEFAULT 'b0000000-0000-0000-0000-000000000001',
    tenant_id VARCHAR(100) NOT NULL DEFAULT 'a0000000-0000-0000-0000-000000000001',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO categories (name) VALUES
('Makanan'), ('Minuman'), ('Snack')
ON CONFLICT DO NOTHING;

INSERT INTO products (code, name, description, base_price, unit, category_id, has_variants, tax_rate, sort_order) VALUES
('BRGR-001', 'Beef Burger', 'Burger daging sapi dengan sayuran segar', 35000, 'pcs', 1, false, 11, 1),
('BRGR-002', 'Chicken Burger', 'Burger ayam crispy dengan mayo', 30000, 'pcs', 1, false, 11, 2),
('NASI-001', 'Nasi Goreng Spesial', 'Nasi goreng dengan telur dan ayam', 45000, 'pcs', 1, true, 11, 3),
('NASI-002', 'Nasi Goreng Ayam', 'Nasi goreng ayam sederhana', 38000, 'pcs', 1, false, 11, 4),
('NASI-003', 'Nasi Goreng Seafood', 'Nasi goreng dengan udang dan cumi', 50000, 'pcs', 1, false, 11, 5),
('MIE-001', 'Mie Goreng Spesial', 'Mie goreng dengan telur dan bakso', 40000, 'pcs', 1, false, 11, 6),
('MUM-001', 'Es Teh Manis', 'Teh manis segar dengan es batu', 8000, 'gelas', 2, false, 11, 1),
('MUM-002', 'Kopi Susu', 'Kopi susu gula aren', 25000, 'gelas', 2, true, 11, 2),
('MUM-003', 'Jus Alpukat', 'Jus alpukat segar dengan susu', 20000, 'gelas', 2, false, 11, 3),
('MUM-004', 'Jus Mangga', 'Jus mangga segar', 18000, 'gelas', 2, false, 11, 4),
('MUM-005', 'Air Mineral', 'Air mineral kemasan 600ml', 5000, 'botol', 2, false, 11, 5),
('MUM-006', 'Milkshake Coklat', 'Milkshake coklat creamy', 30000, 'gelas', 2, false, 11, 6),
('SNK-001', 'French Fries', 'Kentang goreng renyah', 18000, 'pcs', 3, false, 11, 1),
('SNK-002', 'Chicken Wings', 'Sayap ayam goreng pedas manis', 35000, 'pcs', 3, false, 11, 2),
('SNK-003', 'Onion Rings', 'Cincin bawang goreng crispy', 20000, 'pcs', 3, false, 11, 3),
('SNK-004', 'Spring Rolls', 'Lumpia isi sayuran', 22000, 'pcs', 3, false, 11, 4),
('SNK-005', 'Potato Wedges', 'Kentang panggang bumbu herbs', 25000, 'pcs', 3, false, 11, 5)
ON CONFLICT DO NOTHING;

INSERT INTO users (username, name, password, role) VALUES
('admin', 'Administrator', '$2b$10$U5VgPNLv9a26FkegqvmmUOrR7.Y/ihJblaIXAV2Q7W8L14heVlVzq', 'admin')
ON CONFLICT (username) DO NOTHING;

CREATE TABLE IF NOT EXISTS shifts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMPTZ,
    opening_balance NUMERIC(15,2) NOT NULL DEFAULT 0,
    closing_balance NUMERIC(15,2),
    status VARCHAR(20) NOT NULL DEFAULT 'open'
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    shift_id BIGINT NOT NULL REFERENCES shifts(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    total NUMERIC(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity INT NOT NULL DEFAULT 1,
    price NUMERIC(15,2) NOT NULL DEFAULT 0,
    subtotal NUMERIC(15,2) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id),
    amount NUMERIC(15,2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL DEFAULT 'cash',
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
