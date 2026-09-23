-- --------------------------------------------------------
-- Migration 011: Add asset type and depreciation method to fixed assets
-- --------------------------------------------------------

ALTER TABLE fixed_assets ADD COLUMN IF NOT EXISTS asset_type VARCHAR(50) NOT NULL DEFAULT 'peralatan';
ALTER TABLE fixed_assets ADD COLUMN IF NOT EXISTS depreciation_method VARCHAR(20) NOT NULL DEFAULT 'garis_lurus';

CREATE INDEX IF NOT EXISTS idx_fixed_assets_asset_type ON fixed_assets(asset_type);
CREATE INDEX IF NOT EXISTS idx_fixed_assets_method ON fixed_assets(depreciation_method);
