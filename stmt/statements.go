package stmt

const CreateTables = `

CREATE TABLE IF NOT EXISTS island
(
	id         text    PRIMARY KEY,
  name       text    NOT NULL,
  email      text    NOT NULL,
  avatar     text    NOT NULL,
  link       text    NOT NULL,
  address    text    NOT NULL UNIQUE,
  note       text    NOT NULL,
  status     text    NOT NULL
);

CREATE TABLE IF NOT EXISTS cluster
(
  id      text    PRIMARY KEY,
  name    text    NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS island_cluster
(
  island_id     text    REFERENCES island(id) ON DELETE CASCADE,
  cluster_id    text    REFERENCES cluster(id) ON DELETE CASCADE,
  UNIQUE (island_id, cluster_id)
);

CREATE TABLE IF NOT EXISTS message
(
  id       text    PRIMARY KEY,
  ctime    int     NOT NULL,
  at       text    NOT NULL,
  body     text    NOT NULL,
  md       int     NOT NULL
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
const InsertIntValue = `INSERT INTO metadata (name, int_value) VALUES (?, ?);`
const GetIntValue = `SELECT int_value FROM metadata WHERE name=?;`
const UpdateIntValue = `UPDATE metadata SET int_value=? WHERE name=?;`

const InsertTextValue = `INSERT INTO metadata (name, text_value) VALUES (?, ?);`
const GetTextValue = `SELECT text_value FROM metadata WHERE name=?;`
const UpdateTextValue = `UPDATE metadata SET text_value=? WHERE name=?;`

const GetIslandByID = `
    SELECT id, name, avatar, email, link, address, note
    FROM island WHERE id=?;`

const GetLastMessage = `
    SELECT message.id, message.ctime, message.at, message.body, message.md
    FROM island INNER JOIN island_msg ON island.id = island_msg.island_id
    INNER JOIN message ON island_msg.msg_id = message.id
    WHERE island.id=? ORDER BY message.ctime DESC LIMIT 1;`

const InsertIsland = `
    INSERT INTO island (id, name, email, avatar, link, address, note, status)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

const UpdateIsland = `
    UPDATE island
    SET name=?, email=?, avatar=?, link=?, address=?, note=?, status=?
    WHERE id=?;`
