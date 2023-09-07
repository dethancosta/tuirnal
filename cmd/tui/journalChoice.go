package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dethancosta/tuirnal/internal/models"
)

// journalChoiceModel holds the model for the page wherein a journal
// is chosen
type journalChoiceModel struct {
	// ChoiceTi         textinput.Model // TODO delete
	NewJournalTi textinput.Model
	// ExistingJournals []string // TODO delete
	ExistingJournals list.Model
	Message          string
	NewJournalMode   bool
}

// journalItem is used in the ExistingJournals list
type journalItem struct {
	title string
}

func (j journalItem) FilterValue() string { return j.title }

var JcKeyMap = KeyMap{

	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),

	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/ctrl+c", "quit"),
	),
}

func initJournalChoice() journalChoiceModel {
	cti := textinput.New()
	cti.CharLimit = 50
	cti.Width = 32 //TODO change to be flexible
	cti.Prompt = "Name: "
	cti.Placeholder = "New journal name..."
	return journalChoiceModel{
		NewJournalTi:     cti,
		ExistingJournals: list.New(nil, list.NewDefaultDelegate(), 0, 0),
		Message:          "",
		NewJournalMode:   false,
	}
}

// journalChoiceView returns the string view to be returned by the
// parent bubbletea application
func journalChoiceView(m model) string {
	jc := &m.JournalChoice
	st := "Which journal would you like to use?\n" +
		jc.ChoiceTi.View() + "\n" +
		jc.Message + "\n\n" +
		jc.ExistingJournals.View() +
		helpStyle(jcHelpString())

	return st
}

// updateJournalChoice updates the parent bubbletea application when
// the JournalChoice page is active.
func updateJournalChoice(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	jc := &m.JournalChoice
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, JcKeyMap.Quit):
			jc.ChoiceTi.Blur()
			jc.ChoiceTi.Reset()
			jc.SelectIdx = 0
			jc.Message = ""
			m.ViewIdx = MenuIdx

		case key.Matches(msg, JcKeyMap.Select):
			i := jc.SelectIdx
			if i == 0 {
				if journalNameAvailable(m.App, m.CurrentAuthor, jc.ChoiceTi.Value()) {
					err := createJournal(m.App, m.CurrentAuthor, jc.ChoiceTi.Value())

					if err != nil {
						jc.Message = fmt.Sprintf("Couldn't create journal: %s", err.Error())
					} else {
						jc.ExistingJournals.InsertItem(journalItem{title: jc.ChoiceTi.Value()})
						jc.ExistingJournals = append(jc.ExistingJournals, jc.ChoiceTi.Value())
						jc.Message = "Journal created. You are now using it."
						m.CurrentJournal = jc.ChoiceTi.Value()
					}

				} else {
					m.CurrentJournal = jc.ExistingJournals.SelectedItem().FilterValue()
					jc.Message = "Now using journal " + m.CurrentJournal
					entries, err := m.App.EntryModel.GetAllEntries(m.CurrentAuthor, m.CurrentJournal)
					if err == nil {
						m.ViewEntry.SetCache(entries)
					}
				}
				jc.ChoiceTi.Reset()

			} else {
				m.CurrentJournal = jc.ExistingJournals[i-1]
				jc.ChoiceTi.Reset()
				jc.Message = "Now using journal " + m.CurrentJournal
				entries, err := m.App.EntryModel.GetAllEntries(m.CurrentAuthor, m.CurrentJournal)
				if err == nil {
					m.ViewEntry.SetCache(entries)
				}
			}

		default:
			var cmd tea.Cmd

			jc.ChoiceTi, cmd = jc.ChoiceTi.Update(msg)
			cmds = append(cmds, cmd)
			jc.ExistingJournals, cmd = jc.ExistingJournals.Update(msg)
			cmds = append(cmds, cmd)

			return m, tea.Batch(cmds...)
		}
	case tea.WindowSizeMsg:
		jc.ExistingJournals.SetSize(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	jc.ChoiceTi, cmd = jc.ChoiceTi.Update(msg)
	cmds = append(cmds, cmd)

	jc.ExistingJournals, cmd = jc.ExistingJournals.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (cm *journalChoiceModel) SetJournalsCache(journals []*models.Journal) {
	cache := make([]list.Item, len(journals))
	for i := range journals {
		cache[i] = journalItem{
			title: journals[i].Name,
		}
	}

	cm.ExistingJournals = list.New(cache, list.NewDefaultDelegate(), 0, 0)
	cm.ExistingJournals.Title = "Choose a journal"
}

/*
func jcHelpString() string {
	st := JcKeyMap.Up.Help().Key + ": "
	st += JcKeyMap.Up.Help().Desc + ",  "
	st += JcKeyMap.Down.Help().Key + ": "
	st += JcKeyMap.Down.Help().Desc + ",  "
	st += JcKeyMap.Select.Help().Key + ": "
	st += JcKeyMap.Select.Help().Desc + ",  "
	st += JcKeyMap.Quit.Help().Key + ": "
	st += JcKeyMap.Quit.Help().Desc + "\n"

	return st
}
*/
