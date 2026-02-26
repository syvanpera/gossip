<div align="center">
    <img src="./docs/gossip-logo.png" alt="gossiplogo" width="400" height="400" />
</div>

# gossip - A Command Line Bookmark Manager

`gossip` is a fast, lightweight, and gorgeous command-line bookmark manager written in Go.
It allows you to save, tag, search, and open links entirely from your terminal. 

## ✨ Features

* **Beautiful UI:** Styled with Charm's Lip Gloss for easy-to-read, color-coded terminal output.
* **Smart Fetching:** Automatically grabs titles and descriptions from HTML meta tags.
* **GitHub Integration:** Recognizes GitHub URLs and automatically fetches repository descriptions, languages, star counts, and topics (as tags).
* **Regex Search:** Search your bookmarks by keyword, URL, or tag with bright visual highlighting.
* **Shell Autocompletion:** Press `TAB` to auto-complete bookmark IDs and view titles directly in your shell prompt.
* **Cross-Platform:** Opens links natively on macOS, Windows, and Linux.
* **Simple Storage:** Data is stored in a clean, human-readable JSON file (`~/.config/gossip/bookmarks.json`).

## 🚀 Installation

Ensure you have [Go](https://go.dev/) installed on your system. 

Clone this repository and run `go install` from the root directory:

```bash
git clone [https://github.com/syvanpera/gossip.git](https://github.com/syvanpera/gossip.git)
cd gossip
go install

```

*Note: Make sure your Go `bin` directory (e.g., `~/go/bin`) is in your system's `PATH`.*

## 📖 Usage

### Adding Bookmarks

Add a URL. The CLI will show a spinner while it asynchronously fetches the page title and description.

```bash
gossip add [https://go.dev](https://go.dev)

```

Add a URL with custom tags:

```bash
gossip add [https://github.com/spf13/cobra](https://github.com/spf13/cobra) -t cli,go,tutorial

```

### Listing Bookmarks

View all saved bookmarks, formatted beautifully.

```bash
gossip list

```

### Searching

Search across URLs, titles, tags, and comments. The matching text will be highlighted in the output.

```bash
gossip search "cli"

```

### Opening Links

Open a bookmark in your default web browser using its ID.

```bash
gossip open a1b2c3d4

```

### Editing Bookmarks

Modify specific fields of an existing bookmark:

```bash
gossip edit a1b2c3d4 --title "New Title" -t "new,tags" -c "Updated comment"

```

You can also force the CLI to refetch the data from the website if it has changed:

```bash
gossip edit a1b2c3d4 --refetch-tags  # Updates only the tags from the source
gossip edit a1b2c3d4 --refetch-all   # Updates title, comment, and tags

```

### Deleting

Remove a bookmark permanently by its ID.

```bash
gossip delete a1b2c3d4

```

## 🪄 Enabling Autocompletion

`gossip` supports dynamic shell autocompletion for bookmark IDs. To enable it, load the completion script into your shell.

**For Fish:**

```bash
gossip completion fish | source

```

**For Zsh:**

```bash
source <(gossip completion zsh)

```

**For Bash:**

```bash
source <(gossip completion bash)

```

## 🛠️ Built With

* [Cobra](https://github.com/spf13/cobra) - CLI framework
* [Viper](https://github.com/spf13/viper) - Configuration management
* [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Bubbles](https://github.com/charmbracelet/bubbles) - Terminal UI and Spinners
* [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
* [Go GitHub](https://github.com/google/go-github) - GitHub API integration
