package api

import (
	"time"
	"fmt"
	"strings"
)

func (s sphinxRow) Short() string {
	return fmt.Sprintf("%s %s", s.Name, s.preAt().String())
}

func (s sphinxRow) preAt() time.Time {
	return time.Unix(s.PreAt, 0)
}

func (s sphinxRow) Formatted() string {
	lines := []string{s.Name, s.preAt().String()};

	return strings.Join(lines, "\n")
}