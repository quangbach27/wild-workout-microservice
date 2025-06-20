-- Recreate ENUM type
CREATE TYPE availability_status AS ENUM (
    'available',
    'not_available',
    'training_scheduled'
);

-- Create dates table
CREATE TABLE dates (
   date DATE NOT NULL UNIQUE,
   PRIMARY KEY(date)
);

-- Create hours table
CREATE TABLE hours (
   hour TIMESTAMP NOT NULL,
   availability availability_status NOT NULL,
   date DATE NOT NULL REFERENCES dates(date) ON DELETE CASCADE,
   PRIMARY KEY (hour)
);
