package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dethancosta/tuirnal/internal/helpers"
)

const (
	LoginIdx         = 0
	MenuIdx          = 1
	JournalChoiceIdx = 2
	NewEntryIdx      = 3
	ViewEntryIdx     = 4
	ViewJournalIdx   = 5
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
		st += m.CurrentAuthor + " || " + m.CurrentJournal + "\n\n\n"
	} else {
		st += "\n\n\n"
	}

	switch m.ViewIdx {
	case 0:
		st += loginView(m)
	case 1:
		st += menuView(m)
	case 2:
		st += journalChoiceView(m)
	case 3:
		st += writeEntryView(m)
	case 4:
		st += viewEntryView(m)
	case 5:
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
