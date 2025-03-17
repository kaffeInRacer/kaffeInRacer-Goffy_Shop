-- Hapus trigger
DROP TRIGGER IF EXISTS update_time_stamp_order_details_trigger ON goffy_order_details;

--
DROP TRIGGER IF EXISTS calculate_weight_total_trigger ON goffy_order_details;

--
    DO $$
        BEGIN
            IF NOT EXISTS(
                SELECT 1 FROM pg_trigger WHERE tgname = 'calculate_weight_total_trigger'
            ) THEN
                DROP FUNCTION IF EXISTS calculate_weight_total();
            END if;
        END $$;


-- Hapus tabel
DROP TABLE IF EXISTS goffy_order_details;