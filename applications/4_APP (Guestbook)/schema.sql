CREATE TABLE users (
    username TEXT PRIMARY KEY,
    password TEXT NOT NULL
);
CREATE TABLE guestbook (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL REFERENCES users(username),
    written_at TIMESTAMP NOT NULL DEFAULT NOW(),
    "comment" TEXT NOT NULL
);
