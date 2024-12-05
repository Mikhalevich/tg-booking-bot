-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE action_state AS ENUM(
    'in_progress',
    'completed',
    'canceled'
);

CREATE TABLE actions(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    employee_id INTEGER NOT NULL,
    action TEXT NOT NULL,
    payload JSONB NOT NULL,
    state action_state NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,

    CONSTRAINT actions_employee_fk FOREIGN KEY(employee_id) REFERENCES employee(id)
);

CREATE INDEX actions_employee_id_is_completed_idx ON actions(employee_id, state);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX actions_employee_id_is_completed_idx;
DROP TABLE actions;
DROP TYPE action_state;
