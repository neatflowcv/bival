package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type app struct {
	Sort    SortCmd    `cmd:"" help:"Sort a BI list by record name."`
	Analyze AnalyzeCmd `cmd:"" help:"Analyze a sorted BI list by building domain entries."`
}

func Run(args []string) error {
	var cli app

	parser, err := kong.New(
		&cli,
		kong.Name("bival"),
		kong.Description("Utilities for sorting and analyzing BI list data."),
	)
	if err != nil {
		return fmt.Errorf("build CLI parser: %w", err)
	}

	ctx, err := parser.Parse(args)
	if err != nil {
		return fmt.Errorf("parse CLI args: %w", err)
	}

	err = ctx.Run()
	if err != nil {
		return fmt.Errorf("run CLI command: %w", err)
	}

	return nil
}
