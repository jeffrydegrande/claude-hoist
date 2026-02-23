package hoist

import (
	"fmt"
	"strings"
)

// UnifiedDiff returns a unified diff between two strings, labeled with the given names.
func UnifiedDiff(aName, bName, a, b string) string {
	if a == b {
		return ""
	}

	aLines := splitLines(a)
	bLines := splitLines(b)

	// Simple Myers-like LCS diff
	ops := diffLines(aLines, bLines)

	var buf strings.Builder
	fmt.Fprintf(&buf, "--- %s\n", aName)
	fmt.Fprintf(&buf, "+++ %s\n", bName)

	// Group changes into hunks
	hunks := buildHunks(ops, 3)
	for _, h := range hunks {
		fmt.Fprintf(&buf, "@@ -%d,%d +%d,%d @@\n", h.aStart+1, h.aCount, h.bStart+1, h.bCount)
		for _, line := range h.lines {
			buf.WriteString(line)
			buf.WriteByte('\n')
		}
	}

	return buf.String()
}

type diffOp struct {
	kind byte // ' ', '+', '-'
	line string
	aIdx int // line index in a (-1 if added)
	bIdx int // line index in b (-1 if removed)
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	// Remove trailing empty line from trailing newline
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// diffLines produces a sequence of diff operations using LCS.
func diffLines(a, b []string) []diffOp {
	n, m := len(a), len(b)

	// Build LCS table
	lcs := make([][]int, n+1)
	for i := range lcs {
		lcs[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else if lcs[i+1][j] >= lcs[i][j+1] {
				lcs[i][j] = lcs[i+1][j]
			} else {
				lcs[i][j] = lcs[i][j+1]
			}
		}
	}

	// Walk the table to produce ops
	var ops []diffOp
	i, j := 0, 0
	for i < n && j < m {
		if a[i] == b[j] {
			ops = append(ops, diffOp{' ', a[i], i, j})
			i++
			j++
		} else if lcs[i+1][j] >= lcs[i][j+1] {
			ops = append(ops, diffOp{'-', a[i], i, -1})
			i++
		} else {
			ops = append(ops, diffOp{'+', b[j], -1, j})
			j++
		}
	}
	for ; i < n; i++ {
		ops = append(ops, diffOp{'-', a[i], i, -1})
	}
	for ; j < m; j++ {
		ops = append(ops, diffOp{'+', b[j], -1, j})
	}
	return ops
}

type hunk struct {
	aStart, aCount int
	bStart, bCount int
	lines          []string
}

func buildHunks(ops []diffOp, context int) []hunk {
	// Find ranges of changes, expand by context, merge overlapping
	type changeRange struct{ start, end int } // indices into ops
	var changes []changeRange
	for i, op := range ops {
		if op.kind != ' ' {
			if len(changes) > 0 && changes[len(changes)-1].end >= i-1 {
				changes[len(changes)-1].end = i
			} else {
				changes = append(changes, changeRange{i, i})
			}
		}
	}

	if len(changes) == 0 {
		return nil
	}

	// Expand each change range by context and merge overlapping
	type expandedRange struct{ start, end int }
	var expanded []expandedRange
	for _, c := range changes {
		s := c.start - context
		if s < 0 {
			s = 0
		}
		e := c.end + context
		if e >= len(ops) {
			e = len(ops) - 1
		}
		if len(expanded) > 0 && expanded[len(expanded)-1].end >= s-1 {
			expanded[len(expanded)-1].end = e
		} else {
			expanded = append(expanded, expandedRange{s, e})
		}
	}

	var hunks []hunk
	for _, r := range expanded {
		var h hunk
		h.aStart = -1
		h.bStart = -1
		for i := r.start; i <= r.end; i++ {
			op := ops[i]
			h.lines = append(h.lines, string(op.kind)+op.line)
			switch op.kind {
			case ' ':
				if h.aStart == -1 {
					h.aStart = op.aIdx
				}
				if h.bStart == -1 {
					h.bStart = op.bIdx
				}
				h.aCount++
				h.bCount++
			case '-':
				if h.aStart == -1 {
					h.aStart = op.aIdx
				}
				if h.bStart == -1 {
					// Find next context or add line to figure out bStart
					for j := i + 1; j <= r.end; j++ {
						if ops[j].bIdx >= 0 {
							h.bStart = ops[j].bIdx
							break
						}
					}
					if h.bStart == -1 {
						h.bStart = 0
					}
				}
				h.aCount++
			case '+':
				if h.bStart == -1 {
					h.bStart = op.bIdx
				}
				if h.aStart == -1 {
					for j := i + 1; j <= r.end; j++ {
						if ops[j].aIdx >= 0 {
							h.aStart = ops[j].aIdx
							break
						}
					}
					if h.aStart == -1 {
						h.aStart = 0
					}
				}
				h.bCount++
			}
		}
		hunks = append(hunks, h)
	}
	return hunks
}
