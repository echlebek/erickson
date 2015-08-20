package diff

import "strings"

// SideBySide takes a unified diff and produces a side-by-side
// rendering of it, split into lines.
func SideBySide(diff string) (lhs, rhs []string) {
	lines := strings.Split(diff, "\n")

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

	for i := range lines {
		lhsL, rhsL := lhsM[i], rhsM[i]
		del, ins := deletion(lhsL), insertion(rhsL)
		if !del && !ins {
			lhs = append(lhs, lhsL)
			rhs = append(rhs, rhsL)
			delete(lhsM, i)
			delete(rhsM, i)
		} else if !del && ins {
			// Try to find a matching deletion
			rhs = append(rhs, rhsL)
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
			lhs = append(lhs, deletn)
		} else if del && !ins {
			var insertn string
			lhs = append(lhs, lhsL)
			delete(lhsM, i)
			j, k := hunkRange(i, lines)
			for j := j; j < k; j++ {
				if match, ok := rhsM[j]; ok {
					insertn = match
					delete(rhsM, j)
					break
				}
			}
			rhs = append(rhs, insertn)
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
