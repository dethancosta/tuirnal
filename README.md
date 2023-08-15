<h1>tuirnal</h1>
<h3>A journal app in the terminal</h3>


With tuirnal, you can can write and read journal entries without having to leave the terminal.

Features:
- written in Go, using charmbracelet's bubbletea module for the TUI
- multiple user profiles, each with multiple journals

This project is currently a barebones MVP, with much more functionality on the roadmap.
It should currently only work on Unix-like systems due to the way the storage file is saved; I've only tested it on my Macbook.
You can read the roadmap.txt file in this repository for a to-do list of features.
The tasks with the highest priority right now:
- refactor file creation/storage to be cross-platform (Windows/Mac/Linux)
- better authentication and the option to encrypt journals/entries
- a prettier, more fluid TUI. Most likely with charmbracelet's libgloss library
- better support for searching with and using tags
