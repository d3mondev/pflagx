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

	// DefaultAlignUsagePerFlagSet determines whether usage text
	// alignment is calculated per FlagSet (true) or globally (false)
	// by default.
	DefaultAlignUsagePerFlagSet = false
)

// Command manages multiple FlagSets and provides unified parsing and help output.
type Command struct {
	// Name is the program name shown in help output.
	Name string

	// Version is the program version shown in help output.
	Version string

	// Description appears at the top of help output.
	Description string

	// AlignUsagePerFlagSet determines if usage text alignment is calculated
	// per FlagSet  (true) or globally across all FlagSets (false).
	AlignUsagePerFlagSet bool

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

// New creates a new Command with default settings.
func New() *Command {
	cmd := &Command{
		AlignUsagePerFlagSet: DefaultAlignUsagePerFlagSet,
		Indentation:          DefaultIndentation,
		Padding:              DefaultPadding,
		SortFlags:            DefaultSortFlags,

		Writer: os.Stderr,

		flagSets: make([]*FlagSet, 0, 8),
	}

	return cmd
}

// NewFlagSet creates a new FlagSet group with the given name and adds it to the Command.
func (cmd *Command) NewFlagSet(name string) *FlagSet {
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
func (cmd *Command) Parse() error {
	pflag.CommandLine = pflag.NewFlagSet("", pflag.ContinueOnError)
	pflag.Usage = cmd.Usage

	for _, fs := range cmd.flagSets {
		pflag.CommandLine.AddFlagSet(fs.FlagSet)
	}

	if err := pflag.CommandLine.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			os.Exit(0)
		}
		return err
	}

	return nil
}

// NArg returns the number of arguments remaining after flags have been processed.
func (cmd *Command) NArg() int {
	return pflag.CommandLine.NArg()
}

// Arg returns the nth argument remaining after flags have been processed.
func (cmd *Command) Arg(n int) string {
	return pflag.CommandLine.Arg(n)
}

// Args returns the non-flag positional arguments.
func (cmd *Command) Args() []string {
	return pflag.CommandLine.Args()
}

// Usage prints formatted help text to the configured Writer.
func (cmd *Command) Usage() {
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

	// Calculate the length of the longest flag name in all the FlagSets
	var maxNameLen int
	for _, fs := range cmd.flagSets {
		maxNameLen = max(maxNameLen, fs.maxNameLength())
	}

	for _, fs := range cmd.flagSets {
		// Calculate the length of the longest flag name in the current FlagSet
		fsMaxNameLen := fs.maxNameLength()

		// Skip the FlagSet if there is nothing to output.
		if fsMaxNameLen == 0 && fs.Description == "" && fs.Footer == "" {
			continue
		}

		if cmd.AlignUsagePerFlagSet {
			maxNameLen = fsMaxNameLen
		}

		// Apply the proper padding
		fs.computePadding(maxNameLen)

		// Write the FlagSet
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
