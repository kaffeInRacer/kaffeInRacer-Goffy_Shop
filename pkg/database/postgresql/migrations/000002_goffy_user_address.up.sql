-- CREATE TABLE IF NOT EXISTS
CREATE TABLE IF NOT EXISTS goffy_address_users (
    id BIGSERIAL PRIMARY KEY,
    users_id BIGINT NOT NULL REFERENCES goffy_users(id) ON DELETE CASCADE,
    provincy VARCHAR(255) NOT NULL,
    city_name VARCHAR(255) NOT NULL,
    district_name VARCHAR(255),
    zip_code VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

-- CREATE TRIGGER FOR goffy_address_users
CREATE TRIGGER update_time_stamp_users_address_trigger
    BEFORE UPDATE ON goffy_address_users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();


