-- =====================================================
-- Migration 016: users.is_active
-- Dipakai gate login (AuthUsecase.Login) + toggle nonaktifkan user di admin.
-- =====================================================

ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true;
