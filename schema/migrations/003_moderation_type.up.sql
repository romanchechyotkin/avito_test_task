DO $$
BEGIN

IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'moderation_type') THEN
    CREATE TYPE moderation_type AS ENUM ('created', 'approved', 'declined', 'on moderation');
END IF;

END $$;

