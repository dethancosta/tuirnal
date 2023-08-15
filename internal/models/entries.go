package models

import (
	"bytes"
	"encoding/json"
	"time"

	bolt "go.etcd.io/bbolt"
)

type JournalEntry struct {
	WrittenAt     time.Time `json:"writtenAt"`
	Author        string    `json:"authorName"`
	ParentJournal string    `json:"parentJournal"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Tags          []string  `json:"tags"`
}

type EntryModel struct {
	DB *bolt.DB
}

func (m *EntryModel) Insert(author, journal, title, content string, tags []string) error {
	t := time.Now()
	entryJson, err := json.Marshal(JournalEntry{
		WrittenAt:     t,
		Author:        author,
		ParentJournal: journal,
		Title:         title,
		Content:       content,
		Tags:          tags,
	})

	if err != nil {
		return err
	}

	err = m.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("entries"))
		if err != nil {
			return err
		}
		k := append([]byte(author), []byte(journal)...)
		k = append(k, []byte(title)...)
		err = b.Put(k, entryJson)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (m *EntryModel) Get(author, journal, name string) (*JournalEntry, error) {
	var entry JournalEntry
	err := m.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("entries"))
		key := append([]byte(author), []byte(journal)...)
		key = append(key, []byte(name)...)
		row := b.Get(key)
		if row == nil {
			return &NoEntryError{}
		}
		err = json.Unmarshal(row, &entry)
		return err
	})

	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (m *EntryModel) GetAllEntries(author, journal string) ([]*JournalEntry, error) {

	entries := make([]*JournalEntry, 0, 4)
	err := m.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("entries")).Cursor()
		if b == nil {
			return &NoBucketError{}
		}
		keyPfx := append([]byte(author), []byte(journal)...)
		var tempEntry JournalEntry
		for k, v := b.Seek(keyPfx); k != nil && bytes.HasPrefix(k, keyPfx); k, v = b.Next() {
			err := json.Unmarshal(v, &tempEntry)
			if err != nil {
				return err
			}
			entries = append(entries, &tempEntry)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, &NoJournalError{err: "Journal doesn't exist with name " + journal}
	}

	return entries, nil
}
