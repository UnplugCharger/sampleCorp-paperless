ALTER TABLE invoices ADD COLUMN IF NOT EXISTS pdf_url VARCHAR(255);
ALTER TABLE petty_cash ADD COLUMN IF NOT EXISTS pdf_url VARCHAR(255);
ALTER TABLE quotations ADD COLUMN IF NOT EXISTS pdf_url VARCHAR(255);
ALTER TABLE purchase_orders ADD COLUMN IF NOT EXISTS pdf_url VARCHAR(255);
ALTER TABLE payment_requests ADD COLUMN IF NOT EXISTS pdf_url VARCHAR(255);