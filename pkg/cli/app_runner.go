package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/content"
)

// Selection represents the content and persistence hooks for a typing session.
type Selection struct {
	Text     string
	Content  *content.Content
	Provider StateProvider
}

// AppConfig coordinates shared CLI behavior between binaries.
type AppConfig struct {
	Name            string
	Version         string
	ListDescription string
	ListItems       func() ([]string, error)
	Configure       []func() error
	SelectAndLoad   func(width, height int) (*Selection, error)
	Width           int
	Height          int
	Args            []string
	Stdout          io.Writer
	Stderr          io.Writer
}

// RunApp wires common flag handling, listing, and program execution for CLI binaries.
func RunApp(cfg AppConfig) error {
	if cfg.SelectAndLoad == nil {
		return fmt.Errorf("cli: SelectAndLoad handler is required")
	}
	if cfg.ListItems == nil {
		return fmt.Errorf("cli: ListItems handler is required")
	}

	args := cfg.Args
	if args == nil {
		args = os.Args[1:]
	}

	stdout := cfg.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}

	stderr := cfg.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	width := cfg.Width
	if width <= 0 {
		width = 80
	}
	height := cfg.Height
	if height <= 0 {
		height = 24
	}

	fs := flag.NewFlagSet(cfg.Name, flag.ContinueOnError)
	fs.SetOutput(stderr)

	listDesc := cfg.ListDescription
	if listDesc == "" {
		listDesc = "List available entries"
	}

	list := fs.Bool("l", false, listDesc)
	listLong := fs.Bool("list", false, listDesc+" (long form)")
	version := fs.Bool("v", false, "Show application version")
	versionLong := fs.Bool("version", false, "Show application version (long form)")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	if *version || *versionLong {
		_, _ = fmt.Fprintln(stdout, cfg.Version)
		return nil
	}

	for _, configure := range cfg.Configure {
		if configure == nil {
			continue
		}
		if err := configure(); err != nil {
			return err
		}
	}

	if *list || *listLong {
		items, err := cfg.ListItems()
		if err != nil {
			return err
		}
		for _, item := range items {
			_, _ = fmt.Fprintln(stdout, item)
		}
		return nil
	}

	selection, err := cfg.SelectAndLoad(width, height)
	if err != nil {
		return err
	}
	if selection == nil {
		return nil
	}
	if selection.Content == nil {
		return fmt.Errorf("cli: selection missing content metadata")
	}
	if selection.Provider == nil {
		return fmt.Errorf("cli: selection missing state provider")
	}

	model := NewModel(selection.Text, selection.Content, width, height, selection.Provider)
	program := tea.NewProgram(model)
	_, err = program.Run()
	return err
}
