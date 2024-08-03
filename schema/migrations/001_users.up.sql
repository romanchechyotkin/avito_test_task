CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    user_type user_type NOT NULL
);