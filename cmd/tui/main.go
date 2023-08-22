package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dethancosta/tuirnal/internal/helpers"
)

const (
	LoginIdx         = iota
	MenuIdx
	JournalChoiceIdx
	NewEntryIdx
	ViewEntryIdx
	ViewJournalIdx
)

type model struct {
	App            helpers.Application
	CurrentAuthor  string //username of the current user
	CurrentJournal string //name of the current journal
	ViewEntry      viewEntryModel
	Login          loginModel
	Menu           MenuPage
	WriteEntryPage WriteEntry
	JournalChoice  journalChoiceModel
	ViewIdx        int //current view to be displayed
}

func initialModel() model {

	app := *helpers.InitApp("tuirnal.db")

	return model{
		App:            app,
		CurrentAuthor:  "",
		CurrentJournal: "",
		ViewEntry:      initViewEntry(),
		Login:          initLogin(),
		Menu:           initMenuPage(),
		WriteEntryPage: initWriteEntry(),
		JournalChoice:  initJournalChoice(),
		ViewIdx:        0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch m.ViewIdx {
	case LoginIdx:
		return updateLogin(msg, m)
	case MenuIdx:
		return updateMenu(msg, m)
	case JournalChoiceIdx:
		return updateJournalChoice(msg, m)
	case NewEntryIdx:
		return updateWriteEntry(msg, m)
	case ViewEntryIdx:
		return updateViewEntry(msg, m)
	case ViewJournalIdx:
		return updateEntries(msg, m)
	}

	return m, nil
}

func (m model) View() string {
	st := ""
	if m.ViewIdx != LoginIdx {
		st += m.CurrentAuthor + " || "
		if len(strings.TrimSpace(m.CurrentJournal)) > 0 {
			st += m.CurrentJournal + "\n\n\n"
		} else {
			st += "*No journal selected*" + "\n\n\n"
		}
	} else {
		st += "\n\n\n"
	}

	switch m.ViewIdx {
	case LoginIdx:
		st += loginView(m)
	case MenuIdx:
		st += menuView(m)
	case JournalChoiceIdx:
		st += journalChoiceView(m)
	case NewEntryIdx:
		st += writeEntryView(m)
	case ViewEntryIdx:
		st += viewEntryView(m)
	case ViewJournalIdx:
		st += entriesView(m)
	}
	return st
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
