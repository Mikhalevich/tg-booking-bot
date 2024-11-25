-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE role_name AS ENUM(
    'owner',
    'manager',
    'employee'
);

CREATE TABLE role(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name role_name NOT NULL,
    CONSTRAINT role_name_unique UNIQUE(name)
);

INSERT INTO role(name) VALUES
    ('owner'), -- 1
    ('manager'), -- 2
    ('employee'); -- 3

CREATE TABLE employee(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    first_name TEXT NOT NULL,
    second_name TEXT NOT NULL,
    role_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT employee_role_fk FOREIGN KEY(role_id) REFERENCES role(id) ON DELETE RESTRICT ON UPDATE CASCADE
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

CREATE TABLE role_inher(
    role_id INTEGER NOT NULL,
    parent_id INTEGER NOT NULL,

    CONSTRAINT role_inher_pk PRIMARY KEY(role_id, parent_id),
    CONSTRAINT role_inher_fk_by_role FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT role_inher_fk_by_parent FOREIGN KEY (parent_id) REFERENCES role(id) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO role_inher(role_id, parent_id) VALUES
    (2, 3),
    (1, 2);

CREATE TYPE permission_action AS ENUM(
    'view_schedule_template',
    'add_employee',
    'edit_employee_schedule'
);

CREATE TABLE permission(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    action permission_action NOT NULL,

    CONSTRAINT permission_action_unique UNIQUE(action)
);

CREATE TABLE role_perm(
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,

    CONSTRAINT role_perm_pk PRIMARY KEY(role_id, permission_id),
    CONSTRAINT role_perm_role_id_fk FOREIGN KEY(role_id) REFERENCES role(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT role_perm_permission_id_fk FOREIGN KEY(permission_id) REFERENCES permission(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE role_perm;
DROP TABLE permission;
DROP TABLE role_inher;
DROP TABLE employee;
DROP TABLE role;
