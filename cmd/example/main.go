package main

import (
	"fmt"
	"os"

	"github.com/d3mondev/pflagx"
)

func main() {
	// Create a new command line parser
	cmd := pflagx.New()
	cmd.Name = "myapp"
	cmd.Version = "v1.0.0"
	cmd.Description = "A demonstration of the pflagx package capabilities.\nThis program shows how to organize flags into logical groups."

	// Create different flag groups
	generalFlags := cmd.NewFlagSet("General Options")
	generalFlags.Description = "This is a description for the General Options group."
	verbose := generalFlags.BoolP("verbose", "v", false, "Enable verbose output")
	config := generalFlags.StringP("config", "c", "", "Path to configuration file")
	generalFlags.Bool("dry-run", false, "Perform a trial run with no changes made")

	databaseFlags := cmd.NewFlagSet("Database Options")
	dbHost := databaseFlags.String("db-host", "localhost", "Database server hostname")
	dbPort := databaseFlags.Int("db-port", 5432, "Database server port")
	dbUser := databaseFlags.String("db-user", "postgres", "Database username")
	dbPass := databaseFlags.String("db-password", "", "Database password")
	dbName := databaseFlags.String("db-name", "myapp", "Database name")
	dbSSL := databaseFlags.Bool("db-ssl", false, "Use SSL for database connection")

	outputFlags := cmd.NewFlagSet("Output Options")
	outputFlags.StringP("format", "f", "text", "Output format (text, json, yaml)")
	outputFlags.StringP("output", "o", "-", "Output file (- for stdout)")
	outputFlags.Bool("color", true, "Enable colorized output")
	outputFlags.Int("indent", 2, "Indentation level for structured output")

	advancedFlags := cmd.NewFlagSet("Advanced Options")
	advancedFlags.SortFlags = true
	advancedFlags.Duration("timeout", 0, "Operation timeout (0 for no timeout)")
	advancedFlags.Int("retry", 3, "Number of retry attempts")
	advancedFlags.Float64("factor", 1.5, "Exponential backoff factor")
	advancedFlags.StringSlice("tags", []string{}, "List of tags to apply")
	advancedFlags.Footer = "The previous flags are sorted alphabetically."

	debugFlags := cmd.NewFlagSet("Debug")
	debugFlags.Bool("trace", false, "Enable tracing")
	debugFlags.Lookup("trace").Hidden = true

	exampleFlags := cmd.NewFlagSet("Examples")
	exampleFlags.Footer = `# Basic usage with verbose mode
myapp --verbose

# Specifying database connection parameters
myapp --db-host db.example.com --db-port 3306 --db-user admin --db-password secret --db-ssl

# Using output options with tags
myapp -f json -o output.json --tags frontend,backend,testing

# Dry run with advanced options
myapp --dry-run --timeout 30s --retry 5 --factor 2.0

# Using a configuration file
myapp -c /etc/myapp/config.yaml`

	// Parse command line arguments
	if err := cmd.Parse(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Example of using the parsed flags
	if *verbose {
		fmt.Println("Verbose mode enabled")
	}

	if *config != "" {
		fmt.Printf("Using configuration file: %s\n", *config)
	}

	// Print database connection info
	fmt.Printf("Database connection: %s@%s:%d/%s (SSL: %v)\n",
		*dbUser, *dbHost, *dbPort, *dbName, *dbSSL)

	if *dbPass != "" {
		fmt.Println("Database password is set")
	} else {
		fmt.Println("No database password provided")
	}
}
