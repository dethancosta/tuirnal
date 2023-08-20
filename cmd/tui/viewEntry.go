package main

import (
	"strings"

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
		if mt, ok := msg.(tea.KeyMsg); ok {
			kmsg := mt.String()
			switch kmsg {
			case "j", "down":
				vem.Vp.Update(vem.Vp.KeyMap.HalfPageDown)
			case "k", "up":
				vem.Vp.Update(vem.Vp.KeyMap.HalfPageUp)
			case "q":
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
		if kt, ok := msg.(tea.KeyMsg); ok {
			switch kt.String() {
			case "enter":
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
			case "ctrl+c":
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
		return st
	} else {
		searchList := vem.getSearchList()

		return vem.TitleInput.View() + "\n\n" + vem.Message + "\n\n" + searchList
	}
}

func (vem *viewEntryModel) getSearchList() string {
	titles := make([]string, len(vem.EntriesCache))
	for i, e := range vem.EntriesCache {
		titles[i] = e.Title
	}

	sb := strings.Builder{}
	search := strings.ToLower(vem.TitleInput.Value())
	for i := range titles {
		tLower := strings.ToLower(titles[i])
		if strings.HasPrefix(tLower, search) {
			//TODO get appropriate entry from db
			sb.WriteString(titles[i] + "\n")
		}
	}
	return sb.String()
}

func (vem *viewEntryModel) SetCache(cache []*models.JournalEntry) {
	vem.EntriesCache = cache
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
