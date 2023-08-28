CREATE OR REPLACE FUNCTION check_approval_status()
    RETURNS TRIGGER AS $$
BEGIN
    IF lower(NEW.status) = 'approved' AND lower(OLD.status) = 'approved' THEN
        RAISE EXCEPTION 'Payment request is already approved.';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER payment_request_update BEFORE UPDATE
    ON payment_requests FOR EACH ROW EXECUTE PROCEDURE check_approval_status();
