package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const version = "0.0.9"
const LIMIT = 32

var binPath = os.Getenv("EXE_ROOT")
var storePath = binPath + "/store.json"

const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiItalic = "\x1b[3m"
	ansiGreen  = "\x1b[32m"
	ansiRed    = "\x1b[31m"
	ansiYellow = "\x1b[33m"
	ansiGray   = "\x1b[90m"
)

func style(text string, codes ...string) string {
	return strings.Join(codes, "") + text + ansiReset
}

const corruptedStoreErrorMsg = `anchor store is corrupted
  -> "j -R" to reset it
  -> "j -h" for more info`

func cliMsg(isErr bool, msg string, args ...any) {
	color := ansiGreen
	if isErr {
		color = ansiRed
	}
	fmt.Fprintln(os.Stderr, style("j:", ansiBold, color), fmt.Sprintf(msg, args...))
}

func printUsage() {
	cmds := []struct {
		cmd  string
		desc string
	}{
		{"<name>", "jump to anchor <name>"},
		{"-a <name>", "anchor the current directory as <name>"},
		{"-r <name>", "remove anchor <name>"},
		{"-l", "list anchors"},
		{"-h", "show help"},
		{"-R", "reset the anchor store"},
	}
	topCmdLen := 0
	for _, cmd := range cmds {
		if len(cmd.cmd) > topCmdLen {
			topCmdLen = len(cmd.cmd)
		}
	}
	fmt.Fprintf(
		os.Stderr,
		"%s %s\n\n%s\n",
		style("Javelin", ansiBold),
		style(fmt.Sprintf("[v%s]", version), ansiGray),
		"usage:",
	)
	for _, cmd := range cmds {
		pad := topCmdLen - len(cmd.cmd) + 4
		fmt.Fprintf(
			os.Stderr,
			"  %s %s%s%s\n",
			style("j", ansiBold),
			style(cmd.cmd, ansiYellow),
			strings.Repeat(" ", pad),
			style(cmd.desc, ansiGray),
		)
	}
}

type Store struct {
	Version string   `json:"version"`
	Anchors []Anchor `json:"anchors"`
}
type Anchor struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func getStore() *Store {
	data, err := os.ReadFile(storePath)
	store := &Store{}
	if err != nil {
		cliMsg(true, corruptedStoreErrorMsg)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &store)
	if err != nil {
		cliMsg(true, corruptedStoreErrorMsg)
		os.Exit(1)
	}

	return store
}

func (store *Store) addAnchor(name string) {
	if len(store.Anchors) == LIMIT {
		cliMsg(true, "anchor limit reached (%d)\n  -> remove one with \"j -r <name>\"\n  -> or reset all with \"j -R\"", LIMIT)
		os.Exit(1)
	}
	absPath, err := filepath.Abs(".")
	if err != nil {
		cliMsg(true, "could not resolve current directory")
		os.Exit(1)
	} else if name[0] == '-' {
		cliMsg(true, "anchor name cannot start with a dash")
		os.Exit(1)
	}
	for _, anchor := range store.Anchors {
		if anchor.Name == name {
			cliMsg(true, "anchor with name %q (%q) already exists", name, absPath)
			os.Exit(1)
		}
		if anchor.Path == absPath {
			cliMsg(true, "anchor with path %q (%q) already exists", absPath, anchor.Name)
			os.Exit(1)
		}
	}

	store.Anchors = append(store.Anchors, Anchor{Name: name, Path: absPath})
	cliMsg(false, "added anchor %q (%q)", name, absPath)
}

func (store *Store) removeAnchor(name string) {
	for i, anchor := range store.Anchors {
		if anchor.Name == name {
			store.Anchors = append(store.Anchors[:i], store.Anchors[i+1:]...)
			cliMsg(false, "removed anchor %q (%q)", name, anchor.Path)
			return
		}
	}
	cliMsg(true, "anchor %q does not exist", name)
	os.Exit(1)
}

func (store *Store) list() {
	topNameLen := 0
	for _, anchor := range store.Anchors {
		if len(anchor.Name) > topNameLen {
			topNameLen = len(anchor.Name)
		}
	}
	for _, anchor := range store.Anchors {
		padAmt := topNameLen - len(anchor.Name) + 2
		pad := style(fmt.Sprintf("%s➜", strings.Repeat("–", padAmt-1)), ansiGray)
		fmt.Fprintf(os.Stderr, "%s%s%s\n", style(anchor.Name, ansiBold, ansiYellow), pad, anchor.Path)
	}
}

func (store *Store) save() {
	store.Version = version
	data, err := json.MarshalIndent(store, "", "\t")
	if err != nil {
		cliMsg(true, corruptedStoreErrorMsg)
		os.Exit(1)
	}
	err = os.WriteFile(storePath, data, 0644)
	if err != nil {
		cliMsg(true, "could not save store: %v", err)
		os.Exit(1)
	}
}

func main() {
	args := os.Args[1:]
	numArgs := len(args)

	if numArgs == 0 {
		printUsage()
		os.Exit(0)
	}

	cmd := args[0]
	cmdArgs := args[1:]
	numCmdArgs := len(cmdArgs)
	switch cmd {
	case "-h":
		printUsage()
	case "-a":
		store := getStore()
		name := ""
		switch numCmdArgs {
		case 1:
			name = cmdArgs[0]
		case 0:
			cliMsg(true, "provide an anchor name")
			os.Exit(1)
		default:
			cliMsg(true, "unexpected extra arguments: %s", strings.Join(cmdArgs[1:], " "))
			os.Exit(1)
		}

		store.addAnchor(name)
		store.save()
	case "-r":
		store := getStore()
		if numCmdArgs != 1 {
			cliMsg(true, "provide an anchor name")
			os.Exit(1)
		}
		name := cmdArgs[0]
		store.removeAnchor(name)
		store.save()
	case "-l":
		store := getStore()
		store.list()
	case "-R":
		if confirm("Reset anchor store (this will delete all anchors)?") {
			store := Store{Anchors: []Anchor{}}
			store.save()
			cliMsg(false, "reset anchor store")
		}
	default:
		// args is at least 1 long, checked above
		if cmd[0] == '-' {
			cliMsg(true, "invalid command: %q", cmd)
			os.Exit(1)
		}

		if numArgs > 1 {
			cliMsg(true, "unexpected extra arguments: %s", strings.Join(cmdArgs, " "))
			os.Exit(1)
		}
		store := getStore()
		name := args[0]
		for _, anchor := range store.Anchors {
			if anchor.Name == name {
				fmt.Print(anchor.Path)
				os.Exit(0)
			}
		}
		cliMsg(true, "anchor %q does not exist", name)
		os.Exit(1)
	}

	os.Exit(0)
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprintf(os.Stderr, "%s [y/n]: ", prompt)

		input, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Fprintln(os.Stderr, "Please enter y or n.")
		}
	}
}
