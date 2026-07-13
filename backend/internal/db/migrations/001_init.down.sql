DROP TRIGGER IF EXISTS trigger_update_cards ON "card-platform".cards;
DROP TRIGGER IF EXISTS trigger_update_orders ON "card-platform".orders;

DROP FUNCTION IF EXISTS "card-platform".update_updated_at();

DROP TABLE IF EXISTS "card-platform".mifare_data;
DROP TABLE IF EXISTS "card-platform".magnet_data;
DROP TABLE IF EXISTS "card-platform".laser_data;
DROP TABLE IF EXISTS "card-platform".card_status_history;
DROP TABLE IF EXISTS "card-platform".cards;
DROP TABLE IF EXISTS "card-platform".orders;

DROP SCHEMA IF EXISTS "card-platform";

DROP EXTENSION IF EXISTS pgcrypto;
