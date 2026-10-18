// Package styled renders report documents for an interactive terminal using lipgloss.
package styled

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/example/review-metrics/internal/report"
)

var (
	accent    = lipgloss.AdaptiveColor{Light: "#5A3FC0", Dark: "#B4A2FF"}
	muted     = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#8A8A8A"}
	border    = lipgloss.AdaptiveColor{Light: "#BBBBBB", Dark: "#4A4A4A"}
	titleSt   = lipgloss.NewStyle().Bold(true).Foreground(accent)
	sectionSt = lipgloss.NewStyle().Bold(true).Underline(true).MarginTop(1)
	keySt     = lipgloss.NewStyle().Foreground(muted)
	noteSt    = lipgloss.NewStyle().Foreground(muted).Italic(true)
	headerSt  = lipgloss.NewStyle().Bold(true).Foreground(accent).Padding(0, 1)
	cellSt    = lipgloss.NewStyle().Padding(0, 1)
)

// Render writes docs with lipgloss styling. width is the terminal width (0 = unknown).
func Render(w io.Writer, docs []report.Document, width int) error {
	var b strings.Builder
	for i, d := range docs {
		if i > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(titleSt.Render(d.Title) + "\n")
		b.WriteString(keySt.Render(d.Subtitle) + "\n")
		for _, s := range d.Sections {
			b.WriteString(sectionSt.Render(s.Title) + "\n")
			keyWidth := 0
			for _, kv := range s.Lines {
				keyWidth = max(keyWidth, lipgloss.Width(kv.Key)+1)
			}
			for _, kv := range s.Lines {
				b.WriteString("  " + keySt.Width(keyWidth+2).Render(kv.Key+":") + kv.Value + "\n")
			}
			for _, t := range s.Tables {
				if t.Title != "" {
					b.WriteString("  " + lipgloss.NewStyle().Bold(true).Render(t.Title) + "\n")
				}
				if len(t.Rows) == 0 {
					b.WriteString("  " + noteSt.Render("(none)") + "\n")
					continue
				}
				b.WriteString(indent(renderTable(t, width)) + "\n")
			}
			for _, l := range s.List {
				b.WriteString("  • " + l + "\n")
			}
			if s.Note != "" {
				b.WriteString("  " + noteSt.Render(s.Note) + "\n")
			}
		}
	}
	if _, err := io.WriteString(w, b.String()); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}
	return nil
}

func renderTable(t report.Table, width int) string {
	tb := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(border)).
		Headers(t.Headers...).
		Rows(t.Rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerSt
			}
			if col > 0 && col < len(t.Headers) && !strings.HasPrefix(t.Headers[col], "Top") {
				return cellSt.Align(lipgloss.Right)
			}
			return cellSt
		})
	if width > 4 {
		natural := lipgloss.Width(tb.String())
		if natural > width-2 {
			tb = tb.Width(width - 2)
		}
	}
	return tb.String()
}

func indent(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = "  " + l
	}
	return strings.Join(lines, "\n")
}
