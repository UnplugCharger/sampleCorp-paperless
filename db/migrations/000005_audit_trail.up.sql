CREATE TABLE audit_petty_cash (
audit_id bigserial PRIMARY KEY,
operation varchar(6) CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
transaction_id INTEGER,
employee_id INTEGER,
amount VARCHAR,
description VARCHAR(255),
transaction_date DATE,
updated_at DATE,
folio VARCHAR(50),
debit_account VARCHAR(50),
status VARCHAR(20),
authorised_by INTEGER,
changed_at timestamptz NOT NULL DEFAULT (now()),
changed_by INTEGER
);

-- Trigger function
CREATE OR REPLACE FUNCTION audit_petty_cash()
    RETURNS TRIGGER AS $audit_petty_cash$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO audit_petty_cash (operation, transaction_id, employee_id, amount, description, transaction_date, updated_at, folio, debit_account, status, authorised_by, changed_at, changed_by)
        VALUES ('DELETE', OLD.transaction_id, OLD.employee_id, OLD.amount, OLD.description, OLD.transaction_date, OLD.updated_at, OLD.folio, OLD.debit_account, OLD.status, OLD.authorised_by, now(), current_setting('audit.current_user_id')::integer);
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO audit_petty_cash (operation, transaction_id, employee_id, amount, description, transaction_date, updated_at, folio, debit_account, status, authorised_by, changed_at, changed_by)
        VALUES ('UPDATE', NEW.transaction_id, NEW.employee_id, NEW.amount, NEW.description, NEW.transaction_date, NEW.updated_at, NEW.folio, NEW.debit_account, NEW.status, NEW.authorised_by, now(), current_setting('audit.current_user_id')::integer);
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO audit_petty_cash (operation, transaction_id, employee_id, amount, description, transaction_date, updated_at, folio, debit_account, status, authorised_by, changed_at, changed_by)
        VALUES ('INSERT', NEW.transaction_id, NEW.employee_id, NEW.amount, NEW.description, NEW.transaction_date, NEW.updated_at, NEW.folio, NEW.debit_account, NEW.status, NEW.authorised_by, now(), current_setting('audit.current_user_id')::integer);
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$audit_petty_cash$ LANGUAGE plpgsql;

-- Trigger
CREATE TRIGGER audit_petty_cash_trigger
    AFTER INSERT OR UPDATE OR DELETE ON petty_cash
    FOR EACH ROW EXECUTE PROCEDURE audit_petty_cash();

-- Payment Requests
CREATE TABLE audit_payment_requests (
audit_id bigserial PRIMARY KEY,
operation varchar(6) CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
request_id INTEGER,
employee_id INTEGER,
amount FLOAT,
description VARCHAR(255),
request_date DATE,
status VARCHAR(20),
admin_id INTEGER,
approval_date DATE,
invoice_id INTEGER,
changed_at timestamptz NOT NULL DEFAULT (now()),
changed_by INTEGER
);

-- Trigger function
CREATE OR REPLACE FUNCTION audit_payment_requests()
    RETURNS TRIGGER AS $audit_payment_requests$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO audit_payment_requests (operation, request_id, employee_id, amount, description, request_date, status, admin_id, approval_date, invoice_id, changed_at, changed_by)
        VALUES ('DELETE', OLD.request_id, OLD.employee_id, OLD.amount, OLD.description, OLD.request_date, OLD.status, OLD.admin_id, OLD.approval_date, OLD.invoice_id, now(), current_setting('audit.current_user_id')::integer);
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO audit_payment_requests (operation, request_id, employee_id, amount, description, request_date, status, admin_id, approval_date, invoice_id, changed_at, changed_by)
        VALUES ('UPDATE', NEW.request_id, NEW.employee_id, NEW.amount, NEW.description, NEW.request_date, NEW.status, NEW.admin_id, NEW.approval_date, NEW.invoice_id, now(), current_setting('audit.current_user_id')::integer);
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO audit_payment_requests (operation, request_id, employee_id, amount, description, request_date, status, admin_id, approval_date, invoice_id, changed_at, changed_by)
        VALUES ('INSERT', NEW.request_id, NEW.employee_id, NEW.amount, NEW.description, NEW.request_date, NEW.status, NEW.admin_id, NEW.approval_date, NEW.invoice_id, now(), current_setting('audit.current_user_id')::integer);
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$audit_payment_requests$ LANGUAGE plpgsql;

