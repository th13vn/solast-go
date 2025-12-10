package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
	"github.com/th13vn/solast-go/pkg/parser"
	"github.com/th13vn/solast-go/pkg/version"
)

var (
	// Version information (set during build via ldflags, or detected from build info)
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func init() {
	// Try to get version from Go module build info (works with go install)
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Version != "" && info.Main.Version != "(devel)" {
				Version = info.Main.Version
			}
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					if len(setting.Value) >= 7 {
						GitCommit = setting.Value[:7]
					}
				case "vcs.time":
					BuildTime = setting.Value
				}
			}
		}
	}
}

// Parse command flags
var (
	outputFile  string
	withLoc     bool
	withRange   bool
	tolerant    bool
	prettyPrint bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "solast",
		Short: "Solast-go: Solidity AST Parser",
		Long: `Solast-go is a Solidity AST parser written in Go.
It parses Solidity smart contracts and outputs an AST compatible
with the TypeScript solidity-parser.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, GitCommit, BuildTime),
	}

	// Parse command
	parseCmd := &cobra.Command{
		Use:   "parse [file]",
		Short: "Parse a Solidity file and output AST",
		Long: `Parse a Solidity file and output the Abstract Syntax Tree (AST).
If no file is specified or '-' is given, reads from stdin.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runParse,
	}

	parseCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	parseCmd.Flags().BoolVar(&withLoc, "loc", false, "Include location information (line/column)")
	parseCmd.Flags().BoolVar(&withRange, "range", false, "Include character range information")
	parseCmd.Flags().BoolVar(&tolerant, "tolerant", false, "Tolerant mode (collect errors)")
	parseCmd.Flags().BoolVarP(&prettyPrint, "pretty", "p", true, "Pretty print JSON output")

	// Validate command
	validateCmd := &cobra.Command{
		Use:   "validate [file]",
		Short: "Validate Solidity syntax",
		Long: `Validate the syntax of a Solidity file without producing AST output.
Returns exit code 0 if valid, 1 if there are syntax errors.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runValidate,
	}

	// Version-detect command
	versionCmd := &cobra.Command{
		Use:   "version-detect [file]",
		Short: "Detect Solidity version from pragma",
		Long:  `Detect the Solidity version constraints from a file's pragma directive.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  runVersionDetect,
	}

	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runParse(cmd *cobra.Command, args []string) error {
	input, err := readInput(args)
	if err != nil {
		return err
	}

	opts := &parser.Options{
		Tolerant: tolerant,
		Loc:      withLoc,
		Range:    withRange,
	}

	ast, err := parser.Parse(input, opts)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	var output []byte
	if prettyPrint {
		output, err = json.MarshalIndent(ast, "", "  ")
	} else {
		output, err = json.Marshal(ast)
	}
	if err != nil {
		return fmt.Errorf("JSON encoding error: %w", err)
	}

	return writeOutput(output)
}

func runValidate(cmd *cobra.Command, args []string) error {
	input, err := readInput(args)
	if err != nil {
		return err
	}

	opts := &parser.Options{
		Tolerant: true,
	}

	_, err = parser.Parse(input, opts)
	if err != nil {
		if parserErr, ok := err.(*parser.ParserError); ok {
			fmt.Fprintf(os.Stderr, "Syntax errors found:\n")
			for _, e := range parserErr.Errors {
				fmt.Fprintf(os.Stderr, "  line %d:%d: %s\n", e.Line, e.Column, e.Message)
			}
			os.Exit(1)
		}
		return fmt.Errorf("parse error: %w", err)
	}

	fmt.Println("Syntax OK")
	return nil
}

func runVersionDetect(cmd *cobra.Command, args []string) error {
	input, err := readInput(args)
	if err != nil {
		return err
	}

	detected, err := version.Detect(input)
	if err != nil {
		return fmt.Errorf("version detection error: %w", err)
	}

	fmt.Printf("Pragma: %s\n", detected.Raw)
	fmt.Printf("Version: %s\n", detected.Version)
	if detected.Constraint != "" {
		fmt.Printf("Constraint: %s\n", detected.Constraint)
	}

	return nil
}

func readInput(args []string) (string, error) {
	var reader io.Reader

	if len(args) == 0 || args[0] == "-" {
		reader = os.Stdin
	} else {
		file, err := os.Open(args[0])
		if err != nil {
			return "", fmt.Errorf("cannot open file: %w", err)
		}
		defer file.Close()
		reader = file
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("cannot read input: %w", err)
	}

	return string(content), nil
}

func writeOutput(data []byte) error {
	var writer io.Writer

	if outputFile == "" {
		writer = os.Stdout
	} else {
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("cannot create output file: %w", err)
		}
		defer file.Close()
		writer = file
	}

	_, err := writer.Write(data)
	if err != nil {
		return fmt.Errorf("cannot write output: %w", err)
	}

	// Add newline for stdout
	if outputFile == "" {
		fmt.Println()
	}

	return nil
}

