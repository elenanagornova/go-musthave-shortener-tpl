CREATE SCHEMA IF NOT EXISTS shortener;
-- DROP SCHEMA shortener CASCADE ;
--CREATE SCHEMA shortener;
SET SEARCH_PATH TO shortener;

CREATE TABLE IF NOT EXISTS links
(
    id            serial primary key,
    short_link    varchar,
    original_link varchar,
    user_uid      varchar,
    removed       boolean
);
CREATE UNIQUE INDEX IF NOT EXISTS original_link_idx ON links USING btree (original_link);
CREATE UNIQUE INDEX IF NOT EXISTS short_link_idx ON links USING btree (short_link);