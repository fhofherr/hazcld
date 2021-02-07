package process

import (
	"fmt"
	"regexp"

	"github.com/mitchellh/go-ps"
)

// HasChildProcess checks if the passed process ID has a child process with
// a command matching re.
func HasChildProcess(pid int, re *regexp.Regexp) (bool, error) {
	const op = "process/HasChildProcess"

	procs, err := ps.Processes()
	if err != nil {
		return false, fmt.Errorf("%s: get procs: %v", op, err)
	}

	// Build lookup table from process id to respective process. Collect
	// possible candidate processes along the way.
	procMap := make(map[int]ps.Process, len(procs))
	candidates := make([]ps.Process, 0, 10) // should be more than enough in most cases

	for _, p := range procs {
		isCandidate := re.MatchString(p.Executable())
		// Bail out early if we found a direct child that is a candidate.
		if isCandidate && p.PPid() == pid {
			return true, nil
		}
		if isCandidate {
			candidates = append(candidates, p)
		}
		procMap[p.Pid()] = p
	}

	for _, candidate := range candidates {
		if hasAncestor(procMap, candidate, pid) {
			return true, nil
		}
	}
	return false, nil
}

func hasAncestor(procMap map[int]ps.Process, p ps.Process, pid int) bool {
	for p.PPid() != 1 {
		var ok bool

		if p.PPid() == pid {
			return true
		}
		p, ok = procMap[p.PPid()]
		if !ok {
			return false
		}
	}
	return false
}
