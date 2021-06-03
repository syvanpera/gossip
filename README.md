# GOSSIP - Snippet Manager

Command line text snippet/bookmark manager written in Go.

Inspired by
[Pet](https://github.com/knqyf263/pet) and
[Buku](https://github.com/jarun/Buku).

## Features

* Supports generic text snippets, code snippets, shell commands and bookmarks
* Syntax coloring for code snippets
* Running stored shell commands
* Opening bookmarks in browser
* Copying snippet contents to the clipboard
* Automatically fetches description and tags for bookmarks from the HTML in the
  page or from Github API (in case of links pointing to github repos)

## Demo

### Bookmarks
<img src="demo-bookmarks.gif" width="700">

### Code snippets
<img src="demo-code.gif" width="700">

### Commands
<img src="demo-commands.gif" width="700">

## TODO

* Switch to using https://github.com/peterh/liner (instead of github.com/chzyer/readline)
