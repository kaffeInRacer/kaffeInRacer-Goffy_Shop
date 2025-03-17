-- CREATE TABLE IF NOT EXISTS
CREATE TABLE IF NOT EXISTS goffy_products (
    id BIGSERIAL PRIMARY KEY,
    qty BIGINT NOT NULL CHECK (qty >= 0),
    slug VARCHAR(255) NOT NULL CHECK (slug <> ''),
    name VARCHAR(255) NOT NULL CHECK (name <> ''),
    description TEXT NOT NULL CHECK (description <> ''),
    weight BIGINT NOT NULL CHECK ( weight >= 0 ),
    price BIGINT NOT NULL CHECK (price >= 0),
    is_delete BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- CREATE TRIGGER FOR goffy_products
CREATE TRIGGER update_time_stamp_products_trigger
    BEFORE UPDATE ON goffy_products
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

