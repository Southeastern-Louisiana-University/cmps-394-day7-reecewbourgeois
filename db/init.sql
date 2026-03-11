-- Create the messages table
CREATE TABLE IF NOT EXISTS messages (
    id      SERIAL PRIMARY KEY,
    content TEXT NOT NULL
);

-- Seed an initial message
INSERT INTO messages (content) VALUES ('Hello, World!');
