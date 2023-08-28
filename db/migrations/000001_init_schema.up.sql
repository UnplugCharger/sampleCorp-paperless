CREATE TABLE users (
  id bigserial PRIMARY KEY NOT NULL,
  username varchar UNIQUE NOT NULL,
  hashed_password varchar NOT NULL,
  full_name varchar NOT NULL,
  email varchar UNIQUE NOT NULL,
  department varchar NOT NULL CHECK (department IN ('ACCOUNTS', 'PROCUREMENT','FINANCE','ADMIN')),
  password_changed_at timestamptz NOT NULL DEFAULT (now()),
  created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE user_roles (
  id bigserial PRIMARY KEY NOT NULL,
  user_id bigserial NOT NULL,
  role_id bigserial NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now()),
  terminated_at timestamptz DEFAULT (now())
);

CREATE TABLE roles (
  id bigserial PRIMARY KEY NOT NULL,
  name varchar UNIQUE NOT NULL,
  description varchar
);

CREATE TABLE companies (
id bigserial PRIMARY KEY NOT NULL,
name varchar UNIQUE NOT NULL,
initials varchar NOT NULL,
address varchar
);

CREATE TABLE signatories (
 id bigserial PRIMARY KEY NOT NULL,
 name varchar NOT NULL,
 title varchar NOT NULL
);


CREATE TABLE petty_cash (
transaction_id SERIAL PRIMARY KEY,
petty_cash_no VARCHAR(50) UNIQUE NOT NULL,
employee_id INTEGER NOT NULL REFERENCES users (id),
amount numeric(14,4),
currency_code char(3),
description VARCHAR(255) NOT NULL,
transaction_date DATE NOT NULL,
updated_at DATE NOT NULL DEFAULT CURRENT_DATE,
approved_at DATE,
folio VARCHAR(50) NOT NULL DEFAULT 'PETTY CASH' CHECK (folio IN ('PETTY CASH', 'BANK', 'CASH','FUEL')),
debit_account VARCHAR(50) NOT NULL,
status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'DECLINED')),
authorised_by INTEGER  REFERENCES users (id)
);
CREATE TABLE bank_details (
id SERIAL PRIMARY KEY,
bank_name VARCHAR(50) NOT NULL,
account_name VARCHAR(50) NOT NULL,
account_number VARCHAR(50) UNIQUE NOT NULL,
branch VARCHAR(50) NOT NULL,
swift_code VARCHAR(50) NOT NULL,
address VARCHAR(255) NOT NULL,
country VARCHAR(50) NOT NULL,
currency VARCHAR(50) NOT NULL,
account_type VARCHAR(50) NOT NULL,
company_id INTEGER NOT NULL REFERENCES companies(id)

);

CREATE TABLE quotations (
id SERIAL PRIMARY KEY,
quotation_no VARCHAR(50) NOT NULL UNIQUE,
quotation_revision INTEGER DEFAULT 0 NOT NULL,
quotation_revision_number INTEGER DEFAULT 0 NOT NULL,
date DATE DEFAULT CURRENT_DATE,
attn VARCHAR(100) NOT NULL,
company_id INTEGER NOT NULL REFERENCES companies(id),
site VARCHAR(100) NOT NULL,
validity INTEGER NOT NULL,
warranty INTEGER NOT NULL,
payment_terms TEXT NOT NULL,
delivery_terms TEXT NOT NULL,
signatory_id INTEGER NOT NULL REFERENCES signatories(id),
status VARCHAR(50) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'SUCCESSFUL', 'UNSUCCESSFUL')),
sent_or_received VARCHAR(10) CHECK (sent_or_received IN ('sent', 'received'))
);

CREATE TABLE purchase_orders (
 id SERIAL PRIMARY KEY,
 po_no VARCHAR(50) UNIQUE NOT NULL,
 date DATE DEFAULT CURRENT_DATE,
 attn VARCHAR(100) NOT NULL,
 company_id INTEGER NOT NULL REFERENCES companies(id),
 address VARCHAR(255) NOT NULL,
 signatory_id INTEGER NOT NULL REFERENCES signatories(id),
 quotation_id INTEGER REFERENCES quotations(id),
 po_status VARCHAR(50) DEFAULT 'PENDING',
 approved_by INTEGER REFERENCES users (id),
 sent_or_received VARCHAR(10) CHECK (sent_or_received IN ('sent', 'received'))
);



