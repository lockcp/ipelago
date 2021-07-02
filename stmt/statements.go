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
  status     text    NOT NULL,
  checked    int     NOT NULL
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
  id           text    PRIMARY KEY,
  island_id    text    REFERENCES island(id) ON DELETE CASCADE,
  time         int     NOT NULL,
  body         text    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_message_time ON message(time);
CREATE INDEX IF NOT EXISTS idx_message_id_time ON message(id, time);

CREATE TABLE IF NOT EXISTS denylist
(
  ctime      int     PRIMARY KEY,
  address    text    NOT NULL UNIQUE
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
  SELECT id, name, email, avatar, link, address, note, status, checked
  FROM island WHERE id=?;`

const AllIslands = `
  SELECT id, name, email, avatar, link, address, note, status, checked
  FROM island WHERE id<>? ORDER BY id DESC;`

const GetMoreMessagesByIsland = `
  SELECT id, island_id, time, body FROM message
  WHERE island_id=? AND time<? ORDER BY time DESC LIMIT ?;`

const GetMoreMessages = `
  SELECT msg.id, island_id, msg.time, msg.body FROM message AS msg
  INNER JOIN island ON msg.island_id = island.id
  WHERE msg.time<? and island.status<>"unfollowed" ORDER BY msg.time DESC LIMIT ?;`

const DeleteIsland = `
  DELETE FROM island WHERE id=?;`

const InsertIsland = `
  INSERT INTO island (id, name, email, avatar, link, address, note, status, checked)
  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

const UpdateIsland = `
  UPDATE island
  SET name=?, email=?, avatar=?, link=?, status=?
  WHERE id=?;`

const UpdateIslandChecked = `
  UPDATE island SET checked=? WHERE id=?;`

const SetStatus = `
  UPDATE island SET status=? WHERE id=?;`

const UpdateNote = `
  UPDATE island SET note=? WHERE id=?;`

const InsertMsg = `
  INSERT INTO message (id, island_id, time, body)
  VALUES (?, ?, ?, ?);`

const DeleteMessage = `
  DELETE FROM message WHERE id=?;`

const CountMessages = `
  SELECT count(*) FROM message WHERE island_id=?;`

const InsertDeny = `
  INSERT INTO denylist (ctime, address) VALUES (?, ?);`

const DeleteDeny = `
  DELETE FROM denylist WHERE address=?;`

const GetDenyList = `
  SELECT address FROM denylist ORDER BY ctime DESC;`

const CountDeny = `
  SELECT count(*) FROM denylist WHERE address=?;`

const CountIsland = `
  SELECT count(*) FROM island WHERE address=?;`
