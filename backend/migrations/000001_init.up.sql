CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(200) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  phone VARCHAR(30),
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS families (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(200) NOT NULL,
  owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS family_members (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role VARCHAR(50) NOT NULL,
  joined_at TIMESTAMP NOT NULL DEFAULT now(),
  UNIQUE(family_id, user_id)
);

CREATE TABLE IF NOT EXISTS wallets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
  name VARCHAR(200) NOT NULL,
  type VARCHAR(50) NOT NULL,
  balance NUMERIC(18,2) NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
  name VARCHAR(200) NOT NULL,
  type VARCHAR(50) NOT NULL,
  icon VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
  wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE RESTRICT,
  category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
  amount NUMERIC(18,2) NOT NULL,
  type VARCHAR(50) NOT NULL,
  note TEXT,
  transaction_date DATE NOT NULL,
  created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS budgets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
  category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
  amount NUMERIC(18,2) NOT NULL,
  month INT NOT NULL,
  year INT NOT NULL,
  UNIQUE(family_id, category_id, month, year)
);

CREATE TABLE IF NOT EXISTS goals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
  name VARCHAR(200) NOT NULL,
  target_amount NUMERIC(18,2) NOT NULL,
  current_amount NUMERIC(18,2) NOT NULL DEFAULT 0,
  target_date DATE
);

CREATE TABLE IF NOT EXISTS bills (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
  name VARCHAR(200) NOT NULL,
  amount NUMERIC(18,2) NOT NULL,
  due_day INT NOT NULL,
  recurring BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title VARCHAR(200) NOT NULL,
  message TEXT NOT NULL,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS invitations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  family_id UUID NOT NULL REFERENCES families(id) ON DELETE CASCADE,
  email VARCHAR(255) NOT NULL,
  token TEXT NOT NULL UNIQUE,
  status VARCHAR(50) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_families_owner_id ON families(owner_id);
CREATE INDEX IF NOT EXISTS idx_family_members_family_id ON family_members(family_id);
CREATE INDEX IF NOT EXISTS idx_family_members_user_id ON family_members(user_id);
CREATE INDEX IF NOT EXISTS idx_wallets_family_id ON wallets(family_id);
CREATE INDEX IF NOT EXISTS idx_categories_family_id ON categories(family_id);
CREATE INDEX IF NOT EXISTS idx_transactions_family_id ON transactions(family_id);
CREATE INDEX IF NOT EXISTS idx_transactions_wallet_id ON transactions(wallet_id);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(transaction_date);
CREATE INDEX IF NOT EXISTS idx_budgets_family_id ON budgets(family_id);
CREATE INDEX IF NOT EXISTS idx_goals_family_id ON goals(family_id);
CREATE INDEX IF NOT EXISTS idx_bills_family_id ON bills(family_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_invitations_family_id ON invitations(family_id);
CREATE INDEX IF NOT EXISTS idx_invitations_email ON invitations(email);
