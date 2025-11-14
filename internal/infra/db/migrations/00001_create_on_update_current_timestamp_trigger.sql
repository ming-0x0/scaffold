-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION on_update_current_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
    t TEXT;
BEGIN
    FOR t IN
        SELECT table_name
        FROM information_schema.columns
        WHERE column_name = 'updated_at'
        AND table_schema = 'public'
    LOOP
        EXECUTE format('
            DROP TRIGGER IF EXISTS trg_on_update_current_timestamp_%I ON %I;
            CREATE TRIGGER trg_on_update_current_timestamp_%I
            BEFORE UPDATE ON %I
            FOR EACH ROW EXECUTE FUNCTION on_update_current_timestamp();
        ', t, t, t, t);
    END LOOP;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
DECLARE
    t TEXT;
BEGIN
    FOR t IN
        SELECT table_name
        FROM information_schema.columns
        WHERE column_name = 'updated_at'
        AND table_schema = 'public'
    LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS trg_on_update_current_timestamp_%I ON %I;', t, t);
    END LOOP;
END $$;

DROP FUNCTION IF EXISTS on_update_current_timestamp();
-- +goose StatementEnd
