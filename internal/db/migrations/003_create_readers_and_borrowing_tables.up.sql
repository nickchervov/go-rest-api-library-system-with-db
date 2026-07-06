CREATE TABLE readers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL DEFAULT '',
			email TEXT UNIQUE NOT NULL DEFAULT '',
			phone TEXT
		);
CREATE TABLE borrowing (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			book_id INTEGER NOT NULL DEFAULT 0,
			readers_id INTEGER NOT NULL DEFAULT 0,
			borrow_date TEXT NOT NULL DEFAULT '',
			return_date TEXT,
			FOREIGN KEY (book_id) REFERENCES books(id),
			FOREIGN KEY (readers_id) REFERENCES readers(id)
		);