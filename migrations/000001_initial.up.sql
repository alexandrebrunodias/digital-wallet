CREATE TABLE IF NOT EXISTS customers (
  id UUID PRIMARY KEY,
  name text NOT NULL,
  email text NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL,
    balance DECIMAL(12, 2),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY(customer_id) REFERENCES customers(id)
);

CREATE TABLE IF NOT EXISTS transactions (
  id UUID PRIMARY KEY,
  from_account_id UUID NOT NULL,
  to_account_id UUID NOT NULL,
  amount DECIMAL(14, 2),
  created_at TIMESTAMP,
  FOREIGN KEY(from_account_id) REFERENCES accounts(id),
  FOREIGN KEY(to_account_id) REFERENCES accounts(id)
);