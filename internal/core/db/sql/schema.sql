CREATE TABLE IF NOT EXISTS archive (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL,
    editorial TEXT NOT NULL,
    cover_page INTEGER NOT NULL DEFAULT 1,
    favorite BOOLEAN NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    CHECK (cover_page >= 1)
);
