--setup.sql
CREATE TABLE IF NOT EXISTS userstype(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    rank TEXT NOT NULL
);

INSERT INTO userstype (id, rank) VALUES(0,"guest"),(1, "user"),(2, "moderator"),(4, "administrator");

CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userstypeid INTEGER NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT NOT NULL,
    validation INTEGER NOT NULL,
    time DATETIME NOT NULL,
    session_token TEXT,
    FOREIGN KEY(userstypeid) REFERENCES userstype(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    authorid INTEGER NOT NULL,
    author TEXT NOT NULL,
    category TEXT NOT NULL,
    title TEXT NOT NULL UNIQUE,
    content TEXT NOT NULL,
    like INTEGER NOT NULL,
    dislikes INTEGER NOT NULL,
    FOREIGN KEY(authorid) REFERENCES users(id) ON DELETE CASCADE
);
