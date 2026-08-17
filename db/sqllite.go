package db

import "database/sql"

func initSQLLite(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	initSqlStatement := `
		PRAGMA FOREIGN_KEYS ON
		
		CREATE TABLE IF NOT EXISTS flashcards (
			Id INTEGER PRIMARY KEY AUTOINCREMENT
			Question TEXT NOT NULL
			Answer TEXT NOT NULL
		)

		CREATE TABLE IF NOT EXISTS collections (
			Id INTEGER PRIMARY KEY AUTOINCREMENT
			Name TEXT NOT NULL
		)

		CREATE TABLE IF NOT EXISTS collection_flashcard (
			CollectionId INTEGER REFERENCES collections(Id)
			FlashcardId INTEGER REFERENCES flashcards(Id)
		)

		CREATE TABLE IF NOT EXISTS folders (
			Id INTEGER PRIMARY KEY AUTOINCREMENT
			Name TEXT NOT NULL
		)

		CREATE TABLE IF NOT EXISTS folder_collection (
			FolderId INTEGER REFERENCES folders(Id)
			CollectionId INTEGER REFERENCES collections(Id)
		)
	`

	_, err = db.Exec(initSqlStatement)
	if err != nil {
		return nil, err
	}

	return db, nil
}
