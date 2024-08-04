CREATE TABLE IF NOT EXISTS house_subscriptions (
    house_id SERIAL REFERENCES houses(id),
    user_id UUID REFERENCES users(id),
    PRIMARY KEY (house_id, user_id)
)