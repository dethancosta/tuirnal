## To Do List
- [x] migrate to bbolt (boltdb)
- [x] write basic cli
- [x] write tui with bubble-tea
- [ ] !! implement read-journal page
- [ ] !! refactor file creation/storage
	- [x] use github.com/muesli/go-app-paths
- [ ] include command hints at bottom of each page
- [ ] Include date of creation in entryView
- [x] Include list of entries in current journal for entryView
- [ ] Make it prettier and fluid with charmbracelet/lipgloss
- [ ] Switch from bbolt to sqlite
- [ ] rewrite to better handle tags for querying
	- [ ] make separate 'tags' bucket
	- [ ] make separate 'entryTags' bucket to grab entry ids
- [ ] add option to encrypt journal and/or journal entries
- [ ] add ability to create/use templates (with prompts, or for a coffee journal, etc.)

### Long Term Goals...
- [ ] move from json serialization to protocol buffers
- [ ] introduce vector db integration to query & chat with your journal
- [ ] introduce markdown write/display
- [ ] add ability to listen to entries as TTS?
- [ ] add ability to upload handwritten notes with OCR

