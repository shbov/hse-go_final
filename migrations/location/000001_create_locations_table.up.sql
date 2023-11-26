CREATE TABLE IF NOT EXISTS locations
(
    id   UUID PRIMARY KEY,
    name VARCHAR(255),
    auto VARCHAR(255),
    lat  NUMERIC,
    lng  NUMERIC
);