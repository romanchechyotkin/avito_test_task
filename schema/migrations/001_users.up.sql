CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    user_type user_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);