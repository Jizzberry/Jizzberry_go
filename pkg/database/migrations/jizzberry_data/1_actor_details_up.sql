-- +migrate Up
CREATE TABLE IF NOT EXISTS 'actor_details'(
                                              'generated_id' integer primary key autoincrement,
                                              'thumbnail'    varchar,
                                              'actor_id'     integer,
                                              'name'         varchar,
                                              'born'         varchar,
                                              'birthplace'   varchar,
                                              'height'       varchar,
                                              'weight'       varchar
);

-- +migrate Down
DROP TABLE 'actor_details';