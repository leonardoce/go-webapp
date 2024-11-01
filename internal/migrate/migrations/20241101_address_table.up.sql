BEGIN;
  CREATE TABLE address (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    person_id BIGINT REFERENCES person(id),
    first_line TEXT NOT NULL,
    second_line TEXT,
    postal_code TEXT
  );

  INSERT INTO address (person_id, first_line, second_line, postal_code)
    SELECT 
      id, 
      address_first_line, 
      address_second_line, 
      address_postal_code
    FROM person
    WHERE address_first_line IS NOT NULL
    FOR UPDATE;
    
  ALTER TABLE person 
    DROP COLUMN address_first_line,
    DROP COLUMN address_second_line,
    DROP COLUMN address_postal_code;
COMMIT;
