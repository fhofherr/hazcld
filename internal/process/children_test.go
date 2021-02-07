package process_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/fhofherr/hazcld/internal/process"
	"github.com/stretchr/testify/assert"
)

func TestHasChildProcess_DirectChildren(t *testing.T) {
	tests := []struct {
		name     string
		pid      int
		re       *regexp.Regexp
		expected bool
	}{
		{
			name:     "process has matching child process",
			pid:      os.Getppid(),
			re:       regexp.MustCompile(`.*\.test`),
			expected: true,
		},
		{
			name:     "process does not have matching child process",
			pid:      os.Getppid(),
			re:       regexp.MustCompile(`good luck, have fun`),
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err := process.HasChildProcess(tt.pid, tt.re)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
