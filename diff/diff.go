package diff

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidDiffHeader = errors.New("invalid diff header")

var minusLineRe = regexp.MustCompile(`^--- ([\w\/\.]+)\t?(.+)?`)
var plusLineRe = regexp.MustCompile(`^\+\+\+ ([\w\/\.]+)\t?(.+)?`)
var hunkRangeRe = regexp.MustCompile(`^@@ -(\d+),?(\d+)? \+(\d+),?(\d+)? @@.*`)

// File represents a diff of a single file. Its fields are populated
// by Parse() when it is unmarshaled.
type File struct {
	OldName string     `json:"-"`
	NewName string     `json:"-"`
	OldTime *time.Time `json:"-"`
	NewTime *time.Time `json:"-"`
	Hunks   []Hunk     `json:"-"`
	Lines   []string   `json:"lines"`
}

type DiffLine struct {
	Text string
	Type string
	Line *int
}

func (f *File) UnmarshalJSON(p []byte) error {
	m := make(map[string][]string, 1)
	if err := json.Unmarshal(p, &m); err != nil {
		return err
	}
	lines, ok := m["lines"]
	if !ok {
		return fmt.Errorf("couldn't unmarshal value into File")
	}
	f.Lines = lines
	return f.parse()
}

// Hunk represents a single hunk in a diff.
type Hunk struct {
	Range HunkRange
	LHS   []DiffLine
	RHS   []DiffLine
}

func parseHunkRange(line string) (hunk HunkRange, err error) {
	parts := hunkRangeRe.FindStringSubmatch(line)
	intParts := make([]int, 0, len(parts))
	for _, p := range parts[1:] {
		d, aerr := strconv.Atoi(p)
		if aerr != nil {
			err = aerr
			return
		}
		intParts = append(intParts, d)
	}

	if ln := len(intParts); ln == 2 {
		hunk.MinusL = intParts[0]
		hunk.PlusL = intParts[1]
	} else if ln == 4 {
		hunk.MinusL = intParts[0]
		hunk.MinusS = intParts[1]
		hunk.PlusL = intParts[2]
		hunk.PlusS = intParts[3]
	} else {
		err = fmt.Errorf("invalid hunk index: %s", line)
	}
	return
}

func (h Hunk) Start() int {
	return h.Range.MinusL
}

type HunkRange struct {
	// See https://en.wikipedia.org/wiki/Diff_utility
	MinusL int // Starting line of the hunk in the old file
	MinusS int // Size of the hunk in the old file
	PlusL  int // Starting line of the hunk in the new file
	PlusS  int // Size of the hunk in the new file
}

func (h Hunk) Insertions() (insertions int) {
	for _, line := range h.RHS {
		if strings.HasPrefix(line.Text, "+") {
			insertions++
		}
	}
	return
}

func (h Hunk) Deletions() (deletions int) {
	for _, line := range h.LHS {
		if strings.HasPrefix(line.Text, "-") {
			deletions++
		}
	}
	return
}

// Parse a diff into a set of Files. This function is intended to be
// suitable for use on diff logs from tools like git and hg.
// Unified diff format only.
func ParseFiles(diff string) (files []File, err error) {
	var lines []string
	var file File
	start := strings.Index(diff, "--- ")
	if start < 0 {
		err = fmt.Errorf("invalid diff: %s", diff)
	}
	for _, line := range strings.Split(diff[start:], "\n") {
		if strings.HasPrefix(line, "--- ") {
			if len(lines) > 0 {
				fmt.Println(len(files))
				file, err = NewFile(lines)
				if err != nil {
					return
				}
				files = append(files, file)
				lines = nil
			}
		}
		lines = append(lines, line)
	}

	if len(lines) > 0 {
		file, err = NewFile(lines)
		if err != nil {
			return
		}
		files = append(files, file)
	}

	return
}

func NewFile(lines []string) (File, error) {
	var f File
	f.Lines = lines
	return f, f.parse()
}

func (f *File) parse() error {
	if len(f.Lines) < 3 {
		return fmt.Errorf("invalid diff: %#v", f.Lines)
	}
	minus := minusLineRe.FindStringSubmatch(f.Lines[0])
	plus := plusLineRe.FindStringSubmatch(f.Lines[1])
	if len(minus) > 1 {
		f.OldName = minus[1]
	}
	if len(plus) > 1 {
		f.NewName = plus[1]
	}

	idxLines := indexLines(f.Lines)
	pairs := pairs(idxLines)

	for _, p := range pairs {
		hunkLines := f.Lines[p.x:p.y]
		r, err := parseHunkRange(f.Lines[p.x])
		if err != nil {
			return err
		}
		lhs, rhs := sideBySide(r, hunkLines)
		f.Hunks = append(f.Hunks, Hunk{Range: r, LHS: lhs, RHS: rhs})
	}

	return nil
}

