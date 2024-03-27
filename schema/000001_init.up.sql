CREATE TABLE users
(
    uid         serial PRIMARY KEY,
    email       varchar(255) not null unique,
    password    varchar(255) not null
);

CREATE TABLE sessions
(
    uid    integer       not null,
    token  varchar(255) not null unique,
    exp_at timestamp    not null,
    FOREIGN KEY (uid) REFERENCES users (uid) ON DELETE CASCADE
);

CREATE TABLE tasks
(
    sid             serial PRIMARY KEY,
    maintainer      integer      not null,
    password        varchar(255) not null,
    title           varchar(255) not null,
    description     text,
    last_change     timestamp    not null,
    task_type       varchar(255) not null,
    CHECK (task_type = 'group' OR task_type = 'individuals'),
    method          varchar(255) not null,
    CHECK (method = 'topsis' OR method = 'smart'),
    calc_settings   integer      not null,
    ling_scale      json         not null,
    alternatives    json,
    criteria        json,
    experts_weights json,
    status          boolean      not null,
    FOREIGN KEY (maintainer) REFERENCES users (uid) ON DELETE CASCADE
);

ALTER TABLE tasks
    ADD CONSTRAINT uqc_user_title
        UNIQUE (maintainer, title);

CREATE TABLE matrices (
    mid serial PRIMARY KEY,
    sid integer not null,
    uid integer not null,
    matrix json,
    status boolean not null,
    FOREIGN KEY (sid) REFERENCES tasks (sid) ON DELETE CASCADE,
    FOREIGN KEY (uid) REFERENCES users (uid) ON DELETE  CASCADE
);

ALTER TABLE matrices
    ADD CONSTRAINT uqc_sid_uid
        UNIQUE (uid, sid);

CREATE TABLE final
(
    fid           integer PRIMARY KEY,
    result        json,
    sens_analysis json,
    threshold     double precision not null,
    last_change   timestamp        not null,
    FOREIGN KEY (fid) REFERENCES tasks (sid) ON DELETE CASCADE
);