CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE NOT NULL,
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128),
    CONSTRAINT unique_date_title UNIQUE (date, title)
)