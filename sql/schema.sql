-- Drop in correct order (hours depends on dates)
DROP TABLE IF EXISTS hours;
DROP TABLE IF EXISTS dates;
DROP TYPE IF EXISTS availability_status;

-- Recreate ENUM type
CREATE TYPE availability_status AS ENUM (
    'available',
    'not_available',
    'training_scheduled'
);

-- Create dates table
CREATE TABLE dates (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL UNIQUE
);

-- Create hours table
CREATE TABLE hours (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hour TIMESTAMP NOT NULL,
    availability availability_status NOT NULL,
    date_id UUID NOT NULL REFERENCES dates(uuid) ON DELETE CASCADE,
    UNIQUE (date_id, hour)
);
