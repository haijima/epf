package main

import (
	"fmt"
	"log/slog"

	"github.com/fatih/color"
	"github.com/haijima/cobrax"
	"github.com/haijima/epf"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(v *viper.Viper, fs afero.Fs) *cobra.Command {
	rootCmd := cobrax.NewRoot(v)
	rootCmd.Use = "epf"
	rootCmd.Short = "Show endpoints"
	rootCmd.Args = cobra.NoArgs
	rootCmd.SetGlobalNormalizationFunc(cobrax.SnakeToKebab)
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Colorization settings
		color.NoColor = color.NoColor || v.GetBool("no-color")
		// Set Log level
		lv.Set(cobrax.VerbosityLevel(v))

		// Read config file
		opts := []cobrax.ConfigOption{cobrax.WithConfigFileFlag(cmd, "config"), cobrax.WithOverrideBy(cmd.Name())}
		if err := cobrax.BindConfigs(v, cmd.Root().Name(), opts...); err != nil {
			return err
		}
		// Bind flags (flags of the command to be executed)
		if err := v.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		slog.Debug("bind flags and config values")
		slog.Debug(cobrax.DebugViper(v))
		return nil
	}
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runRoot(cmd, v, fs)
	}

	rootCmd.Flags().StringP("dir", "d", ".", "The directory to analyze")
	rootCmd.Flags().StringP("pattern", "p", "./...", "The pattern to analyze")
	rootCmd.Flags().String("format", "table", "The output format {table|csv|tsv|md|simple}")

	return rootCmd
}

func runRoot(cmd *cobra.Command, v *viper.Viper, _ afero.Fs) error {
	dir := v.GetString("dir")
	pattern := v.GetString("pattern")
	format := v.GetString("format")

	if format != "table" && format != "csv" && format != "tsv" && format != "md" && format != "simple" {
		return fmt.Errorf("invalid format: %s", format)
	}

	ext, err := epf.AutoExtractor(dir, pattern)
	if err != nil {
		return err
	}

	endpoints, err := epf.FindEndpoints(dir, pattern, ext)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(cmd.OutOrStdout())
	header := table.Row{"#", "Method", "Path", "Function", "Declared Package", "Declared Position"}
	t.AppendHeader(header)
	t.AppendRows(endpointsToRows(endpoints))
	switch format {
	case "csv":
		t.RenderCSV()
	case "tsv":
		t.RenderTSV()
	case "table":
		t.Render()
	case "md":
		t.RenderMarkdown()
	case "simple":
		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateHeader = false
		t.Style().Options.SeparateRows = false
		t.Style().Box.MiddleVertical = " "
		t.Render()
	}
	return nil
}

func endpointsToRows(endpoints []*epf.Endpoint) []table.Row {
	rows := make([]table.Row, 0, len(endpoints))
	for i, e := range endpoints {
		row := table.Row{i + 1, e.Method, e.Path, e.FuncName, e.DeclarePos.PackagePath(true), e.DeclarePos.PositionString()}
		rows = append(rows, row)
	}
	return rows
}
