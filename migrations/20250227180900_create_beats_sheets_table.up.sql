CREATE TABLE beats_sheets
(
  id uuid PRIMARY KEY NOT NULL,
  logline_id uuid NOT NULL,

  content jsonb NOT NULL,

  created_at timestamp(6) with time zone NOT NULL
);

CREATE INDEX beats_sheets_logline_id_idx ON beats_sheets (logline_id);
