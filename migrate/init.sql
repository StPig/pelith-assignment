CREATE TABLE users
(
    id serial PRIMARY KEY,
    address varchar(100) UNIQUE NOT NULL,
    points numeric(100, 1) DEFAULT 0,
    onboarding_completed boolean DEFAULT false,
    create_time date
)
WITH (OIDS=FALSE);

CREATE TABLE tasks
(
    id serial PRIMARY KEY,
    type integer NOT NULL,
    total_points numeric DEFAULT 0,
    create_time date,
    start_time date,
    end_time date
)
WITH (OIDS=FALSE);

CREATE TABLE complete_history
(
    id serial PRIMARY KEY,
    user_id integer REFERENCES users (id),
    task_id integer REFERENCES tasks (id),
    earn_points numeric(100, 1) DEFAULT 0,
    create_time date
)
WITH (OIDS=FALSE);

INSERT INTO tasks (type, total_points, create_time, start_time, end_time) VALUES (0, 100, now(), now(), now() + interval '4 weeks');
INSERT INTO tasks (type, total_points, create_time, start_time, end_time) VALUES (1, 10000, now(), now(), now() + interval '1 week');
INSERT INTO tasks (type, total_points, create_time, start_time, end_time) VALUES (1, 10000, now(), now() + interval '1 week', now() + interval '2 weeks');
INSERT INTO tasks (type, total_points, create_time, start_time, end_time) VALUES (1, 10000, now(), now() + interval '2 weeks', now() + interval '3 weeks');
INSERT INTO tasks (type, total_points, create_time, start_time, end_time) VALUES (1, 10000, now(), now() + interval '3 weeks', now() + interval '4 weeks');
