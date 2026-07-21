PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS card_platform_orders (
    uuid TEXT PRIMARY KEY NOT NULL,
    order_name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'in_progress', 'done')),
    description TEXT,
    order_date INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    CONSTRAINT uq_card_platform_orders_order_name UNIQUE (order_name)
);

CREATE TABLE IF NOT EXISTS card_platform_cards (
    uuid TEXT PRIMARY KEY NOT NULL,
    uuid_order TEXT NOT NULL,
    has_laser INTEGER NOT NULL DEFAULT 0 CHECK (has_laser IN (0, 1)),
    has_magnet INTEGER NOT NULL DEFAULT 0 CHECK (has_magnet IN (0, 1)),
    has_mifare_uid INTEGER NOT NULL DEFAULT 0 CHECK (has_mifare_uid IN (0, 1)),
    has_mifare INTEGER NOT NULL DEFAULT 0 CHECK (has_mifare IN (0, 1)),
    has_java_card INTEGER NOT NULL DEFAULT 0 CHECK (has_java_card IN (0, 1)),
    has_press INTEGER NOT NULL DEFAULT 0 CHECK (has_press IN (0, 1)),
    has_temperature INTEGER NOT NULL DEFAULT 0 CHECK (has_temperature IN (0, 1)),
    is_done INTEGER NOT NULL DEFAULT 0 CHECK (is_done IN (0, 1)),
    description TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    CONSTRAINT fk_card_platform_cards_order
        FOREIGN KEY (uuid_order)
        REFERENCES card_platform_orders(uuid)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_card_platform_cards_uuid_order
    ON card_platform_cards(uuid_order);

CREATE TABLE IF NOT EXISTS card_platform_card_status_history (
    uuid TEXT PRIMARY KEY NOT NULL,
    uuid_card TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    CONSTRAINT fk_card_platform_card_status_history_card
        FOREIGN KEY (uuid_card)
        REFERENCES card_platform_cards(uuid)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_card_platform_card_status_history_uuid_card
    ON card_platform_card_status_history(uuid_card);

CREATE TABLE IF NOT EXISTS card_platform_laser_data (
    uuid TEXT PRIMARY KEY NOT NULL,
    uuid_card TEXT NOT NULL,
    side TEXT NOT NULL CHECK (side IN ('front', 'back')),
    row_no INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    content BLOB NOT NULL,
    created_at INTEGER NOT NULL,
    CONSTRAINT uq_card_platform_laser_data_card_side_row
        UNIQUE (uuid_card, side, row_no),
    CONSTRAINT fk_card_platform_laser_data_card
        FOREIGN KEY (uuid_card)
        REFERENCES card_platform_cards(uuid)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_card_platform_laser_data_uuid_card
    ON card_platform_laser_data(uuid_card);

CREATE TABLE IF NOT EXISTS card_platform_magnet_data (
    uuid TEXT PRIMARY KEY NOT NULL,
    uuid_card TEXT NOT NULL,
    track_no INTEGER NOT NULL CHECK (track_no BETWEEN 1 AND 3),
    content BLOB NOT NULL,
    created_at INTEGER NOT NULL,
    CONSTRAINT uq_card_platform_magnet_data_card_track
        UNIQUE (uuid_card, track_no),
    CONSTRAINT fk_card_platform_magnet_data_card
        FOREIGN KEY (uuid_card)
        REFERENCES card_platform_cards(uuid)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_card_platform_magnet_data_uuid_card
    ON card_platform_magnet_data(uuid_card);

CREATE TABLE IF NOT EXISTS card_platform_mifare_data (
    uuid TEXT PRIMARY KEY NOT NULL,
    uuid_card TEXT NOT NULL,
    block_no INTEGER NOT NULL,
    key_a BLOB NOT NULL,
    key_b BLOB NOT NULL,
    content BLOB NOT NULL,
    created_at INTEGER NOT NULL,
    CONSTRAINT uq_card_platform_mifare_data_card_block
        UNIQUE (uuid_card, block_no),
    CONSTRAINT fk_card_platform_mifare_data_card
        FOREIGN KEY (uuid_card)
        REFERENCES card_platform_cards(uuid)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_card_platform_mifare_data_uuid_card
    ON card_platform_mifare_data(uuid_card);

CREATE TRIGGER IF NOT EXISTS trg_card_platform_orders_updated_at
AFTER UPDATE ON card_platform_orders
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE card_platform_orders
    SET updated_at = CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER)
    WHERE uuid = OLD.uuid;
END;

CREATE TRIGGER IF NOT EXISTS trg_card_platform_cards_updated_at
AFTER UPDATE ON card_platform_cards
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE card_platform_cards
    SET updated_at = CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER)
    WHERE uuid = OLD.uuid;
END;
