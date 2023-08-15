package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type MenuPage struct {
	SelectIdx int
	Message   string
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

	switch msg.(tea.KeyMsg).String() {
	case "q":
		// TODO sign out
		lp := &m.Login
		lp.NameTi.Reset()
		lp.PasswordTi.Reset()
		lp.NameTi.Focus()
		mp.SelectIdx = 0
		m.CurrentAuthor = ""
		m.CurrentJournal = ""
		m.ViewIdx = LoginIdx
	case "down", "j":
		mp.SelectIdx++
		if mp.SelectIdx > 3 {
			mp.SelectIdx = 3
		}
	case "up", "k":
		mp.SelectIdx--
		if mp.SelectIdx < 0 {
			mp.SelectIdx = 0
		}
	case "enter":
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
		} else {
			m.ViewIdx = ViewJournalIdx
		}
		mp.SelectIdx = 0
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

	return st
}
