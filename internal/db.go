package internal

import (
	"database/sql"
	"strconv"
	"time"
)

//NOTE: implement tags according to: https://stackoverflow.com/questions/20856/recommended-sql-database-design-for-tags-or-tagging
// > table items with itemid column
// > table tags with tagid column
// > table itemtags with mapping (two columns, one is itemid and other is tagid)

type db struct {
	DB *sql.DB
}

func (db *db) AddJournal(j Journal) error {
	_, err := db.DB.Exec(`INSERT INTO Journals(author, name, createdAt)
				VALUES (?, ?, ?)`, j.authorName, j.name, j.createdAt)
	return err
}

func (db *db) GetJournal(name string, author string) (*Journal, error) {
	row := db.DB.QueryRow(`SELECT * FROM Journals WHERE author=? AND name=name`, author, name)
	journal := NewJournal("")
	err := row.Scan(journal.authorName, journal.name, journal.createdAt)
	if err != nil {
		return journal, nil
	}
	return nil, err
}

func (db *db) GetJournalsByAuthor(author string) ([]Journal, error) {
	rows, err := db.DB.Query(`SELECT * FROM Journals WHERE author=?`, author)
	if err != nil {
		return nil, err
	}
	journals := make([]Journal, 0, 1)
	for rows.Next() {
		journal := Journal{
			name:       "",
			authorName: "",
			createdAt:  time.Now(),
		}
		err = rows.Scan(journal.authorName, journal.name, journal.createdAt)
		if err != nil {
			return nil, err
		}
		journals = append(journals, journal)
	}
	return journals, nil
}

func (db *db) AddJournalEntry(je JournalEntry) error {
	// TODO roll this into a tx so that a rollback is executed if any one of the insertions fails
	var err error
	entryId, err := db.DB.Exec(`INSERT INTO Entries(writtenAt, journal, content) VALUE`)
	if err != nil {
		return err
	}
	// TODO: implement check to see if tag is already in the table
	tagIds := make([]string, len(je.tags))
	for i, t := range je.tags {
		tagID, err := db.DB.Exec(`INSERT INTO Tags(tagName) VALUES (?)`, t)
		if err != nil {
			return err
		}
		tagIdStr, err := tagID.LastInsertId()
		if err != nil {
			return err
		}
		tagIds[i] = strconv.Itoa(int(tagIdStr))
	}
	i, err := entryId.LastInsertId()
	if err != nil {
		return err
	}

	for _, t := range tagIds {
		// TODO make sure that each row in itemTags is its own key
		_, err = db.DB.Exec(`INSERT INTO ItemTags(itemId, tagId) VALUES(?, ?)`, i, t)
		if err != nil {
			return err
		}
	}
	return err
}

func (db *db) GetJournalEntries(journalName string) ([]JournalEntry, error) {
	// TODO implement
	return nil, nil
}

func (db *db) GetEntriesBeforeDate(journalName string, beforeDate time.Time) {}
