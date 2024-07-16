USE snippetbox;
CREATE TABLE sessions(
  token CHAR(43) PRIMARY KEY,
  data BLOB NOT NULL,
  expiry DATETIME NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

