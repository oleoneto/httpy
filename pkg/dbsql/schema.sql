DROP TABLE IF EXISTS responses;

CREATE TABLE responses (
  id INTEGER PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  status_code INTEGER NOT NULL,
  method CHAR(10) NOT NULL,
  url VARCHAR(255) NOT NULL,
  headers TEXT,
  body TEXT,
  timeout_ms INTEGER,

  CONSTRAINT http_status_code CHECK (status_code >= 100 AND status_code <= 599),
  CONSTRAINT http_method      CHECK (method IN ("GET", "PATCH", "POST", "PUT", "DELETE", "TRACE", "HEAD", "OPTIONS", "CONNECT")),
  CONSTRAINT http_timeout     CHECK (timeout_ms > 0)
);
