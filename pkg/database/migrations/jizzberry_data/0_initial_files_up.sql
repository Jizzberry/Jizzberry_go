-- +migrate Up
CREATE TABLE IF NOT EXISTS 'files' (
    'generated_id' integer primary key autoincrement,
    'file_name' varchar,
    'file_path' varchar,
    'date_created' varchar,
    'file_size' varchar,
    'length' varchar,
    'thumbnail' varchar,
    'symlinks' varchar,
    'tags' varchar,
    'studios' varchar,
    'actors' varchar,
    'url' varchar
);

-- +migrate Down
DROP TABLE 'files';