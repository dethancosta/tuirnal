package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// TODO refactor into struct with init function

func updateEntries(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	//TODO implement
	if mt, ok := msg.(tea.KeyMsg); ok {
		kmsg := mt.String()
		switch kmsg {
		case "q":
			m.ViewIdx = MenuIdx
		}
	}
	return m, nil
}

func entriesView(m model) string {
	//TODO implement
	return "🏗️ This page is under construction 🏗️"
}
