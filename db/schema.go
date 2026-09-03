package db

import "database/sql"

const schema = `
	CREATE TABLE IF NOT EXISTS FlashcardSets (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT NOT NULL,
		Description TEXT,
		LastAccessed DATETIME DEFAULT CURRENT_TIMESTAMP,
		TrackProgress INTEGER NOT NULL DEFAULT 0 CHECK (TrackProgress IN (0, 1)),
		Front INTEGER NOT NULL DEFAULT 0 CHECK (Front IN (0, 1)),

		Shuffle INTEGER NOT NULL DEFAULT 0 CHECK (Shuffle IN (0, 1)),
		ShuffleSeed INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS Flashcards (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Question TEXT NOT NULL,
		Answer TEXT NOT NULL,
		FlashcardSetId INTEGER NOT NULL,
		Position INTEGER NOT NULL,

		UNIQUE (FlashcardSetId, Position),

		FOREIGN KEY(FlashcardSetId) REFERENCES FlashcardSets(Id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS Quizzes (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		FlashcardSetId INTEGER UNIQUE NOT NULL,
		CurrentlySelectedIndex INTEGER NOT NULL DEFAULT 0,

		FOREIGN KEY(FlashcardSetId) REFERENCES FlashcardSets(Id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS Folders (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT NOT NULL,
		LastAccessed DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS FolderFlashcardSet (
		FolderId INTEGER NOT NULL,
		FlashcardSetId INTEGER NOT NULL,

		PRIMARY KEY (FolderId, FlashcardSetId),

		FOREIGN KEY(FolderId) REFERENCES Folders(Id) ON DELETE CASCADE,
		FOREIGN KEY(FlashcardSetId) REFERENCES FlashcardSets(Id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS QuizFlashcards (
		QuizId INTEGER NOT NULL,
		FlashcardId INTEGER NOT NULL,
		Position INTEGER NOT NULL,
		IsUnknown INTEGER NOT NULL DEFAULT 0 CHECK (IsUnknown IN (0, 1)),

		PRIMARY KEY (QuizId, Position),

		FOREIGN KEY(FlashcardId) REFERENCES Flashcards(Id) ON DELETE CASCADE,
		FOREIGN KEY(QuizId) REFERENCES Quizzes(Id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_flashcards ON Flashcards(Id);
	CREATE INDEX IF NOT EXISTS idx_flashcardSets ON FlashcardSets(Id);
	CREATE INDEX IF NOT EXISTS idx_folders ON Folders(Id);
	CREATE INDEX IF NOT EXISTS idx_quizzes ON Quizzes(Id);
`

func InitializeSchema() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./app.db?_foreign_keys=ON")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Schema() string {
	return schema
}
