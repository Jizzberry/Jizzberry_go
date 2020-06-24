-- +migrate Up
CREATE TABLE IF NOT EXISTS 'studios' (
   'generated_id' integer primary key autoincrement,
   'studio' varchar,
   'count' integer not null

);

-- +migrate Down
DROP TABLE 'studios';