-- +migrate Up
CREATE TABLE IF NOT EXISTS 'studios' (
                                         'generated_id' integer primary key autoincrement,
                                         'studio' varchar
);

-- +migrate Down
DROP TABLE 'studios';