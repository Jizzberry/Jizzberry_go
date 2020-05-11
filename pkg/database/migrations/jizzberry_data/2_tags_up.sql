-- +migrate Up
CREATE TABLE IF NOT EXISTS 'tags' (
   'generated_id' integer primary key autoincrement,
   'tag' varchar
);

-- +migrate Down
DROP TABLE 'tags';