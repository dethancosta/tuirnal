package main

import (
	"errors"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dethancosta/tuirnal/internal/auth"
	"github.com/dethancosta/tuirnal/internal/helpers"
	"golang.org/x/crypto/bcrypt"
)

// entryNameValidator returns true if the title is available in the given journal,
// and false if it is taken
func entryNameAvailable(app helpers.Application, author, journal, title string) bool {
	_, err := app.EntryModel.Get(author, journal, title)
	//TODO check that err is type NoEntryErr
	if err != nil {
		return true
	}

	return false
}

// saveJournalEntry saves a journal entry in the application's storage layer,
// returning an error if unsuccessful
func saveJournalEntry(app helpers.Application,
	author, journal, title, tags, content string) error {

	err := app.EntryModel.Insert(author, journal, title, content, strings.Fields(tags))

	return err
}

// journalNameAvailable checks whether a journal with a given name exists
func journalNameAvailable(app helpers.Application, name, author string) bool {
	_, err := app.JournalModel.Get(author, name)
	//TODO verify err is NoJournalErr type
	if err != nil {
		return true
	}

	return false
}

func createJournal(app helpers.Application, author, name string) error {

	if journalNameAvailable(app, name, author) {
		err := app.JournalModel.Insert(strings.ToLower(author), name)
		return err
	} else {
		return errors.New("Journal name taken")
	}
}

// loginAuthor attempts to sign a user in, and returns true if successful
func loginAuthor(app helpers.Application, author, password string) (bool, error) {

	// TODO return author object
	_, err := auth.Authenticate(app, author, password)
	return err == nil, err
}

// authorNameAvailable checks whether an account with a given username exists
func authorNameAvailable(app helpers.Application, author string) (bool, error) {
	_, err := app.AuthorModel.Get(author)
	//TODO verify that err is NoAuthorErr
	if err != nil {
		return true, err
	}

	return false, nil
}

// createAuthor creates a new author, returning true if successful
func createAuthor(app helpers.Application, author string, password string) error {
	ok, err := authorNameAvailable(app, author)
	if !ok {
		return errors.New("Author with that username already exists")
	}
	//TODO include error checking in case err isn't NoAuthorErr
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = app.AuthorModel.Insert(author, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

// Members of the SelectorModel interface contain a selection index
type SelectorModel interface {
	GetSelection() int
}

// selected formats a string based on whether it is currently selected or not
func selected(m SelectorModel, idx int, s string) string {
	var res string
	if m.GetSelection() == idx {
		res += "➤ "
	} else {
		res += "  "
	}
	res += s
	return res
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