func indexLines(lines []string) []int {
	var idxLines []int
	for idx := range lines {
		if strings.HasPrefix(lines[idx], "@@") {
			idxLines = append(idxLines, idx)
		}
	}
	idxLines = append(idxLines, len(lines))
	return idxLines
}

type pair struct {
	x int
	y int
}

func pairs(ints []int) []pair {
	if len(ints) < 2 {
		return nil
	}
	pairs := make([]pair, len(ints)-1)
	for i := range pairs {
		pairs[i] = pair{ints[i], ints[i+1]}
	}
	return pairs
}

// sideBySide takes a hunk and produces a side-by-side
// rendering of it, split into lines.
// This function has become quite a mess and needs refactoring.
func sideBySide(r HunkRange, lines []string) (lhs, rhs []DiffLine) {
	type idxLine struct {
		Index int
		Line  string
	}

	lhsM := make(map[int]string)
	rhsM := make(map[int]string)

	for i, line := range lines {
		if len(line) > 0 {
			if line[0] == '+' {
				rhsM[i] = line
			} else if line[0] == '-' {
				lhsM[i] = line
			} else {
				lhsM[i] = line
				rhsM[i] = line
			}
		}
	}

	lhsLineCtr := r.MinusL
	rhsLineCtr := r.PlusL

	for i := range lines {
		if lhsLineCtr > (r.MinusL+r.MinusS)+1 && rhsLineCtr > (r.PlusL+r.PlusS)+1 {
			// Avoid slurping too much of the diff by using the index range
			break
		}
		lhsLineNum, rhsLineNum := new(int), new(int)
		*lhsLineNum, *rhsLineNum = lhsLineCtr, rhsLineCtr
		lhsLineCtr++
		rhsLineCtr++
		lhsL, okL := lhsM[i]
		rhsL, okR := rhsM[i]
		if !(okL || okR) {
			continue
		}
		del, ins := deletion(lhsL), insertion(rhsL)
		if !del && !ins {
			lhs = append(lhs, DiffLine{
				Text: html.EscapeString(lhsL),
				Type: "unchanged",
				Line: lhsLineNum})
			rhs = append(rhs, DiffLine{
				Text: html.EscapeString(rhsL),
				Type: "unchanged",
				Line: rhsLineNum})
			delete(lhsM, i)
			delete(rhsM, i)
		} else if !del && ins {
			// Try to find a matching deletion
			rhsType := "insertion"
			lhsType := "unchanged"
			delete(rhsM, i)
			var deletn string
			j, k := hunkRange(i, lines)
			for j := j; j < k; j++ {
				if match, ok := lhsM[j]; ok {
					deletn = match
					delete(lhsM, j)
					break
				}
			}
			if deletn != "" {
				lhsType = "modification"
				rhsType = "modification"
			} else {
				lhsLineNum = nil
				lhsLineCtr--
			}
			rhs = append(rhs, DiffLine{
				Text: html.EscapeString(rhsL),
				Type: rhsType,
				Line: rhsLineNum})
			lhs = append(lhs, DiffLine{
				Text: html.EscapeString(deletn),
				Type: lhsType,
				Line: lhsLineNum})
		} else if del && !ins {
			rhsType := "unchanged"
			lhsType := "deletion"
			var insertn string
			delete(lhsM, i)
			j, k := hunkRange(i, lines)
			for j := j; j < k; j++ {
				if match, ok := rhsM[j]; ok {
					insertn = match
					delete(rhsM, j)
					break
				}
			}
			if insertn != "" {
				rhsType = "modification"
				lhsType = "modification"
			} else {
				rhsLineNum = nil
				rhsLineCtr--
			}
			lhs = append(lhs, DiffLine{
				Text: html.EscapeString(lhsL),
				Type: lhsType,
				Line: lhsLineNum})
			rhs = append(rhs, DiffLine{
				Text: html.EscapeString(insertn),
				Type: rhsType,
				Line: rhsLineNum})
		}
	}

	return
}

func hunkRange(i int, lines []string) (lhs, rhs int) {
	for lhs = i; lhs > -1; lhs-- {
		if len(lines[lhs]) == 0 || (len(lines[lhs]) > 0 && lines[lhs][0] != '-' && lines[lhs][0] != '+') {
			break
		}
	}
	for rhs = i; rhs < len(lines); rhs++ {
		if len(lines[rhs]) == 0 || (len(lines[rhs]) > 0 && lines[rhs][0] != '-' && lines[rhs][0] != '+') {
			break
		}
	}
	return
}

func insertion(line string) bool {
	return len(line) > 0 && line[0] == '+'
}

func deletion(line string) bool {
	return len(line) > 0 && line[0] == '-'
}
