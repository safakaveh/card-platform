CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS "card-platform";

CREATE TABLE "card-platform".orders (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_name VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'done')),
    description TEXT,
    order_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_order_order_name UNIQUE (order_name)
);

CREATE TABLE "card-platform".cards (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid_order UUID NOT NULL
        REFERENCES "card-platform".orders(uuid) ON DELETE CASCADE,
    has_laser BOOLEAN NOT NULL DEFAULT FALSE,
    has_magnet BOOLEAN NOT NULL DEFAULT FALSE,
    has_mifare_uid BOOLEAN NOT NULL DEFAULT FALSE,
    has_mifare BOOLEAN NOT NULL DEFAULT FALSE,
    has_java_card BOOLEAN NOT NULL DEFAULT FALSE,
    has_press BOOLEAN NOT NULL DEFAULT FALSE,
    has_temperature BOOLEAN NOT NULL DEFAULT FALSE,
    is_done BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cards_order ON "card-platform".cards(uuid_order);

CREATE TABLE "card-platform".card_status_history (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid_card UUID NOT NULL REFERENCES "card-platform".cards(uuid) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_card_status_history_card ON "card-platform".card_status_history(uuid_card);

CREATE TABLE "card-platform".laser_data (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid_card UUID NOT NULL REFERENCES "card-platform".cards(uuid) ON DELETE CASCADE,
    side VARCHAR(20) NOT NULL CHECK (side IN ('front', 'back')),
    row_no INT NOT NULL,
    content_type VARCHAR(50) NOT NULL,
    content BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_laser_card_side_row UNIQUE (uuid_card, side, row_no)
);

CREATE INDEX idx_laser_data_card ON "card-platform".laser_data(uuid_card);

CREATE TABLE "card-platform".magnet_data (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid_card UUID NOT NULL REFERENCES "card-platform".cards(uuid) ON DELETE CASCADE,
    track_no INT NOT NULL CHECK (track_no BETWEEN 1 AND 3),
    content BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_magnet_card_track UNIQUE (uuid_card, track_no)
);

CREATE INDEX idx_magnet_data_card ON "card-platform".magnet_data(uuid_card);

CREATE TABLE "card-platform".mifare_data (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid_card UUID NOT NULL REFERENCES "card-platform".cards(uuid) ON DELETE CASCADE,
    block_no INT NOT NULL,
    key_a BYTEA NOT NULL,
    key_b BYTEA NOT NULL,
    content BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_mifare_card_block UNIQUE (uuid_card, block_no)
);

CREATE INDEX idx_mifare_data_card ON "card-platform".mifare_data(uuid_card);

-- تعریف تابع به‌روزرسانی فیلد updated_at
CREATE OR REPLACE FUNCTION "card-platform".update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_orders
BEFORE UPDATE ON "card-platform".orders
FOR EACH ROW
EXECUTE FUNCTION "card-platform".update_updated_at();

CREATE TRIGGER trigger_update_cards
BEFORE UPDATE ON "card-platform".cards
FOR EACH ROW
EXECUTE FUNCTION "card-platform".update_updated_at();
