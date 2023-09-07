package models

import (
	bolt "go.etcd.io/bbolt"
)

type Author struct {
	Username string
	Password string
}

type AuthorModel struct {
	DB *bolt.DB
}

func (m *AuthorModel) Insert(username string, passwordHash []byte) error {
	err := m.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("authors"))
		if err != nil {
			return err
		}
		err = b.Put([]byte(username), passwordHash)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (m *AuthorModel) Get(name string) (*Author, error) {
	tx, err := m.DB.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte("authors"))
	if b == nil {
		return nil, &NoBucketError{err: "Authors bucket doesn't exist"}
	}

	row := b.Get([]byte(name))
	if row == nil {
		return nil, &NoAuthorError{err: "There is no author with that email"}
	}

	password := string(row)
	return &Author{
		Username: name,
		Password: password,
	}, nil
}