CREATE TABLE quotation_items(
 id SERIAL PRIMARY KEY,
 description TEXT NOT NULL ,
 uom VARCHAR(10) NOT NULL,
 qty INTEGER NOT NULL,
 lead_time VARCHAR(50) NOT NULL,
 item_price FLOAT NOT NULL,
 disc FLOAT NOT NULL DEFAULT 0,
 unit_price FLOAT NOT NULL,
 net_price FLOAT NOT NULL,
 currency VARCHAR(10) NOT NULL,
 quotation_id INTEGER REFERENCES quotations(id) NOT NULL
);



CREATE TABLE invoices (
id SERIAL PRIMARY KEY,
invoice_no VARCHAR(50) UNIQUE NOT NULL,
purchase_order_number VARCHAR(50) REFERENCES purchase_orders(po_no) NOT NULL,
date DATE DEFAULT CURRENT_DATE,
attn VARCHAR(100) NOT NULL,
company_id INTEGER NOT NULL REFERENCES companies(id),
site VARCHAR(100) NOT NULL,
amount_due FLOAT NOT NULL,
bank_details INTEGER NOT NULL REFERENCES bank_details(id),
signatory_id INTEGER NOT NULL REFERENCES signatories(id),
sent_or_received VARCHAR(10) NOT NULL CHECK (sent_or_received IN ('sent', 'received'))
);
CREATE TABLE payment_requests (
request_id SERIAL PRIMARY KEY,
payment_request_no VARCHAR(50) UNIQUE NOT NULL,
amount_in_words VARCHAR(255) NOT NULL,
employee_id INTEGER NOT NULL REFERENCES users (id),
currency VARCHAR(10) NOT NULL,
amount FLOAT NOT NULL,
description VARCHAR(255) NOT NULL,
request_date DATE NOT NULL DEFAULT CURRENT_DATE,
status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'DECLINED')),
admin_id INTEGER REFERENCES users (id),
approval_date DATE ,
invoice_id INTEGER REFERENCES invoices (id)

);

CREATE TABLE invoice_items (
id SERIAL PRIMARY KEY,
description TEXT NOT NULL ,
uom VARCHAR(10) NOT NULL,
qty INTEGER NOT NULL,
unit_price FLOAT NOT NULL,
net_price FLOAT NOT NULL,
currency VARCHAR(10) NOT NULL,
invoice_id INTEGER REFERENCES invoices(id) NOT NULL
);


CREATE TABLE purchase_order_items (
id SERIAL PRIMARY KEY,
description TEXT NOT NULL,
part_no VARCHAR(50) NOT NULL,
uom VARCHAR(10) NOT NULL,
qty INTEGER NOT NULL DEFAULT 1,
item_price FLOAT NOT NULL DEFAULT 0,
discount FLOAT NOT NULL DEFAULT 0 ,
net_price FLOAT NOT NULL DEFAULT 0,
net_value FLOAT NOT NULL DEFAULT 0,
currency VARCHAR(10) NOT NULL,
purchase_order_id INTEGER NOT NULL REFERENCES purchase_orders(id)
);




ALTER TABLE user_roles ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE user_roles ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");


-- CREATE OR REPLACE FUNCTION generate_invoice_no()
CREATE OR REPLACE FUNCTION generate_invoice_no()
    RETURNS TRIGGER AS $$
DECLARE
    max_invoice_seq INTEGER;
    company_initials VARCHAR;
BEGIN
    -- Extract the maximum sequence number from existing invoice_no
    SELECT COALESCE(MAX(CAST(SPLIT_PART(SPLIT_PART(invoice_no, '-', 4), '-', 1) AS INTEGER)), 0) INTO max_invoice_seq FROM invoices;

    -- Get the company initials
    SELECT initials INTO company_initials FROM companies WHERE id = NEW.company_id;

    -- Generate the new invoice_no
    NEW.invoice_no := 'QPL-INV-' || TO_CHAR(current_date, 'YYYYMMDD') || '-' || LPAD((max_invoice_seq + 1)::TEXT, 4, '0') || '-' || company_initials;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER generate_invoice_no_trigger
    BEFORE INSERT ON invoices
    FOR EACH ROW
