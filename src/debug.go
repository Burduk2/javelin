package main

import (
	"fmt"
	"os"
	"strings"
)

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

const corruptedStoreErrorMsg = `anchor store is corrupted; reset it with "j -R"`

func cliMsg(isErr bool, msg string, args ...any) {
	color := ansiGreen
	if isErr {
		color = ansiRed
	}
	fmt.Fprintln(os.Stderr, style("j:", ansiBold, color), fmt.Sprintf(msg, args...))
	if isErr {
		os.Exit(1)
	}
}

func printUsage() {
	cmds := []struct {
		cmd  string
		desc string
	}{
		{"<name>", "jump to anchor <name>"},
		{"-a <name>", "anchor current directory as <name>"},
		{"-r <name>", "remove anchor <name>"},
		{"-l", "list anchors"},
		{"-R", "reset anchors"},
		{"-h", "show help"},
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
