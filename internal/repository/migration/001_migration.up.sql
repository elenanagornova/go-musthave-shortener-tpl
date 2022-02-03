-- +goose Up
CREATE SCHEMA IF NOT EXISTS shortener;
-- DROP SCHEMA shortener CASCADE ;
   CREATE SCHEMA shortener;
SET SEARCH_PATH TO shortener;

CREATE TABLE IF NOT EXISTS links
(
    id           serial primary key,
    short_link      varchar,
    original_link varchar,
    user_uid          varchar,
    created_at   TIMESTAMP
);
ALTER TABLE links
    ALTER COLUMN created_at SET DEFAULT now();
CREATE UNIQUE INDEX original_link_idx ON links USING btree (original_link);
CREATE UNIQUE INDEX short_link_idx ON links USING btree (short_link);