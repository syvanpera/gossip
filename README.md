# GOSSIP - Bookmark Manager

Command line bookmark manager written in Go.

Heavily inspired by [Buku](https://github.com/jarun/Buku).
The structure follows [go-service-example](https://github.com/kott/go-service-example) pretty closely.

## Usage

List bookmarks

- `gossip ls`

Add bookmark

- `gossip add <url>`

Delete bookmark

- `gossip del <id>`

Set tags

- `gossip edit <id> --tag <tag1>,<tag2>,...`

Add tag

- `gossip edit <id> --tag +<tag>`

Remove tag

- `gossip edit <id> --tag -<tag>`

You can add/remove multiple tags by separating them with a comma:

- `gossip edit <id> --tag +<tag1>,<tag2>`
- `gossip edit <id> --tag -<tag1>,<tag2>`

And even add and remove tags in the same command:

- `gossip edit <id> --tag -<tag1>,<tag2>+<tag2>,<tag3>`

## Features

- Open bookmarks in default browser (or any configured browser)
- Copy bookmark URL to the clipboard
- Can automatically fetch the description and/or tags from following sources:
  - The web page itself (using page metadata and title)
  - Github API (if the link is to a github repo)

## Demo

<img src="demo-bookmarks.gif" width="700">

## TODO

- Add confirmation to delete
- Add/remove tags (+tag/-tag)
