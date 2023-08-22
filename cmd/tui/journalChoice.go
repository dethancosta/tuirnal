package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dethancosta/tuirnal/internal/models"
)

// journalChoiceModel holds the model for the page wherein a journal
// is chosen
type journalChoiceModel struct {
	ChoiceTi         textinput.Model
	ExistingJournals []string
	Message          string
	SelectIdx        int
}

// journalChoiceModel satisfies the SelectorModel interface
func (cm journalChoiceModel) GetSelection() int {
	return cm.SelectIdx
}

func initJournalChoice() journalChoiceModel {
	cti := textinput.New()
	cti.CharLimit = 50
	cti.Width = 32 //TODO change to be flexible
	cti.Prompt = "Name: "
	cti.Placeholder = "Journal name..."
	return journalChoiceModel{
		ChoiceTi:         cti,
		ExistingJournals: nil,
		Message:          "",
		SelectIdx:        0,
	}
}

// journalChoiceView returns the string view to be returned by the
// parent bubbletea application
func journalChoiceView(m model) string {
	jc := &m.JournalChoice
	st := "Which journal would you like to use?\n" +
		jc.ChoiceTi.View() + "\n" +
		jc.Message + "\n\n" +
		jc.SelectionString()

	return st
}

// updateJournalChoice updates the parent bubbletea application when
// the JournalChoice page is active.
func updateJournalChoice(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	jc := &m.JournalChoice
	var cmd tea.Cmd

	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+c":
			jc.ChoiceTi.Blur()
			jc.ChoiceTi.Reset()
			jc.SelectIdx = 0
			jc.Message = ""
			m.ViewIdx = MenuIdx
		case "enter":
			i := jc.SelectIdx
			if i == 0 {
				if journalNameAvailable(m.App, m.CurrentAuthor, jc.ChoiceTi.Value()) {
					err := createJournal(m.App, m.CurrentAuthor, jc.ChoiceTi.Value())
					if err != nil {
						jc.Message = fmt.Sprintf("Couldn't create journal: %s", err.Error())
					} else {
						jc.ExistingJournals = append(jc.ExistingJournals, jc.ChoiceTi.Value())
						jc.Message = "Journal created. You are now using it."
						m.CurrentJournal = jc.ChoiceTi.Value()
					}
				} else {
					m.CurrentJournal = jc.ChoiceTi.Value()
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
			}
		case "down":
			jc.SelectIdx++
			if jc.SelectIdx > len(jc.ExistingJournals) {
				jc.SelectIdx = len(jc.ExistingJournals)
			}
			if jc.SelectIdx == 0 {
				jc.ChoiceTi.Focus()
			} else {
				jc.ChoiceTi.Blur()
			}
		case "up":
			jc.SelectIdx--
			if jc.SelectIdx < 0 {
				jc.SelectIdx = 0
			}
			if jc.SelectIdx == 0 {
				jc.ChoiceTi.Focus()
			} else {
				jc.ChoiceTi.Blur()
			}
		default:
			jc.ChoiceTi, cmd = jc.ChoiceTi.Update(msg)
			return m, cmd
		}
	}
	jc.ChoiceTi, cmd = jc.ChoiceTi.Update(msg)
	return m, cmd
}

// helper

// SelectionString is a helper function that formats the
// list of journals in a JournalChoice model
func (cm journalChoiceModel) SelectionString() string {
	sb := strings.Builder{}
	for i, s := range cm.ExistingJournals {
		sb.WriteString(selected(cm, i+1, s) + "\n")
	}

	return sb.String() + "\n"
}

func (cm *journalChoiceModel) SetJournalsCache(journals []*models.Journal) {
	cache := make([]string, len(journals))
	for i := range journals {
		cache[i] = journals[i].Name
	}
	cm.ExistingJournals = cache
}
