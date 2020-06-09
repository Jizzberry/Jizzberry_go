-- +migrate Up
CREATE TABLE IF NOT EXISTS 'tags' (
   'generated_id' integer primary key autoincrement,
   'tag' varchar,
   'count' integer not null
);

-- +migrate Down
DROP TABLE 'tags';