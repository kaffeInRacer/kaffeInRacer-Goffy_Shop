-- CREATE TABLE IF NOT EXISTS
CREATE TABLE IF NOT EXISTS goffy_order_details (
   id BIGSERIAL PRIMARY KEY,
   order_id BIGINT NOT NULL REFERENCES goffy_orders(id) ON DELETE RESTRICT,
   product_id BIGINT NOT NULL REFERENCES goffy_products(id) ON DELETE RESTRICT,
   qty BIGINT NOT NULL CHECK (qty >= 0),
   price BIGINT NOT NULL CHECK (price >= 0),
   weight_total BIGINT NOT NULL CHECK (weight_total >= 0),
   grand_total BIGINT GENERATED ALWAYS AS (price * qty) STORED,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TRIGGER FOR goffy_order_details
CREATE TRIGGER update_time_stamp_order_details_trigger
    BEFORE UPDATE ON goffy_order_details
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- CREATE FUNCTION CALCULATE WEIGHT TOTAL
CREATE OR REPLACE FUNCTION calculate_weight_total()
    RETURNS TRIGGER AS $$
BEGIN
    SELECT weight INTO NEW.weight_total
    FROM goffy_products
    WHERE id = NEW.product_id;
    NEW.weight_total = NEW.weight_total * NEW.qty;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- CREATE TRIGGER INSERT CALCULATE WEIGHT TOTAL AFTER INSERT AND UPDATE
CREATE TRIGGER calculate_weight_total_trigger
    BEFORE INSERT OR UPDATE ON goffy_order_details
    FOR EACH ROW
    EXECUTE FUNCTION calculate_weight_total();