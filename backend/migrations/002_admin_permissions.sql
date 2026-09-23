CREATE TABLE IF NOT EXISTS permissions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_permissions (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (name, description) VALUES
    ('pos_access', 'Akses POS Kasir'),
    ('menu_manage', 'Kelola Menu Produk'),
    ('report_view', 'Lihat Laporan'),
    ('user_manage', 'Kelola User'),
    ('settings_manage', 'Kelola Pengaturan')
ON CONFLICT DO NOTHING;

ALTER TABLE products ADD COLUMN IF NOT EXISTS image_url TEXT NOT NULL DEFAULT '';

UPDATE users SET password = '$2b$10$U5VgPNLv9a26FkegqvmmUOrR7.Y/ihJblaIXAV2Q7W8L14heVlVzq' WHERE username = 'admin';
