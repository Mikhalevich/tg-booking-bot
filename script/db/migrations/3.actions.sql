-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE actions(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    employee_id INTEGER NOT NULL,
    action TEXT NOT NULL,
    payload JSONB NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,

    CONSTRAINT actions_employee_fk FOREIGN KEY(employee_id) REFERENCES employee(id)
);

CREATE INDEX actions_employee_id_is_completed_idx ON actions(employee_id, is_completed);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE actions;
DROP INDEX actions_employee_id_is_completed_idx;
