BEGIN;

create table if not exists road_type
(
    character_id integer primary key,
    road_type integer
);

END