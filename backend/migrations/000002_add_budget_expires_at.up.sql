ALTER TABLE budgets
ADD COLUMN IF NOT EXISTS expires_at DATE;

CREATE INDEX IF NOT EXISTS idx_budgets_expires_at ON budgets(expires_at);
