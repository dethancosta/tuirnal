package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type MenuPage struct {
	SelectIdx int
	Message   string
}

var MenuKeyMap = KeyMap{

	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "move up"),
	),

	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "move down"),
	),

	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),

	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "quit"),
	),
}

func (mp MenuPage) GetSelection() int {
	return mp.SelectIdx
}

func initMenuPage() MenuPage {
	return MenuPage{
		SelectIdx: 0,
	}
}

func updateMenu(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	mp := &m.Menu
	if _, ok := msg.(tea.KeyMsg); !ok {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, MenuKeyMap.Quit):
			lp := &m.Login
			lp.NameTi.Reset()
			lp.PasswordTi.Reset()
			lp.NameTi.Focus()
			lp.SelectIdx = 0
			mp.SelectIdx = 0
			m.CurrentAuthor = ""
			m.CurrentJournal = ""
			m.JournalChoice.ExistingJournals = nil
			m.ViewIdx = LoginIdx

		case key.Matches(msg, MenuKeyMap.Down):
			mp.SelectIdx++
			if mp.SelectIdx > 3 {
				mp.SelectIdx = 3
			}
		case key.Matches(msg, MenuKeyMap.Up):
			mp.SelectIdx--
			if mp.SelectIdx < 0 {
				mp.SelectIdx = 0
			}
		case key.Matches(msg, MenuKeyMap.Select):
			if mp.SelectIdx == 0 {
				m.JournalChoice.ChoiceTi.Focus()
				m.ViewIdx = JournalChoiceIdx
			} else if mp.SelectIdx == 1 {
				if m.CurrentJournal == "" {
					mp.Message = "Please select a journal to use first"
				} else {
					m.WriteEntryPage.TitleTi.Focus()
					m.ViewIdx = NewEntryIdx
				}
			} else if mp.SelectIdx == 2 {
				m.ViewIdx = ViewEntryIdx
				m.ViewEntry.TitleInput.Focus()
			} else {
				m.ViewIdx = ViewJournalIdx
			}
			mp.SelectIdx = 0
		}
	}

	return m, nil
}

func menuView(m model) string {
	mp := m.Menu
	st := "Choose an action\n\n"
	st += selected(mp, 0, "Select a journal to use") + "\n"
	st += selected(mp, 1, "Write a new journal entry") + "\n"
	st += selected(mp, 2, "View a journal entry") + "\n"
	st += selected(mp, 3, "Read through a journal by name") + "\n"
	st += "\n" + mp.Message
	st += helpStyle(menuHelpString())

	return st
}

func menuHelpString() string {
	st := MenuKeyMap.Up.Help().Key + ": "
	st += MenuKeyMap.Up.Help().Desc + ",  "
	st += MenuKeyMap.Down.Help().Key + ": "
	st += MenuKeyMap.Down.Help().Desc + ",  "
	st += MenuKeyMap.Select.Help().Key + ": "
	st += MenuKeyMap.Select.Help().Desc + ",  "
	st += MenuKeyMap.Quit.Help().Key + ": "
	st += MenuKeyMap.Quit.Help().Desc + "\n"

	return st
}
