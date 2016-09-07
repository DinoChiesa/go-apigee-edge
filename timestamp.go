package apigee

import (
  "strconv"
  "time"
  "fmt"
)

// Timestamp represents a time that can be unmarshalled from a JSON string
// formatted as "java time" = milliseconds-since-unix-epoch.
type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
  ms := int64(time.Nanosecond) * int64(time.Time(*t).UnixNano()) /
    int64(time.Millisecond)
  stamp := fmt.Sprintf("%d", ms)
  return []byte(stamp), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Time is expected in RFC3339 or Unix format.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
  ms, err := strconv.ParseInt(string(b), 10, 64)
  if err != nil {
    return err
  }

  *t = Timestamp(time.Unix(0, ms * int64(time.Millisecond) / int64(time.Nanosecond)))
  return nil
}


func (t Timestamp) String() string {
  return time.Time(t).String()
}

// Equal reports whether t and u are equal based on time.Equal
func (t Timestamp) Equal(u Timestamp) bool {
  return time.Time(t).Equal(time.Time(u))
}
