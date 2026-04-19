package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const version = "0.0.16"
const LIMIT = 32

var binPath = os.Getenv("EXE_ROOT")
var storePath = binPath + "/store"

type Store []Anchor
type Anchor struct {
	name string
	path string
}

// json perf = 130 || my perf = 67
func getStore() *Store {
	rawData, err := os.ReadFile(storePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Store{}
		}
		cliMsg(true, corruptedStoreErrorMsg+"\n"+err.Error())
	}
	data := string(rawData)
	if data == "" {
		return &Store{}
	}

	store := Store{}
	for _, line := range strings.Split(data, "\n") {
		if line == "" {
			continue
		}
		name, path, ok := strings.Cut(line, " ")
		if !ok {
			cliMsg(true, corruptedStoreErrorMsg)
		}
		store = append(store, Anchor{name: name, path: path})
	}
	return &store
}

func (store *Store) save() {
	var b bytes.Buffer
	for _, anchor := range *store {
		b.WriteString(anchor.name)
		b.WriteByte(' ')
		b.WriteString(anchor.path)
		b.WriteByte('\n')
	}

	err := os.WriteFile(storePath, b.Bytes(), 0644)
	if err != nil {
		cliMsg(true, "could not save anchors: %v", err)
	}
}

func getAnchorName(cmdArgs []string) string {
	numArgs := len(cmdArgs)
	name := ""
	switch numArgs {
	case 0:
		cliMsg(true, "anchor name expected")
	case 1:
		name = cmdArgs[0]
	default:
		cliMsg(true, "unexpected extra arguments: %s", strings.Join(cmdArgs[1:], " "))
	}

	return name
}

func (store *Store) addAnchor(name string) {
	anchors := *store
	if len(anchors) >= LIMIT {
		cliMsg(true, "anchor limit reached (%d)", LIMIT)
	}
	absPath, err := filepath.Abs(".")
	if err != nil {
		cliMsg(true, "could not resolve current directory")
	} else if name[0] == '-' {
		cliMsg(true, "anchor name cannot start with a dash")
	} else if len(name) > 6 {
		cliMsg(false,
			"%s prefer short anchor names for faster jumping",
			style("warn:", ansiYellow),
		)
	}
	for _, anchor := range anchors {
		if anchor.name == name {
			cliMsg(true, "anchor %q already exists", name)
		} else if anchor.path == absPath {
			cliMsg(true, "directory %q is already anchored as %q", absPath, anchor.name)
		}
	}

	*store = append(anchors, Anchor{name: name, path: absPath})
	cliMsg(false,
		"added anchor: %s%s%s",
		style(name, ansiBold, ansiYellow),
		style("-✓->", ansiGreen),
		style(absPath),
	)
}

func (store *Store) removeAnchor(name string) {
	anchors := *store
	for i, anchor := range anchors {
		if anchor.name == name {
			*store = append(anchors[:i], anchors[i+1:]...)
			cliMsg(false,
				"removed anchor: %s%s%s",
				style(name, ansiBold, ansiYellow),
				style("-×->", ansiRed),
				style(anchor.path),
			)
			return
		}
	}
	cliMsg(true, "anchor %q does not exist", name)
}

func (store *Store) list() {
	anchors := *store
	topNameLen := 0
	for _, anchor := range anchors {
		if len(anchor.name) > topNameLen {
			topNameLen = len(anchor.name)
		}
	}
	for _, anchor := range anchors {
		padAmt := topNameLen - len(anchor.name) + 2
		pad := style(fmt.Sprintf("%s➜", strings.Repeat("–", padAmt-1)), ansiGray)
		fmt.Fprintf(os.Stderr, "%s%s%s\n", style(anchor.name, ansiBold, ansiYellow), pad, anchor.path)
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
	switch cmd {
	case "-h":
		printUsage()
	case "-a":
		name := getAnchorName(cmdArgs)
		store := getStore()
		store.addAnchor(name)
		store.save()
	case "-r":
		name := getAnchorName(cmdArgs)
		store := getStore()
		store.removeAnchor(name)
		store.save()
	case "-l":
		getStore().list()
	case "-R":
		if confirm("Reset all anchors?") {
			(&Store{}).save()
			cliMsg(false, "reset anchors")
		}
	default:
		if cmd[0] == '-' {
			cliMsg(true, "invalid command: %q", cmd)
		}

		if numArgs > 1 {
			cliMsg(true, "unexpected extra arguments: %s", strings.Join(cmdArgs, " "))
		}
		anchors := *getStore()
		name := args[0]
		for _, anchor := range anchors {
			if anchor.name == name {
				info, err := os.Stat(anchor.path)
				if err != nil || !info.IsDir() {
					cliMsg(true, "could not access directory %q", anchor.path)
				}
				fmt.Print(anchor.path)
				os.Exit(0)
			}
		}
		cliMsg(true, "anchor %q does not exist", name)
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
