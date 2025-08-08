ALTER TABLE loglines
ADD COLUMN lang text NOT NULL DEFAULT 'en';

ALTER TABLE beats_sheets
ADD COLUMN lang text NOT NULL DEFAULT 'en';
