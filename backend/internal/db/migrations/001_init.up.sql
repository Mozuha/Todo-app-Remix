CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,  -- Assuming single DB
  user_id UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),  -- Public-facing user identifier
  username VARCHAR(30) NOT NULL DEFAULT 'No Name', -- User-editable
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash BYTEA(60) NOT NULL,  -- Assuming bcrypt
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  CHECK (LENGTH(TRIM(username)) > 0)
);

CREATE TABLE todos (
  id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  description TEXT NOT NULL,
  position NUMERIC NOT NULL DEFAULT 0,
  completed BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  CHECK (LENGTH(TRIM(description)) > 0),
  CHECK (position >= 0)
);

-- Index on email for fast unique lookups during authentication and registration
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- Index on username for username-based searches
-- Could be changed to CREATE UNIQUE INDEX and use it for uniqueness checks if decided to force the username to be unique
CREATE INDEX idx_users_username ON users(username);

-- Index on user_id for efficient filtering of todos per user
-- This is especially important as most todo queries will be user-specific
CREATE INDEX idx_todos_user_id ON todos(user_id);

-- Composite index for efficient filtering and sorting of todos
-- Useful for queries that filter by user and want to sort/filter by completion or position
CREATE INDEX idx_todos_user_id_completed_position 
ON todos(user_id, completed, position);

-- Full-text search index on description
CREATE INDEX idx_todos_description_search 
ON todos USING GIN (to_tsvector('english', description));

-- Automatically updating updated_at column
-- https://zenn.dev/mpyw/articles/rdb-ids-and-timestamps-best-practices
CREATE FUNCTION refresh_updated_at_step1() RETURNS trigger AS
$$
BEGIN
  IF NEW.updated_at = OLD.updated_at THEN
    NEW.updated_at := NULL;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
    
CREATE FUNCTION refresh_updated_at_step2() RETURNS trigger AS
$$
BEGIN
  IF NEW.updated_at IS NULL THEN
    NEW.updated_at := OLD.updated_at;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION refresh_updated_at_step3() RETURNS trigger AS
$$
BEGIN
  IF NEW.updated_at IS NULL THEN
    NEW.updated_at := CURRENT_TIMESTAMP;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER refresh_users_updated_at_step1
  BEFORE UPDATE ON users FOR EACH ROW
  EXECUTE PROCEDURE refresh_updated_at_step1();
CREATE TRIGGER refresh_users_updated_at_step2
  BEFORE UPDATE OF updated_at ON users FOR EACH ROW
  EXECUTE PROCEDURE refresh_updated_at_step2();
CREATE TRIGGER refresh_users_updated_at_step3
  BEFORE UPDATE ON users FOR EACH ROW
  EXECUTE PROCEDURE refresh_updated_at_step3();

CREATE TRIGGER refresh_todos_updated_at_step1
  BEFORE UPDATE ON todos FOR EACH ROW
  EXECUTE PROCEDURE refresh_updated_at_step1();
CREATE TRIGGER refresh_todos_updated_at_step2
  BEFORE UPDATE OF updated_at ON todos FOR EACH ROW
  EXECUTE PROCEDURE refresh_updated_at_step2();
CREATE TRIGGER refresh_todos_updated_at_step3
  BEFORE UPDATE ON todos FOR EACH ROW
  EXECUTE PROCEDURE refresh_updated_at_step3();

-- Automatically rebalancing todo positions when needed
CREATE OR REPLACE FUNCTION rebalance_todo_positions()
RETURNS TRIGGER AS $$
DECLARE
    gap NUMERIC;
    min_gap NUMERIC := 1; -- Minimum allowed gap
BEGIN
    -- Find the smallest gap between consecutive todos
    SELECT MIN(t2.position - t1.position) INTO gap
    FROM todos t1
    JOIN todos t2 ON t1.user_id = t2.user_id AND t1.position < t2.position
    WHERE t1.user_id = NEW.user_id;

    -- If the smallest gap is too small, rebalance
    IF gap IS NOT NULL AND gap < min_gap THEN
        UPDATE todos
        SET position = (ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY position) * 100)
        WHERE user_id = NEW.user_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_rebalance_positions
AFTER UPDATE ON todos
FOR EACH ROW
WHEN (OLD.position IS DISTINCT FROM NEW.position)
EXECUTE FUNCTION rebalance_todo_positions();