CREATE TRIGGER audit_payment_requests_trigger
    AFTER INSERT OR UPDATE OR DELETE ON payment_requests
    FOR EACH ROW EXECUTE PROCEDURE audit_payment_requests();


-- Audit Table for Quotations
CREATE TABLE audit_quotations (
audit_id bigserial PRIMARY KEY,
operation varchar(6) CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
id INTEGER,
quotation_no VARCHAR(50),
quotation_revision INTEGER,
quotation_revision_number INTEGER,
date DATE,
attn VARCHAR(100),
company_id INTEGER,
site VARCHAR(100),
validity INTEGER,
warranty INTEGER,
payment_terms TEXT,
delivery_terms TEXT,
signatory_id INTEGER,
status VARCHAR(50),
sent_or_received VARCHAR(10),
changed_at timestamptz NOT NULL DEFAULT (now()),
changed_by INTEGER
);

-- Trigger function for quotations
CREATE OR REPLACE FUNCTION audit_quotations_trigger_function()
    RETURNS TRIGGER AS $audit_quotations$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO audit_quotations (operation, id, quotation_no, quotation_revision, quotation_revision_number, date, attn, company_id, site, validity, warranty, payment_terms, delivery_terms, signatory_id, status, sent_or_received, changed_at, changed_by)
        VALUES ('DELETE', OLD.id, OLD.quotation_no, OLD.quotation_revision, OLD.quotation_revision_number, OLD.date, OLD.attn, OLD.company_id, OLD.site, OLD.validity, OLD.warranty, OLD.payment_terms, OLD.delivery_terms, OLD.signatory_id, OLD.status, OLD.sent_or_received, now(), current_setting('audit.current_user_id')::integer);
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO audit_quotations (operation, id, quotation_no, quotation_revision, quotation_revision_number, date, attn, company_id, site, validity, warranty, payment_terms, delivery_terms, signatory_id, status, sent_or_received, changed_at, changed_by)
        VALUES ('UPDATE', NEW.id, NEW.quotation_no, NEW.quotation_revision, NEW.quotation_revision_number, NEW.date, NEW.attn, NEW.company_id, NEW.site, NEW.validity, NEW.warranty, NEW.payment_terms, NEW.delivery_terms, NEW.signatory_id, NEW.status, NEW.sent_or_received, now(), current_setting('audit.current_user_id')::integer);
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO audit_quotations (operation, id, quotation_no, quotation_revision, quotation_revision_number, date, attn, company_id, site, validity, warranty, payment_terms, delivery_terms, signatory_id, status, sent_or_received, changed_at, changed_by)
        VALUES ('INSERT', NEW.id, NEW.quotation_no, NEW.quotation_revision, NEW.quotation_revision_number, NEW.date, NEW.attn, NEW.company_id, NEW.site, NEW.validity, NEW.warranty, NEW.payment_terms, NEW.delivery_terms, NEW.signatory_id, NEW.status, NEW.sent_or_received, now(), current_setting('audit.current_user_id')::integer);
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$audit_quotations$ LANGUAGE plpgsql;

-- Trigger for quotations
CREATE TRIGGER audit_quotations_trigger
    AFTER INSERT OR UPDATE OR DELETE ON quotations
    FOR EACH ROW EXECUTE PROCEDURE audit_quotations_trigger_function();


-- Audit Table for Purchase Orders
CREATE TABLE audit_purchase_orders (
audit_id bigserial PRIMARY KEY,
operation varchar(6) CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
id INTEGER,
po_no VARCHAR(50),
date DATE,
attn VARCHAR(100),
company_id INTEGER,
address VARCHAR(255),
signatory_id INTEGER,
quotation_id INTEGER,
po_status VARCHAR(50),
sent_or_received VARCHAR(10),
changed_at timestamptz NOT NULL DEFAULT (now()),
changed_by INTEGER
);

