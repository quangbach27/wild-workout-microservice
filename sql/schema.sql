CREATE TYPE availability_status AS ENUM (
    'available',
    'not_available',
    'training_scheduled'
);

CREATE TABLE dates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL UNIQUE
);

CREATE TABLE hours (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hour TIMESTAMP NOT NULL,
    availability availability_status NOT NULL,

    date_id UUID NOT NULL REFERENCES dates(id) ON DELETE CASCADE,
);
