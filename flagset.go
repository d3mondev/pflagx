package pflagx

import (
	"strings"

	"github.com/spf13/pflag"
)

// FlagSet represents a group of flags with additional formatting options.
// It extends spf13/pflag.FlagSet with descriptive text and layout controls.
type FlagSet struct {
	*pflag.FlagSet

	// Name is the title of the flag group.
	Name string

	// Description appears below the Name and before any flags.
	Description string

	// Footer appears after all flags in the group.
	Footer string

	// Indentation is the number of spaces to indent all content in the group.
	Indentation int

	// Padding is the minimum number of spaces between flag names and their usage text.
	Padding int

	// SortFlags determines if flags should be sorted alphabetically.
	SortFlags bool

	// padding is the computed total padding for aligning usage text.
	padding int
}

// ToString returns the formatted string representation of the FlagSet,
// including the name, description, flags, and footer text.
// Unfortunately, we can't use String as a method name because that would
// override the pflag.FlagSet String method.
func (s *FlagSet) ToString() string {
	sb := strings.Builder{}

	// Indentation
	indentation := strings.Repeat(" ", s.Indentation)

	// Name of the FlagSet
	if s.Name != "" {
		sb.WriteString(s.Name)
		sb.WriteString(":\n")
	}

	// Description of the FlagSet
	if s.Description != "" {
		writeWithPrefix(&sb, s.Description, indentation)
	}

	// Parse all the flags
	s.FlagSet.SortFlags = s.SortFlags
	s.FlagSet.VisitAll(func(f *pflag.Flag) {
		// Skip flags that are hidden
		if f.Hidden {
			return
		}

		// Indentation
		flagBuilder := strings.Builder{}
		flagBuilder.WriteString(indentation)

		// Shorthand flag
		if f.Shorthand != "" {
			flagBuilder.WriteByte('-')
			flagBuilder.WriteString(f.Shorthand)
			flagBuilder.WriteString(", ")
		} else {
			flagBuilder.WriteString("    ")
		}

		// Long flag
		flagBuilder.WriteString("--")
		flagBuilder.WriteString(f.Name)

		// Padding between flag name and usage
		repeat := max(s.padding-flagBuilder.Len(), 0)
		flagBuilder.WriteString(strings.Repeat(" ", repeat))

		// Usage
		if f.Usage != "" {
			addPadding := false
			for line := range strings.SplitSeq(f.Usage, "\n") {
				if addPadding {
					flagBuilder.WriteByte('\n')
					flagBuilder.WriteString(strings.Repeat(" ", s.padding))
				}
				flagBuilder.WriteString(line)
				addPadding = true
			}

			// Default value
			if shouldPrintDefault(f) {
				quotes := false
				switch f.Value.Type() {
				case "string":
					quotes = true
				}

				flagBuilder.WriteString(" (default: ")
				if quotes {
					flagBuilder.WriteByte('"')
				}
				flagBuilder.WriteString(f.DefValue)
				if quotes {
					flagBuilder.WriteByte('"')
				}
				flagBuilder.WriteByte(')')
			}
		}

		// Write the flag string to the main string builder
		sb.WriteString(flagBuilder.String())
		sb.WriteByte('\n')
	})

	// Footer
	if s.Footer != "" {
		writeWithPrefix(&sb, s.Footer, indentation)
	}

	return sb.String()
}

// maxNameLength returns the length of the longest flag name in the FlagSet.
func (s *FlagSet) maxNameLength() int {
	maxLen := 0
	s.FlagSet.VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		maxLen = max(maxLen, len(f.Name))
	})
	return maxLen
}

// setPadding computes and sets the total padding needed to align usage text.
// The padding is calculated as: indentation + shorthand flag space +
// double slash + maximum name length + extra padding.
func (fs *FlagSet) setPadding(maxNameLen int) {
	padding := fs.Indentation // Length of the indentation
	padding += 4              // Shorthand flag "-a, "
	padding += 2              // Double slash of the flag name
	padding += maxNameLen     // Name length
	padding += fs.Padding     // Padding between the name and the usage
	fs.padding = padding
}

// writePrefixedLines writes the string s to the StringBuilder, adding the prefix string
// to the start of each line. A newline is appended after each line, including the last one.
func writeWithPrefix(sb *strings.Builder, s string, prefix string) {
	// Indent each line of text
	lines := strings.SplitSeq(s, "\n")
	for line := range lines {
		sb.WriteString(prefix)
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
}

// shouldPrintDefault returns whether the default value for a flag should
// appear in its usage string.
func shouldPrintDefault(f *pflag.Flag) bool {
	switch f.Value.Type() {
	case "bool":
		return f.DefValue == "true"
	case "stringSlice":
		fallthrough
	case "intSlice":
		fallthrough
	case "uintSlice":
		fallthrough
	case "boolSlice":
		return f.DefValue != "[]"
	default:
		return f.DefValue != ""
	}
}
