DROP TRIGGER IF EXISTS trg_cards_updated_at;
DROP TRIGGER IF EXISTS trg_orders_updated_at;

DROP INDEX IF EXISTS idx_mifare_data_uuid_card;
DROP INDEX IF EXISTS idx_magnet_data_uuid_card;
DROP INDEX IF EXISTS idx_laser_data_uuid_card;
DROP INDEX IF EXISTS idx_card_status_history_uuid_card;
DROP INDEX IF EXISTS idx_cards_uuid_order;

DROP TABLE IF EXISTS mifare_data;
DROP TABLE IF EXISTS magnet_data;
DROP TABLE IF EXISTS laser_data;
DROP TABLE IF EXISTS card_status_history;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS orders;

-- The migration bookkeeping table is intentionally removed as the final step
-- so a subsequent startup can recreate the complete schema from scratch.
DROP TABLE IF EXISTS schema_migrations;
