## To Do List
- [x] migrate to bbolt (boltdb)
- [x] write basic cli
- [x] write tui with bubble-tea
- [x] refactor file creation/storage
- [ ] format entry-view list so that dates are aligned
- [ ] !! implement read-journal page
- [x] !! include command hints at bottom of each page
- [x] Include date of creation in entryView
- [x] Include list of entries in current journal for entryView
- [ ] Check for and deny blank titles/names
- [ ] Make it prettier and resizable with charmbracelet/lipgloss
- [ ] Switch from bbolt to sqlite
- [ ] rewrite to better handle tags for querying
	- [ ] ? migrate to sqlite
	- [ ] make separate 'tags' bucket
	- [ ] make separate 'entryTags' bucket to grab entry ids
- [ ] add ability to create/use templates (with prompts, or for a coffee journal, etc.)
- [ ] change page message to be between top status bar and main content, with coloured text
- [x] !! encrypt passwords

### Long Term Goals...
- [ ] introduce vector db integration to query & chat with your journal
- [ ] introduce markdown write/display
- [ ] add ability to listen to entries as TTS?
- [ ] add ability to upload handwritten notes with OCR
