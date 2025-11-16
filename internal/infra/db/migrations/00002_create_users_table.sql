-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id              CHAR(26) NOT NULL PRIMARY KEY,
    email           VARCHAR(255) NOT NULL UNIQUE,
    username        VARCHAR(50) NOT NULL UNIQUE,
    full_name       VARCHAR(50) NOT NULL,
    status          SMALLINT NOT NULL DEFAULT 1,
    is_admin        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMP(0) NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP(0) NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_on_update_current_timestamp_users
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION on_update_current_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_on_update_current_timestamp_users ON users;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
