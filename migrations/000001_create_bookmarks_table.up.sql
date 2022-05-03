CREATE TABLE IF NOT EXISTS bookmarks(
    id INTEGER NOT NULL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    tags TEXT NOT NULL DEFAULT '',
    flags INTEGER NOT NULL DEFAULT 0,

    created_at DATETIME DEFAULT (datetime('now')),
    updated_at DATETIME DEFAULT (datetime('now'))
);