EXECUTE FUNCTION generate_invoice_no();




CREATE OR REPLACE FUNCTION generate_quotation_no()
    RETURNS TRIGGER AS $$
DECLARE
    max_quotation_seq INTEGER;
BEGIN
    -- Extract the maximum sequence number from existing quotation_no
    SELECT COALESCE(MAX(CAST(SPLIT_PART(SPLIT_PART(quotation_no, '-', 4), '-', 1) AS INTEGER)), 0) INTO max_quotation_seq FROM quotations;

    -- Generate the new quotation_no
    NEW.quotation_no := 'QPL-QTT-' || TO_CHAR(current_date, 'YYYYMMDD') || '-' || LPAD((max_quotation_seq + 1)::TEXT, 4, '0');

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER generate_quotation_no_trigger
    BEFORE INSERT ON quotations
    FOR EACH ROW
EXECUTE FUNCTION generate_quotation_no();




CREATE OR REPLACE FUNCTION generate_purchase_order_no()
    RETURNS TRIGGER AS $$
DECLARE
    max_po_seq INTEGER;
BEGIN
    -- Extract the maximum sequence number from existing po_no
    SELECT COALESCE(MAX(CAST(SPLIT_PART(SPLIT_PART(po_no, '-', 4), '-', 1) AS INTEGER)), 0) INTO max_po_seq FROM purchase_orders;

    -- Generate the new po_no
    NEW.po_no := 'QPL-PO-' || TO_CHAR(current_date, 'YYYYMMDD') || '-' || LPAD((max_po_seq + 1)::TEXT, 4, '0');

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER generate_purchase_order_no_trigger
    BEFORE INSERT ON purchase_orders
    FOR EACH ROW
EXECUTE FUNCTION generate_purchase_order_no();




-- After update quotations
-- CREATE OR REPLACE FUNCTION increment_quotation_revision_number()
CREATE OR REPLACE FUNCTION increment_quotation_revision_number()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.quotation_revision_number := OLD.quotation_revision_number + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER increment_quotation_revision_number_trigger
    BEFORE UPDATE ON quotations
    FOR EACH ROW
EXECUTE FUNCTION increment_quotation_revision_number();


-- Payment request trigger
-- CREATE OR REPLACE FUNCTION generate_payment_request_no()
-- Create a sequence for the unique serial number part of the payment request number
CREATE SEQUENCE payment_request_number_seq AS INTEGER START 1;

-- Create the function to generate the payment request number
CREATE OR REPLACE FUNCTION generate_payment_request_no() RETURNS TRIGGER AS $$
DECLARE
    current_day_seq INTEGER;
BEGIN
    -- Get the current value of the sequence
    SELECT nextval('payment_request_number_seq') INTO current_day_seq;

    -- Generate the payment request number
    NEW.payment_request_no := 'QPL-PR-' || to_char(current_date, 'YYYYMMDD') || '-' || lpad(current_day_seq::text, 4, '0');

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER generate_payment_request_no_trigger
    BEFORE INSERT ON payment_requests
    FOR EACH ROW
EXECUTE FUNCTION generate_payment_request_no();


-- Create a sequence for the unique serial number part of the petty_cash_no
CREATE SEQUENCE petty_cash_number_seq AS INTEGER START 1;

-- Create the function to generate the petty_cash_no
CREATE OR REPLACE FUNCTION generate_petty_cash_no() RETURNS TRIGGER AS $$
DECLARE
    current_petty_cash_seq INTEGER;
BEGIN
    -- Get the current value of the sequence
    SELECT nextval('petty_cash_number_seq') INTO current_petty_cash_seq;

    -- Generate the petty_cash_no
    NEW.petty_cash_no := 'QPL-PC-' || to_char(NEW.transaction_date, 'YYYYMMDD') || '-' || lpad(current_petty_cash_seq::text, 4, '0');

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER generate_petty_cash_no_trigger
    BEFORE INSERT ON petty_cash
    FOR EACH ROW
EXECUTE FUNCTION generate_petty_cash_no();
