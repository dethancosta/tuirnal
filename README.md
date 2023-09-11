<h1>tuirnal</h1>
<h3>A journal app in the terminal</h3>


With tuirnal, you can can write and read journal entries without having to leave the terminal.

Features:
- written in Go, using charmbracelet's [bubbletea](https://github.com/charmbracelet/bubbletea) for the TUI
- multiple user profiles, each with multiple journals

This project is currently a barebones MVP, with more functionality on the roadmap.
You can read the roadmap.txt file in this repository for a to-do list of features.
The tasks with the highest priority right now:
- better authentication and the option to encrypt journals/entries
- a prettier, more fluid TUI. Most likely with [lipgloss](https://github.com/charmbracelet/lipgloss)
- better support for searching with and using tags
