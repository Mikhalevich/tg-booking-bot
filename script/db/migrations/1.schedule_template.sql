-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE schedule_template(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    schedule_payload JSONB NOT NULL,
    payload_version INTEGER NOT NULL DEFAULT 0,
    payload_schema_version INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE schedule_template;
