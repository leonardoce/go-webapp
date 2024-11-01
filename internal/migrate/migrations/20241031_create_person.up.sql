BEGIN;

CREATE TABLE person (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name TEXT NOT NULL,
  surname TEXT NOT NULL,
  phone_number TEXT,
  address_first_line TEXT,
  address_second_line TEXT,
  address_postal_code TEXT
);

COMMIT;