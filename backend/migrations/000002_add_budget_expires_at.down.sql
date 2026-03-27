DROP INDEX IF EXISTS idx_budgets_expires_at;
ALTER TABLE budgets DROP COLUMN IF EXISTS expires_at;
