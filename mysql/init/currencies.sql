CREATE TABLE IF NOT EXISTS currencies (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  type VARCHAR(255) NOT NULL,
  chain VARCHAR(255) DEFAULT NULL,
  name VARCHAR(255) NOT NULL,
  iso VARCHAR(255) NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE (chain, iso)
);

INSERT INTO currencies (type, chain, name, iso) VALUES 
  ('fiat', NULL, 'EURO', 'EUR'),
  ('fiat', NULL, 'US Dollar', 'USD'),
  ('crypto', 'ethereum', 'Ethereum', 'ETH'),
  ('crypto', 'ethereum', 'Tether ERC-20', 'USDT'),
  ('crypto', 'binance', 'BNB Binance', 'BNB'),
  ('crypto', 'binance', 'Tether BEP-20', 'USDT')
;
