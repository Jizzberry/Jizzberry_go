-- +migrate Up
CREATE TABLE IF NOT EXISTS 'actors' (
    'generated_id' integer primary key autoincrement,
    'name' varchar,
    'website' varchar,
    'urlid' varchar
);

-- +migrate Down
DROP TABLE 'actors';