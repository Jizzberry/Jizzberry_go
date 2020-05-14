-- +migrate Up
CREATE TABLE IF NOT EXISTS 'auth' (
    'username' varchar primary key,
    'password' varchar
);

-- +migrate Down
DROP TABLE 'auth';

