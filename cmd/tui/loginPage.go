package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type loginModel struct {
	NameTi     textinput.Model
	PasswordTi textinput.Model
	Message    string
	SelectIdx  int
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
	if _, ok := msg.(tea.KeyMsg); !ok {
		var cmd2 tea.Cmd
		lm.PasswordTi, cmd = lm.PasswordTi.Update(msg)
		lm.NameTi, cmd2 = lm.NameTi.Update(msg)

		return m, tea.Batch(cmd, cmd2)
	}

	switch msg.(tea.KeyMsg).String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "enter":
		name := lm.NameTi.Value()
		pw := lm.PasswordTi.Value()

		if lm.SelectIdx == 2 {
			ok, err := loginAuthor(m.App, name, pw)
			if !ok && err != nil {
				lm.Message = "Author with that username doesn't exist."
			} else if ok {
				m.ViewIdx = MenuIdx
				m.CurrentAuthor = lm.NameTi.Value()
				lm.NameTi.Reset()
				lm.NameTi.Blur()
				lm.PasswordTi.Reset()
				lm.PasswordTi.Blur()
				lm.SelectIdx = 0
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
	case "tab":
		lm.SelectIdx++
		if lm.SelectIdx > 3 {
			lm.SelectIdx = 3
		}
	case "down":
		lm.SelectIdx++
		if lm.SelectIdx > 3 {
			lm.SelectIdx = 3
		}
	case "up":
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

func loginView(m model) string {
	lm := &m.Login
	st := "Login/Signup\n\n"
	st += lm.NameTi.View() + "\n\n"
	st += lm.PasswordTi.View() + "\n\n"
	st += selected(lm, 2, "Log In") + "\n\n"
	st += selected(lm, 3, "Sign Up")
	st += lm.Message

	return st
}
