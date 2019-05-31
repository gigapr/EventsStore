CREATE TABLE IF NOT EXISTS Events
(
  Id SERIAL PRIMARY KEY,
  SourceId VARCHAR(255)  NOT NULL,
  EventType VARCHAR(255)  NOT NULL,
  EventData BYTEA NOT NULL,
  Received timestamptz NOT NULL DEFAULT now()
)
