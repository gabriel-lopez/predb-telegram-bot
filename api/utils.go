package api

import (
	"time"
	"fmt"
)

func (s sphinxRow) Short() string {
	return fmt.Sprintf("%s %s", s.Name, s.preAt().String())
}

func (s sphinxRow) preAt() time.Time {
	return time.Unix(s.PreAt, 0)
}