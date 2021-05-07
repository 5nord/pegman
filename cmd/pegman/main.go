package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/5nord/pegman"
	"github.com/spf13/cobra"
	"golang.org/x/exp/ebnf"
)

var (
	rootCmd = &cobra.Command{
		Use:   "pegman",
		Short: "pegman parses EBNF grammar files and generates output based on the options given.",
		Long: `pegan parses EBNF grammar files and generates output based on the options given.

Grammar File Format
-------------------

Pegman uses the format described in this package: https://pkg.go.dev/golang.org/x/exp/ebnf:

Production  = name "=" [ Expression ] "." .
Expression  = Alternative { "|" Alternative } .
Alternative = Term { Term } .
Term        = name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .

Example:

     Start = ("foo" "bar" | Bar) .
     Bar   = {"bar"} .

`,
		RunE: run,
	}

	format = "json"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&format, "generator", "G", "", "generator to use")
	rootCmd.MarkPersistentFlagRequired("generator")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func run(cmd *cobra.Command, args []string) error {

	grammar, err := pegman.Parse(args[0])
	if err != nil {
		return err
	}

	grammar, err := ebnf.Parse(path, bufio.NewReader(file))
	if err != nil {
		return fmt.Errorf("error parsing %q: %w", path, err)
	}

	prods := pegman.Productions(grammar)
	if len(prods) == 0 {
		return fmt.Errorf("no productions found in %q", path)
	}

	generator, err := findGenerator(format)
	if err != nil {
		return fmt.Errorf("could not find generator for %q: %w", format, err)
	}

	proc := exec.Command(generator)
	proc.Stdin = bytes.NewBuffer(out)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	if err := proc.Start(); err != nil {
		fatal(err)
	}

	proc.Wait()
	return nil
}

// findGenerator returns the path the generator based on the given name.
func findGenerator(name string) (string, error) {
	return exec.LookPath(fmt.Sprintf("pegman-gen-%s", name))
}
