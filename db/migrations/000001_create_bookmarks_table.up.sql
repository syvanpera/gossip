CREATE TABLE IF NOT EXISTS bookmarks(
    id INTEGER PRIMARY KEY,
    url TEXT UNIQUE,
    description TEXT,
    tags TEXT,
    flags INTEGER,

    created_at datetime,
    updated_at datetime
);
