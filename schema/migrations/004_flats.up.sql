CREATE TABLE IF NOT EXISTS flats (
    id SERIAL PRIMARY KEY,
    number INT NOT NULL,
    house_id SERIAL REFERENCES houses(id),
    price INT NOT NULL,
    rooms_amount INT NOT NULL,
    moderation_status moderation_type DEFAULT 'created',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (house_id, number)
);