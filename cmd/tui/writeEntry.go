package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var WeKeyMap = struct {
	Up   key.Binding
	Down key.Binding
	Save key.Binding
	Quit key.Binding
}{
	Up: key.NewBinding( //up
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),

	Down: key.NewBinding( //down
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),

	Save: key.NewBinding( //ctrl+s
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),

	Quit: key.NewBinding( //esc
		key.WithKeys("esc"),
		key.WithHelp("esc", "quit"),
	),
}

type WriteEntry struct {
	TitleTi   textinput.Model
	EntryTa   textarea.Model
	TagsTi    textinput.Model
	SelectIdx int
	Message   string
}

func initWriteEntry() WriteEntry {

	titleInput := textinput.New()
	titleInput.CharLimit = 64
	titleInput.Width = 50
	titleInput.Prompt = ""
	titleInput.Placeholder = "Give your entry a title..."

	entryTa := textarea.New()
	entryTa.Placeholder = "Write your entry here..."
	entryTa.Prompt = "┃ "
	entryTa.CharLimit = 0 // no limit to input length
	entryTa.ShowLineNumbers = false

	entryTa.SetWidth(50)  // TODO adjust as needed (make fluid?)
	entryTa.SetHeight(15) // TODO adjust as needed (make fluid?)

	tagsTi := textinput.New()
	tagsTi.Width = 50
	tagsTi.Prompt = "Tags: "
	tagsTi.Placeholder = "Tags are separated by spaces"

	return WriteEntry{
		TitleTi:   titleInput,
		EntryTa:   entryTa,
		TagsTi:    tagsTi,
		SelectIdx: 0,
		Message:   "",
	}
}

func updateWriteEntry(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	wep := &m.WriteEntryPage

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, WeKeyMap.Quit):
			m.ViewIdx = MenuIdx
			wep.TitleTi.Reset()
			wep.TitleTi.Blur()
			wep.EntryTa.Reset()
			wep.EntryTa.Blur()
			wep.TagsTi.Reset()
			wep.TagsTi.Blur()
			wep.Message = ""
			wep.SelectIdx = 0

		case key.Matches(msg, WeKeyMap.Save):
			if len(strings.TrimSpace(wep.TitleTi.Value())) == 0 {
				wep.Message = "Please enter a title for the entry."
			} else {
				ok := entryNameAvailable(
					m.App,
					m.CurrentAuthor,
					m.CurrentJournal,
					wep.TitleTi.Value(),
				)

				if ok {
					err := saveJournalEntry(
						m.App,
						m.CurrentAuthor,
						m.CurrentJournal,
						wep.TitleTi.Value(),
						wep.TagsTi.Value(),
						wep.EntryTa.Value(),
					)
					if err != nil {
						wep.Message = "Couldn't save the entry :("
					} else {
						wep.Message = "Entry saved successfully."
						newCacheEntry, err := m.App.EntryModel.Get(m.CurrentAuthor, m.CurrentJournal, wep.TitleTi.Value())
						if err != nil {
							wep.Message += " Couldn't update cache."
						}
						m.ViewEntry.SetCache(append(m.ViewEntry.EntriesCache, newCacheEntry))
					}
				} else {
					wep.Message = "This title is already taken."
				}
			}

		case key.Matches(msg, WeKeyMap.Down):
			wep.SelectIdx++
			if wep.SelectIdx > 2 {
				wep.SelectIdx = 2
			}

		case key.Matches(msg, WeKeyMap.Up):
			wep.SelectIdx--
			if wep.SelectIdx < 0 {
				wep.SelectIdx = 0
			}

		default:
			if wep.TitleTi.Focused() {
				wep.TitleTi, cmd = wep.TitleTi.Update(msg)
			} else if wep.EntryTa.Focused() {
				wep.EntryTa, cmd = wep.EntryTa.Update(msg)
			} else {
				wep.TagsTi, cmd = wep.TagsTi.Update(msg)
			}
		}
	default:
		cmds := make([]tea.Cmd, 3)
		wep.EntryTa, cmds[0] = wep.EntryTa.Update(msg)
		wep.TitleTi, cmds[1] = wep.TitleTi.Update(msg)
		wep.TagsTi, cmds[2] = wep.TagsTi.Update(msg)

		return m, tea.Batch(cmds...)
	}

	if wep.SelectIdx == 0 {
		wep.TitleTi.Focus()
		wep.EntryTa.Blur()
		wep.TagsTi.Blur()
	} else if wep.SelectIdx == 1 {
		wep.EntryTa.Focus()
		wep.TitleTi.Blur()
		wep.TagsTi.Blur()
	} else {
		wep.TagsTi.Focus()
		wep.EntryTa.Blur()
		wep.TitleTi.Blur()
	}
	return m, cmd
}

func writeEntryView(m model) string {
	wep := &m.WriteEntryPage
	st := "Journal Entry\n"
	st += wep.TitleTi.View() + "\n"
	st += wep.EntryTa.View() + "\n"
	st += wep.TagsTi.View() + "\n\n"
	st += wep.Message
	st += helpStyle(weHelpString())

	return st
}

func weHelpString() string {
	st := WeKeyMap.Up.Help().Key + ": "
	st += WeKeyMap.Up.Help().Desc + ",  "
	st += WeKeyMap.Down.Help().Key + ": "
	st += WeKeyMap.Down.Help().Desc + ",  "
	st += WeKeyMap.Save.Help().Key + ": "
	st += WeKeyMap.Save.Help().Desc + ",  "
	st += WeKeyMap.Quit.Help().Key + ": "
	st += WeKeyMap.Quit.Help().Desc + "\n"

	return st
}
