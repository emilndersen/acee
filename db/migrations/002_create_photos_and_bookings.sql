-- migrations/002_create_photos_and_bookings.sql

CREATE TABLE photos (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    title       TEXT        NOT NULL DEFAULT '',
    album       TEXT        NOT NULL DEFAULT 'general',
    description TEXT        NOT NULL DEFAULT '',
    url         TEXT        NOT NULL,
    sort_order  INT         NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_photos_album ON photos(album);

-- ──────────────────────────────────────────

CREATE TABLE bookings (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL,
    contact     TEXT        NOT NULL,
    shoot_type  TEXT        NOT NULL DEFAULT '',
    date        TEXT        NOT NULL DEFAULT '',
    idea        TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
