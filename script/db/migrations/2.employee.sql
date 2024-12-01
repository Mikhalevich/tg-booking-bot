-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE role_name AS ENUM(
    'unspecified',
    'owner',
    'manager',
    'employee'
);

CREATE TABLE role(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name role_name NOT NULL,
    CONSTRAINT role_name_unique UNIQUE(name)
);

CREATE TYPE employee_state AS ENUM (
    'verification_required',
    'registered'
);

CREATE TABLE employee(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    role_id INTEGER NOT NULL,
    chat_id BIGINT,
    state employee_state NOT NULL DEFAULT 'verification_required',
    verification_code TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT employee_role_fk FOREIGN KEY(role_id) REFERENCES role(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT chat_id_unique UNIQUE(chat_id),
    CONSTRAINT verification_code_unique UNIQUE(verification_code)
);

CREATE TABLE employee_schedule(
    employee_id INTEGER NOT NULL,
    schedule JSONB NOT NULL,
    schedule_version INTEGER NOT NULL DEFAULT 0,
    schedule_schema_version INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT employee_schedule_pk PRIMARY KEY(employee_id),
    CONSTRAINT employee_schedule_employee_fk FOREIGN KEY(employee_id) REFERENCES employee(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE employee_schedule;
DROP TABLE employee;
DROP TYPE employee_state
DROP TABLE role;
DROP TYPE role_name;
