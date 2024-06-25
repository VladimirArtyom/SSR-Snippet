-- Create a new UTF-8 snippetbox database

DROP DATABASE IF EXISTS snippetbox;
DROP USER IF EXISTS web@localhost;

CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE snippetbox;

CREATE TABLE snippets(
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  created DATETIME NOT NULL,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  expires DATETIME NOT NULL
);


CREATE INDEX idx_snippets_created ON snippets(created);


-- Create new User
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
ALTER USER 'web'@'localhost' IDENTIFIED BY '123456';

-- DML snippetbox db table snippets
INSERT INTO snippets(title, content, created, expires) VALUES (
  "An old silent pond",
  "An old silent pond...\nA frog jumps into the pond, \nsplash! Silence again.\n\n - Matsuo",
  UTC_TIMESTAMP(),
  DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
  );

INSERT INTO snippets(title, content, created, expires) VALUES (
"Over the wintry forest",
"Over the wintry forst\n, winds howl in rage\n with no leaves to blow. \n\n- Natsume Soseki",
UTC_TIMESTAMP(),
DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets(title, content, created, expires) VALUES (
"First autumn morning",
"First autumn morning\n the mirror I stare into\nshows my fathers face.\n\n- Murakami Katsushika",
UTC_TIMESTAMP(),
DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);