-- Trigger function for purchase_orders
CREATE OR REPLACE FUNCTION audit_purchase_orders_trigger_function()
    RETURNS TRIGGER AS $audit_purchase_orders$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO audit_purchase_orders (operation, id, po_no, date, attn, company_id, address, signatory_id, quotation_id, po_status, sent_or_received, changed_at, changed_by)
        VALUES ('DELETE', OLD.id, OLD.po_no, OLD.date, OLD.attn, OLD.company_id, OLD.address, OLD.signatory_id, OLD.quotation_id, OLD.po_status, OLD.sent_or_received, now(), current_setting('audit.current_user_id')::integer);
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO audit_purchase_orders (operation, id, po_no, date, attn, company_id, address, signatory_id, quotation_id, po_status, sent_or_received, changed_at, changed_by)
        VALUES ('UPDATE', NEW.id, NEW.po_no, NEW.date, NEW.attn, NEW.company_id, NEW.address, NEW.signatory_id, NEW.quotation_id, NEW.po_status, NEW.sent_or_received, now(), current_setting('audit.current_user_id')::integer);
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO audit_purchase_orders (operation, id, po_no, date, attn, company_id, address, signatory_id, quotation_id, po_status, sent_or_received, changed_at, changed_by)
        VALUES ('INSERT', NEW.id, NEW.po_no, NEW.date, NEW.attn, NEW.company_id, NEW.address, NEW.signatory_id, NEW.quotation_id, NEW.po_status, NEW.sent_or_received, now(), current_setting('audit.current_user_id')::integer);
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$audit_purchase_orders$ LANGUAGE plpgsql;

-- Trigger for purchase_orders
CREATE TRIGGER audit_purchase_orders_trigger
    AFTER INSERT OR UPDATE OR DELETE ON purchase_orders
    FOR EACH ROW EXECUTE PROCEDURE audit_purchase_orders_trigger_function();

-- Audit Table for Invoices
CREATE TABLE audit_invoices (
audit_id bigserial PRIMARY KEY,
operation varchar(6) CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
id INTEGER,
invoice_no VARCHAR(50),
purchase_order_number VARCHAR(50),
date DATE,
attn VARCHAR(100),
company_id INTEGER,
site VARCHAR(100),
amount_due FLOAT,
bank_details INTEGER,
signatory_id INTEGER,
sent_or_received VARCHAR(10),
changed_at timestamptz NOT NULL DEFAULT (now()),
changed_by INTEGER
);

-- Trigger function for invoices
CREATE OR REPLACE FUNCTION audit_invoices_trigger_function()
    RETURNS TRIGGER AS $audit_invoices$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        INSERT INTO audit_invoices (operation, id, invoice_no, purchase_order_number, date, attn, company_id, site, amount_due, bank_details, signatory_id, sent_or_received, changed_at, changed_by)
        VALUES ('DELETE', OLD.id, OLD.invoice_no, OLD.purchase_order_number, OLD.date, OLD.attn, OLD.company_id, OLD.site, OLD.amount_due, OLD.bank_details, OLD.signatory_id, OLD.sent_or_received, now(), current_setting('audit.current_user_id')::integer);
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO audit_invoices (operation, id, invoice_no, purchase_order_number, date, attn, company_id, site, amount_due, bank_details, signatory_id, sent_or_received, changed_at, changed_by)
        VALUES ('UPDATE', NEW.id, NEW.invoice_no, NEW.purchase_order_number, NEW.date, NEW.attn, NEW.company_id, NEW.site, NEW.amount_due, NEW.bank_details, NEW.signatory_id, NEW.sent_or_received, now(), current_setting('audit.current_user_id')::integer);
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO audit_invoices (operation, id, invoice_no, purchase_order_number, date, attn, company_id, site, amount_due, bank_details, signatory_id, sent_or_received, changed_at, changed_by)
        VALUES ('INSERT', NEW.id, NEW.invoice_no, NEW.purchase_order_number, NEW.date, NEW.attn, NEW.company_id, NEW.site, NEW.amount_due, NEW.bank_details, NEW.signatory_id, NEW.sent_or_received, now(), current_setting('audit.current_user_id')::integer);
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$audit_invoices$ LANGUAGE plpgsql;

-- Trigger for invoices
CREATE TRIGGER audit_invoices_trigger
    AFTER INSERT OR UPDATE OR DELETE ON invoices
    FOR EACH ROW EXECUTE PROCEDURE audit_invoices_trigger_function();
