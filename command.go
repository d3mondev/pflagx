// Package pflagx provides extensions to spf13/pflag for better flag organization and help text formatting.
package pflagx

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/spf13/pflag"
)

const (
	// DefaultIndentation specifies the default number of spaces
	// to indent flag groups in help output.
	DefaultIndentation = 2

	// DefaultPadding specifies the default minimum number of spaces
	// between flag names and their usage text.
	DefaultPadding = 4

	// DefaultSortFlags determines whether flags are sorted
	// alphabetically by default in help output.
	DefaultSortFlags = false

	// DefaultAlignUsagePerSection determines whether usage text
	// alignment is calculated per section (true) or globally (false)
	// by default.
	DefaultAlignUsagePerSection = false
)

// CommandLine manages multiple FlagSets and provides unified parsing and help output.
type CommandLine struct {
	// Name is the program name shown in help output.
	Name string

	// Version is the program version shown in help output.
	Version string

	// Description appears at the top of help output.
	Description string

	// AlignUsagePerSection determines if usage text alignment is calculated
	// per section (true) or globally across all sections (false).
	AlignUsagePerSection bool

	// Indentation is the number of spaces to indent flag groups.
	Indentation int

	// Padding is the minimum number of spaces between flag names and usage text.
	Padding int

	// SortFlags determines if flags should be sorted alphabetically.
	SortFlags bool

	// Writer specifies where to write help output.
	Writer io.Writer

	// flagSets holds all flag groups in order of creation.
	flagSets []*FlagSet
}

// New creates a new CommandLine with default settings.
func New() *CommandLine {
	cmd := &CommandLine{
		AlignUsagePerSection: DefaultAlignUsagePerSection,
		Indentation:          DefaultIndentation,
		Padding:              DefaultPadding,
		SortFlags:            DefaultSortFlags,

		Writer: os.Stderr,

		flagSets: make([]*FlagSet, 0, 8),
	}

	return cmd
}

// NewFlagSet creates a new FlagSet group with the given name and adds it to the CommandLine.
func (cmd *CommandLine) NewFlagSet(name string) *FlagSet {
	fs := &FlagSet{
		FlagSet: pflag.NewFlagSet(name, pflag.ContinueOnError),

		Name: name,

		Indentation: cmd.Indentation,
		Padding:     cmd.Padding,
		SortFlags:   cmd.SortFlags,
	}

	cmd.flagSets = append(cmd.flagSets, fs)

	return fs
}

// Parse processes command line arguments according to the defined flags.
// It returns an error if flag parsing fails.
func (cmd *CommandLine) Parse() error {
	flagSet := pflag.NewFlagSet("", pflag.ContinueOnError)
	flagSet.Usage = cmd.Usage

	for _, fs := range cmd.flagSets {
		flagSet.AddFlagSet(fs.FlagSet)
	}

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			os.Exit(0)
		}
		return err
	}

	return nil
}

// Usage prints formatted help text to the configured Writer.
func (cmd *CommandLine) Usage() {
	var n int
	w := bufio.NewWriter(cmd.Writer)

	// Program name
	if cmd.Name != "" {
		n += writeString(w, cmd.Name)
	}

	// Version
	if cmd.Version != "" {
		if n != 0 {
			n += writeByte(w, ' ')
		}
		n += writeString(w, cmd.Version)
	}

	// Description
	if cmd.Description != "" {
		if n != 0 {
			n += writeByte(w, '\n')
		}
		n += writeString(w, cmd.Description)
		n += writeByte(w, '\n')
	}

	// Calculate the length of the longest flag name in all the sections
	var maxNameLen int
	for _, fs := range cmd.flagSets {
		maxNameLen = max(maxNameLen, fs.maxNameLength())
	}

	for _, fs := range cmd.flagSets {
		// Calculate the length of the longest flag name in the current section
		if cmd.AlignUsagePerSection {
			maxNameLen = fs.maxNameLength()
		}

		// Apply the proper padding
		fs.setPadding(maxNameLen)

		// Write the section
		if n != 0 {
			n += writeByte(w, '\n')
		}
		n += writeString(w, fs.ToString())
	}

	w.Flush()
}

func writeString(w *bufio.Writer, s string) int {
	n, _ := w.WriteString(s)
	return n
}

func writeByte(w *bufio.Writer, c byte) int {
	if err := w.WriteByte(c); err != nil {
		return 0
	}
	return 1
}
