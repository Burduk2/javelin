# Javelin

Javelin is a tiny directory jumper for short anchor names.

It is built for the fast path:
- add the current directory as an anchor
- jump to it later with a 1-2 character name
- keep the command surface small

## Commands

```bash
j <name>       jump to anchor <name>
j -a <name>    anchor current directory as <name>
j -r <name>    remove anchor <name>
j -l           list anchors
j -R           reset anchors
j -h           show help
```

## How It Works

The Go binary resolves anchors and prints the target directory path.

The shell wrapper defines `j()`, calls the binary, and then `cd`s into the returned path. That wrapper is required because a standalone binary cannot change the current shell directory.

Anchors are stored in a small custom text file under `EXE_ROOT/store`.

## Install

### With `just`

This is the current full install flow for Unix shells:

```bash
just install
```

Or choose the shell explicitly: (not yet supported)

```bash
just install zsh
just install bash
just install fish
```

What it does:
- builds `javelin`
- installs it to `~/.local/share/javelin/bin/javelin`
- renders a shell wrapper into `~/.config/javelin/`
- appends a `source ...` line to your shell config if needed

After install, reload your shell config.

Examples:

```bash
source ~/.zshrc
source ~/.bashrc
source ~/.config/fish/config.fish
```

### With `go install` (not yet supported)

If you only want the binary:

```bash
go install github.com/you/javelin@latest
```

`go install` does not install the shell wrapper, so you still need a shell function that:
- runs `javelin`
- captures its output
- calls `cd`

## Build

Build all release binaries:

```bash
just build
```

This writes platform-specific binaries into `bin/`

## Notes

- Anchor names cannot start with `-`.
- Javelin allows up to 32 anchors.
- Short anchor names are recommended for fast use.
