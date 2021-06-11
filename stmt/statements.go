package stmt

const CreateTables = `

CREATE TABLE IF NOT EXISTS island
(
	id         text    PRIMARY KEY,
  name       text    NOT NULL,
  avatar     text    NOT NULL,
  email      text    NOT NULL UNIQUE,
  link       text    NOT NULL,
  address    text    NOT NULL,
);

CREATE TABLE IF NOT EXISTS message
(
  id       text    PRIMARY KEY,
  ctime    int     NOT NULL,
  at       text    NOT NULL,
  body     text    NOT NULL,
  md       int     NOT NULL,
);

CREATE TABLE IF NOT EXISTS island_msg
(
  island_id    text    REFERENCES island(id) ON DELETE CASCADE,
  msg_id       text    REFERENCES message(id) ON DELETE CASCADE,
  UNIQUE (island_id, msg_id)
);

CREATE TABLE IF NOT EXISTS metadata
(
  name         text    NOT NULL UNIQUE,
  int_value    int     NOT NULL DEFAULT 0,
  text_value   text    NOT NULL DEFAULT "" 
);
`
