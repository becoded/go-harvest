package harvest

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

var ErrTimeParse = errors.New(`ErrTimeParse: should be a string formatted as "15:04"`)

type Time struct {
	time.Time
}

func (t Time) String() string {
	return strings.ToLower(t.Time.Format(time.Kitchen))
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	str = strings.Trim(str, "\"")
	str = strings.ToUpper(str)

	layouts := []string{time.Kitchen, "15:04"}
	for _, layout := range layouts {
		resp, err := time.ParseInLocation(layout, str, time.Local)
		if err == nil {
			t.Time = resp

			return nil
		}
	}

	return ErrTimeParse
}

func (t *Time) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(fmt.Sprintf("\"%s\"", t.String()))

	return buffer.Bytes(), nil
}

func (t *Time) EncodeValues(key string, v *url.Values) error {
	v.Add(key, t.String())

	return nil
}

// Equal reports whether t and u are equal based on time.Equal.
func (t Time) Equal(u Time) bool {
	return t.Time.Equal(u.Time)
}
