-- Hapus trigger
DROP TRIGGER IF EXISTS update_time_stamp_users_trigger ON goffy_users;

-- Hapus function (jika tidak ada tabel lain yang menggunakannya)
DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM pg_trigger WHERE tgname = 'update_time_stamp_users_trigger'
        ) THEN
            DROP FUNCTION IF EXISTS set_updated_at();
        END IF;
    END $$;

-- Hapus tabel
DROP TABLE IF EXISTS goffy_users;