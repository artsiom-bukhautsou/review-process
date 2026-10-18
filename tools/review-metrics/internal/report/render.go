package report

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// Plain writes documents as aligned plain text (for pipes, files and logs).
func Plain(w io.Writer, docs []Document) error {
	ew := &errWriter{w: w}
	for i, d := range docs {
		if i > 0 {
			ew.printf("\n\n")
		}
		ew.printf("%s\n%s\n%s\n", d.Title, strings.Repeat("=", len(d.Title)), d.Subtitle)
		for _, s := range d.Sections {
			ew.printf("\n%s\n%s\n", s.Title, strings.Repeat("-", len(s.Title)))
			if len(s.Lines) > 0 {
				tw := tabwriter.NewWriter(ew, 0, 0, 2, ' ', 0)
				for _, kv := range s.Lines {
					fmt.Fprintf(tw, "  %s:\t%s\n", kv.Key, kv.Value)
				}
				tw.Flush()
			}
			for _, t := range s.Tables {
				ew.printf("\n")
				if t.Title != "" {
					ew.printf("  %s\n", t.Title)
				}
				if len(t.Rows) == 0 {
					ew.printf("  (none)\n")
					continue
				}
				tw := tabwriter.NewWriter(ew, 0, 0, 2, ' ', 0)
				fmt.Fprintf(tw, "  %s\n", strings.Join(t.Headers, "\t"))
				dashes := make([]string, len(t.Headers))
				for i, h := range t.Headers {
					dashes[i] = strings.Repeat("-", len(h))
				}
				fmt.Fprintf(tw, "  %s\n", strings.Join(dashes, "\t"))
				for _, r := range t.Rows {
					fmt.Fprintf(tw, "  %s\n", strings.Join(r, "\t"))
				}
				tw.Flush()
			}
			if len(s.List) > 0 {
				ew.printf("\n")
				for _, l := range s.List {
					ew.printf("  - %s\n", l)
				}
			}
			if s.Note != "" {
				ew.printf("  (%s)\n", s.Note)
			}
		}
	}
	return ew.err
}

// Markdown writes documents as GitHub-flavoured Markdown.
func Markdown(w io.Writer, docs []Document) error {
	ew := &errWriter{w: w}
	for i, d := range docs {
		if i > 0 {
			ew.printf("\n")
		}
		ew.printf("## Review metrics: %s\n\n_%s_\n", d.Title, d.Subtitle)
		for _, s := range d.Sections {
			ew.printf("\n### %s\n\n", s.Title)
			for _, kv := range s.Lines {
				ew.printf("- **%s:** %s\n", mdEscape(kv.Key), mdEscape(kv.Value))
			}
			for _, t := range s.Tables {
				if len(s.Lines) > 0 || t.Title != "" {
					ew.printf("\n")
				}
				if t.Title != "" {
					ew.printf("**%s**\n\n", mdEscape(t.Title))
				}
				if len(t.Rows) == 0 {
					ew.printf("_none_\n")
					continue
				}
				ew.printf("| %s |\n", strings.Join(mdCells(t.Headers), " | "))
				seps := make([]string, len(t.Headers))
				for i := range seps {
					seps[i] = "---"
					if i > 0 {
						seps[i] = "---:"
					}
				}
				if len(t.Headers) > 0 && strings.HasPrefix(t.Headers[len(t.Headers)-1], "Top") {
					seps[len(seps)-1] = "---"
				}
				ew.printf("| %s |\n", strings.Join(seps, " | "))
				for _, r := range t.Rows {
					ew.printf("| %s |\n", strings.Join(mdCells(r), " | "))
				}
			}
			if len(s.List) > 0 {
				ew.printf("\n")
				for _, l := range s.List {
					ew.printf("- %s\n", mdEscape(l))
				}
			}
			if s.Note != "" {
				ew.printf("\n<sub>%s</sub>\n", mdEscape(s.Note))
			}
		}
	}
	return ew.err
}

func mdCells(cells []string) []string {
	out := make([]string, len(cells))
	for i, c := range cells {
		out[i] = mdEscape(c)
	}
	return out
}

func mdEscape(s string) string {
	return strings.NewReplacer("|", `\|`, "<", "&lt;").Replace(s)
}

// errWriter remembers the first write error so rendering code stays linear.
type errWriter struct {
	w   io.Writer
	err error
}

func (e *errWriter) Write(p []byte) (int, error) {
	if e.err != nil {
		return 0, e.err
	}
	n, err := e.w.Write(p)
	if err != nil {
		e.err = fmt.Errorf("writing report: %w", err)
	}
	return n, err
}

func (e *errWriter) printf(format string, args ...interface{}) {
	fmt.Fprintf(e, format, args...)
}
