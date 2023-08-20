package models

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Journal struct {
	AuthorName string    `json:"author"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
}

type JournalModel struct {
	DB *bolt.DB
}

func (m *JournalModel) Insert(authorName string, journalName string) error {
	t := time.Now()
	journalJson, err := json.Marshal(Journal{
		AuthorName: authorName,
		Name:       journalName,
		CreatedAt:  t,
	})
	if err != nil {
		return err
	}

	err = m.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("journals"))
		if err != nil {
			return err
		}
		err = b.Put(append([]byte(strings.ToLower(authorName)+" "), []byte(strings.ToLower(journalName))...), journalJson)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *JournalModel) Get(authorName string, name string) (*Journal, error) {
	key := append([]byte(strings.ToLower(authorName)+" "), []byte(strings.ToLower(name))...)
	var row []byte

	err := m.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("journals"))
		if b == nil {
			return &NoBucketError{err: "no journals bucket"}
		}

		row = b.Get(key)
		if row == nil || len(row) == 0 {
			return &NoJournalError{err: "no journal with that name"}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var journal Journal
	err = json.Unmarshal(row, &journal)
	if err != nil {
		return nil, err
	}
	return &journal, nil
}

func (m *JournalModel) GetAllJournals(author string) ([]*Journal, error) {

	journals := make([]*Journal, 0)

	err := m.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("journals")).Cursor()
		if b == nil {
			return &NoJournalError{}
		}
		for k, v := b.Seek([]byte(author)); k != nil && bytes.HasPrefix(k, []byte(author)); k, v = b.Next() {
			var tempJournal Journal
			err := json.Unmarshal(v, &tempJournal)
			if err != nil {
				return err
			}
			journals = append(journals, &tempJournal)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(journals) == 0 {
		return nil, &NoAuthorError{err: "No journals found with author " + author}
	}

	return journals, nil
}
