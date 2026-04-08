ALTER TABLE stocks DROP COLUMN IF EXISTS data_source;
ALTER TABLE stocks DROP COLUMN IF EXISTS last_polygon_update;
DROP TABLE IF EXISTS market_status;
