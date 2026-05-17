-- +goose Up
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    user_id UUID NOT NULL,
    start_date VARCHAR(7) NOT NULL, -- MM-YYYY format
    end_date VARCHAR(7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_id ON subscriptions(user_id);
CREATE INDEX idx_start_date ON subscriptions(start_date);
CREATE INDEX idx_service_name ON subscriptions(service_name);

-- +goose Down
DROP INDEX IF EXISTS idx_service_name;
DROP INDEX IF EXISTS idx_start_date;
DROP INDEX IF EXISTS idx_user_id;
DROP TABLE IF EXISTS subscriptions;
