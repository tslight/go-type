package main

import (
"flag"
"fmt"
"os"

tea "github.com/charmbracelet/bubbletea"
"github.com/tobe/go-type/internal/godocgen"
"github.com/tobe/go-type/pkg/cli"
)

// doctype - typing practice app based on Go documentation
var Version = "unknown"

func main() {
	listDocs := flag.Bool("l", false, "List available Go documentation modules")
	docList := flag.Bool("list", false, "List available Go documentation modules (long form)")
	version := flag.Bool("version", false, "Show application version")
	flag.Parse()

	if *version {
		fmt.Println(Version)
		return
	}

	// Handle list docs flag
	if *listDocs || *docList {
		docs := godocgen.GetDocumentationNames()
		for _, doc := range docs {
			fmt.Println(doc)
		}
		os.Exit(0)
	}

	// For now, just get a random documentation
	// TODO: Add menu for selecting specific documentation
	text, err := godocgen.GetRandomDocumentation()
	if err != nil {
		fmt.Printf("Error loading documentation: %v\n", err)
		os.Exit(1)
	}

	// Create and run the Bubble Tea model for typing test
	m := cli.NewModel(text, nil, 80, 24)
	p := tea.NewProgram(m)

	_, err = p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
