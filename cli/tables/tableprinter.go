package tables

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

const (
	colMaxWidth = 50 // Set a maximum width for each column
)

// TableOutputWriter is the interface to write out tables
type TableOutputWriter interface {
	SetHeaders(headers ...string)
	AddRow(items ...interface{}) error
	Render() error
}

// NewTableWriter gets a new instance of our table output writer
func NewTableWriter(output io.Writer, headers ...string) TableOutputWriter {
	opts := []tablewriter.Option{
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Settings: tw.Settings{
				Separators: tw.Separators{
					BetweenRows: tw.On,
				},
			},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignLeft}},
			Row: tw.CellConfig{
				Alignment:    tw.CellAlignment{Global: tw.AlignLeft},
				Formatting:   tw.CellFormatting{AutoWrap: tw.WrapNormal}, // Wrap long content
				ColMaxWidths: tw.CellWidth{Global: colMaxWidth},
				Padding:      tw.CellPadding{Global: tw.Padding{Left: " ", Right: " "}},
			},
			Footer: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignLeft}},
		}),
	}

	table := tablewriter.NewTable(output, opts...)
	table.Header(headers)

	t := &tableoutputwriter{}
	t.out = output
	t.table = table

	return t
}

// tableoutputwriter is our internal implementation of the TableOutputWriter
type tableoutputwriter struct {
	out   io.Writer
	table *tablewriter.Table
}

// SetHeaders sets the headers for our table and coverts to UPPERCASE for everyone's viewing pleasure
func (t *tableoutputwriter) SetHeaders(headers ...string) {
	t.table.Header(headers)
}

// AddRow appends a new row to our table
func (t *tableoutputwriter) AddRow(items ...interface{}) error {
	row := []string{}

	for _, item := range items {
		row = append(row, fmt.Sprintf("%v", item))
	}

	return t.table.Append(row)
}

// Render emits the generated table to the output once ready
func (t *tableoutputwriter) Render() error {
	if err := t.table.Render(); err != nil {
		return err
	}

	// ensures a break line after we flush the tabwriter
	fmt.Fprintln(t.out)

	return nil
}
