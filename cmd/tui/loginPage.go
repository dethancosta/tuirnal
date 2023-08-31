package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

type loginModel struct {
	NameTi     textinput.Model
	PasswordTi textinput.Model
	Message    string
	SelectIdx  int
}

var DefaultKeyMap = KeyMap{

	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),

	Down: key.NewBinding(
		key.WithKeys("tab", "down"),
		key.WithHelp("⇥/↓", "move down"),
	),

	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),

	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/ctrl+c", "quit"),
	),
}

func (m loginModel) GetSelection() int {
	return m.SelectIdx
}

func initLogin() loginModel {
	n := textinput.New()
	n.CharLimit = 32
	n.Width = 32 // TODO make flexible?
	n.Prompt = "Username: "
	n.Placeholder = "..."

	//TODO format password input appropriately
	p := textinput.New()
	p.CharLimit = 32
	p.Width = 32
	p.Prompt = "Password: "
	p.Placeholder = "..."
	p.EchoMode = textinput.EchoPassword

	n.Focus()
	return loginModel{
		NameTi:     n,
		PasswordTi: p,
		Message:    "",
		SelectIdx:  0,
	}
}

func updateLogin(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	lm := &m.Login

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			m.ViewEntry.EntriesCache = nil
			return m, tea.Quit

		case key.Matches(msg, DefaultKeyMap.Select):
			name := lm.NameTi.Value()
			pw := lm.PasswordTi.Value()

			if lm.SelectIdx == 2 {
				ok, err := loginAuthor(m.App, name, pw)
				if !ok && err != nil {
					lm.Message = "Author with that username doesn't exist."
				} else if ok {
					lm.Message = ""
					m.ViewIdx = MenuIdx
					m.CurrentAuthor = lm.NameTi.Value()
					lm.NameTi.Reset()
					lm.NameTi.Blur()
					lm.PasswordTi.Reset()
					lm.PasswordTi.Blur()

					journals, err := m.App.JournalModel.GetAllJournals(strings.ToLower(m.CurrentAuthor))
					if err != nil {
						m.JournalChoice.Message = fmt.Sprintf("Couldn't load journals for user:\n%s", err.Error())
					} else {
						m.JournalChoice.SetJournalsCache(journals)
					}
					lm.SelectIdx = 0
				} else {
					lm.Message = "Username or password is incorrect."
				}

			} else if lm.SelectIdx == 3 {
				//create new user
				err := createAuthor(m.App, name, pw)
				if err != nil {
					lm.Message = err.Error()
				} else {
					m.CurrentAuthor = lm.NameTi.Value()
					lm.NameTi.Reset()
					lm.NameTi.Blur()
					lm.PasswordTi.Reset()
					lm.PasswordTi.Blur()

					m.ViewIdx = MenuIdx
				}
			}

		case key.Matches(msg, DefaultKeyMap.Down):
			lm.SelectIdx++
			if lm.SelectIdx > 3 {
				lm.SelectIdx = 3
			}

		case key.Matches(msg, DefaultKeyMap.Up):
			lm.SelectIdx--
			if lm.SelectIdx < 0 {
				lm.SelectIdx = 0
			}

		default:
			if lm.NameTi.Focused() {
				lm.NameTi, cmd = lm.NameTi.Update(msg)
			} else if lm.PasswordTi.Focused() {
				lm.PasswordTi, cmd = lm.PasswordTi.Update(msg)
			}
		}
	default:
		var cmd2 tea.Cmd
		lm.PasswordTi, cmd = lm.PasswordTi.Update(msg)
		lm.NameTi, cmd2 = lm.NameTi.Update(msg)

		return m, tea.Batch(cmd, cmd2)

	}

	if lm.SelectIdx == 0 {
		lm.NameTi.Focus()
		lm.PasswordTi.Blur()
	} else if lm.SelectIdx == 1 {
		lm.PasswordTi.Focus()
		lm.NameTi.Blur()
	} else { // login or signup button
		lm.NameTi.Blur()
		lm.PasswordTi.Blur()
	}

	return m, cmd
}

func helpString() string {
	st := DefaultKeyMap.Up.Help().Key + ": "
	st += DefaultKeyMap.Up.Help().Desc + ",  "
	st += DefaultKeyMap.Down.Help().Key + ": "
	st += DefaultKeyMap.Down.Help().Desc + ",  "
	st += DefaultKeyMap.Select.Help().Key + ": "
	st += DefaultKeyMap.Select.Help().Desc + ",  "
	st += DefaultKeyMap.Quit.Help().Key + ": "
	st += DefaultKeyMap.Quit.Help().Desc + "\n"

	return st
}

func loginView(m model) string {
	lm := &m.Login
	st := "Login/Signup\n\n"
	st += lm.NameTi.View() + "\n\n"
	st += lm.PasswordTi.View() + "\n\n"
	st += selected(lm, 2, "Log In") + "\n\n"
	st += selected(lm, 3, "Sign Up") + "\n\n"
	st += lm.Message
	st += helpString()

	return st
}
