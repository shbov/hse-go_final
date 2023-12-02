CREATE TABLE IF NOT EXISTS locations
(
    id         SERIAL PRIMARY KEY,
    driver_id  UUID    NOT NULL,

    lat        NUMERIC NOT NULL,
    lng        NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX driver_index ON locations (driver_id);