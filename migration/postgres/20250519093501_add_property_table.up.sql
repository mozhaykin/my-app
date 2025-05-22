BEGIN;

CREATE TABLE IF NOT EXISTS property
(
    profile_id UUID REFERENCES profile (id) ON DELETE CASCADE,
    tags       TEXT[]
);

-- Добаить индекс на property (profile_id)
CREATE INDEX IF NOT EXISTS idx_property_profile_id ON property (profile_id);

COMMIT;