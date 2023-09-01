package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dethancosta/tuirnal/internal/models"
	"github.com/muesli/reflow/wordwrap"
)

type viewEntryModel struct {
	// CreatedAt string
	Tags         string
	Title        string
	Vp           viewport.Model
	TitleInput   textinput.Model
	Message      string
	ReadingMode  bool
	EntriesCache []*models.JournalEntry
}

var VeReadingKeyMap = KeyMap{

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

var VeSelectKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),

	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
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

func initViewEntry() viewEntryModel {
	vp := viewport.New(80, 10) // TODO make flexible
	ti := textinput.New()
	ti.Prompt = "Title: "
	return viewEntryModel{
		Tags:         "",
		Title:        "",
		Vp:           vp,
		TitleInput:   ti,
		Message:      "",
		ReadingMode:  false,
		EntriesCache: nil,
	}
}

func updateViewEntry(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	vem := &m.ViewEntry
	var cmd tea.Cmd
	if vem.ReadingMode {
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch {

			case key.Matches(msg, VeReadingKeyMap.Down):
				vem.Vp.Update(vem.Vp.KeyMap.HalfPageDown)
			case key.Matches(msg, VeReadingKeyMap.Up):
				vem.Vp.Update(vem.Vp.KeyMap.HalfPageUp)
			case key.Matches(msg, VeReadingKeyMap.Quit):
				vem.TitleInput.Focus()
				vem.Title = ""
				vem.Tags = ""
				vem.Vp.SetContent("")
				vem.ReadingMode = false
			}
		}
		vem.Vp, cmd = vem.Vp.Update(msg)
		return m, cmd

	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			//TODO implement up and down to select entries in list
			case key.Matches(msg, VeSelectKeyMap.Select):
				exists := !entryNameAvailable(
					m.App,
					m.CurrentAuthor,
					m.CurrentJournal,
					vem.TitleInput.Value(),
				)
				if exists {
					vem.Message = ""
					je, err := m.App.EntryModel.Get(m.CurrentAuthor, m.CurrentJournal, vem.TitleInput.Value())
					if err != nil {
						vem.Message = "ERR: Couldn't get entry from storage."
						vem.TitleInput.Reset()
						vem.TitleInput.Focus()
					}
					vem.TitleInput.Blur()
					vem.Title = je.Title
					vem.Vp.SetContent(wordwrap.String(je.Content, vem.Vp.Width))
					vem.Tags = strings.Join(je.Tags, ", ")
					vem.ReadingMode = true
				} else {
					vem.Message = "No entry with that title."
				}

			case key.Matches(msg, VeSelectKeyMap.Quit):
				vem.TitleInput.Blur()
				m.ViewIdx = MenuIdx
			}
		}
		vem.TitleInput, cmd = vem.TitleInput.Update(msg)
	}
	return m, cmd
}

func viewEntryView(m model) string {
	vem := &m.ViewEntry
	if vem.ReadingMode {
		st := vem.Title + "\n\n" +
			vem.Vp.View() + "\n\n" +
			vem.Tags
		return st + helpStyle(m.ViewEntry.veHelpString())
	} else {
		searchList := vem.getSearchList()

		return vem.TitleInput.View() + "\n\n" + vem.Message + "\n\n" + searchList + helpStyle(m.ViewEntry.veHelpString())
	}
}

func (vem *viewEntryModel) getSearchList() string {
	titles := make([]string, len(vem.EntriesCache))
	for i := range vem.EntriesCache {
		titles[i] = vem.EntriesCache[i].Title + "\t" + vem.EntriesCache[i].WrittenAt.Format(time.DateTime)
	}

	sb := strings.Builder{}
	//search := strings.ToLower(vem.TitleInput.Value())
	for i := range titles {
		//tLower := strings.ToLower(titles[i])
		//if strings.HasPrefix(tLower, search) {
		//TODO get appropriate entry from db
		sb.WriteString(titles[i] + "\n")
		//}
	}
	return sb.String()
}

func (vem *viewEntryModel) SetCache(cache []*models.JournalEntry) {
	vem.EntriesCache = cache
}

func (vem viewEntryModel) veHelpString() string {
	var st string

	if vem.ReadingMode {
		st += VeReadingKeyMap.Up.Help().Key + ": "
		st += VeReadingKeyMap.Up.Help().Desc + ",  "
		st += VeReadingKeyMap.Down.Help().Key + ": "
		st += VeReadingKeyMap.Down.Help().Desc + ",  "
		st += VeReadingKeyMap.Quit.Help().Key + ": "
		st += VeReadingKeyMap.Quit.Help().Desc + "\n"
	} else {
		st += VeSelectKeyMap.Up.Help().Key + ": "
		st += VeSelectKeyMap.Up.Help().Desc + ",  "
		st += VeSelectKeyMap.Down.Help().Key + ": "
		st += VeSelectKeyMap.Down.Help().Desc + ",  "
		st += VeSelectKeyMap.Select.Help().Key + ": "
		st += VeSelectKeyMap.Select.Help().Desc + ",  "
		st += VeSelectKeyMap.Quit.Help().Key + ": "
		st += VeSelectKeyMap.Quit.Help().Desc + "\n"
	}

	return st
}

// for testing purposes
const fillerText = `Lorem ipsum dolor sit amet, 
consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut 
labore et dolore magna aliqua. Ut 
enim ad minim veniam, quis nostrud 
exercitation ullamco laboris nisi ut 
aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit 
in voluptate velit esse cillum dolore
eu fugiat nulla pariatur. Excepteur 
sint occaecat cupidatat non proident,
sunt in culpa qui officia deserunt 
mollit anim id est laborum.`
