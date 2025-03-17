-- CREATE TABLE IF NOT EXISTS
CREATE TABLE IF NOT EXISTS goffy_shipments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES goffy_orders(id) ON DELETE CASCADE,
    courier VARCHAR(255) NOT NULL CHECK (courier <> ''),
    tracking_number VARCHAR(255) UNIQUE NOT NULL CHECK (tracking_number <> ''),
    shipping_status VARCHAR(255) NOT NULL CHECK (shipping_status <> ''),
    shipped_at TIMESTAMP DEFAULT NULL,
    delivered_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- CREATE TRIGGER FOR goffy_order_details
CREATE TRIGGER update_time_stamp_goffy_shipments_trigger
    BEFORE UPDATE ON goffy_shipments
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();