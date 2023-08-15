package models

import "fmt"

type NoAuthorError struct {
	err string
}

func (e *NoAuthorError) Error() string {
	return fmt.Sprintf("%s", e.err)
}

type NoBucketError struct {
	err string
}

func (e *NoBucketError) Error() string {
	return e.err
}

type NoJournalError struct {
	err string
}

func (e *NoJournalError) Error() string {
	return e.err
}

type NoEntryError struct {
	err string
}

func (e *NoEntryError) Error() string {
	return e.err
}
