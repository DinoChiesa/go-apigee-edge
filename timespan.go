package apigee

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Timespan represents a timespan that can be parsed from a string like "3d"
// meaning "3 days". It will typically be serialized as milliseconds =
// milliseconds-since-unix-epoch.
type Timespan struct {
	time.Duration
}

func NewTimespan(initializer string) *Timespan {
	initializer = strings.ToLower(initializer)
	if strings.HasSuffix(initializer, "d") {
		s := strings.TrimSuffix(initializer, "d")
		i, _ := strconv.Atoi(s)
		i = i * 24
		initializer = fmt.Sprintf("%dh", i)
	}

	duration, _ := time.ParseDuration(initializer)
	return &Timespan{duration}
}

func (span Timespan) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%0.f", span.Duration.Seconds()*1000)
	return []byte(s), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (span *Timespan) UnmarshalJSON(b []byte) error {
	e := errors.New("not supported")
	return e
}

func (span Timespan) String() string {
	return span.Duration.String()
}
