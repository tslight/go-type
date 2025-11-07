package runner

import (
	"flag"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/model"
	"github.com/tobe/go-type/internal/selection"
)

// AppConfig coordinates shared CLI behavior between binaries.
type AppConfig struct {
	Name            string
	Version         string
	ListDescription string
	ListItems       func() ([]string, error)
	Configure       []func() error
	SelectAndLoad   func(width, height int) (*selection.Selection, error)
	Width           int
	Height          int
	Args            []string
	Stdout          io.Writer
	Stderr          io.Writer
}

// RunApp wires common flag handling, listing, and program execution for CLI binaries.
func RunApp(cfg AppConfig) error {
	if cfg.SelectAndLoad == nil {
		return fmt.Errorf("runner: SelectAndLoad handler is required")
	}
	if cfg.ListItems == nil {
		return fmt.Errorf("runner: ListItems handler is required")
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

	selectionResult, err := cfg.SelectAndLoad(width, height)
	if err != nil {
		return err
	}
	if selectionResult == nil {
		return nil
	}
	if selectionResult.Content == nil {
		return fmt.Errorf("runner: selection missing content metadata")
	}
	if selectionResult.Provider == nil {
		return fmt.Errorf("runner: selection missing state provider")
	}

	modelInstance := model.NewModel(selectionResult.Text, selectionResult.Content, width, height, selectionResult.Provider)
	_, err = runModelProgram(modelInstance)
	return err
}

// runModelProgram is a hook so tests can stub Bubble Tea run for model execution.
var runModelProgram = func(m tea.Model) (tea.Model, error) { return tea.NewProgram(m).Run() }
