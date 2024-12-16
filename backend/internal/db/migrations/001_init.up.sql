CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id SERIAL PRIMARY KEY,  -- Assuming single DB
  user_id UUID UNIQUE NOT NULL DEFAULT uuid_generate_v4(),  -- Public-facing user identifier
  username VARCHAR(30) NOT NULL DEFAULT 'No Name', -- User-editable
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE todos (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  description TEXT NOT NULL,
  position INTEGER NOT NULL DEFAULT 0,
  completed BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
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