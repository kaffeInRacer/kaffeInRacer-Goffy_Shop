-- CREATE TABLE IF NOT EXISTS
CREATE TABLE IF NOT EXISTS goffy_orders (
    id BIGSERIAL PRIMARY KEY,
    order_code UUID DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES goffy_users(id) ON DELETE RESTRICT,
    paid_bank VARCHAR(255) NOT NULL CHECK (paid_bank <> ''),
    paid_status VARCHAR(255) NOT NULL CHECK (paid_status <> ''),
    paid_at TIMESTAMP DEFAULT NULL,
    grand_total BIGINT NOT NULL CHECK (grand_total >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TRIGGER FOR goffy_orders
CREATE TRIGGER update_time_stamp_orders_trigger
    BEFORE UPDATE ON goffy_orders
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
