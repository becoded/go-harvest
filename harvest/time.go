package harvest

import (
	"net/url"
	"strings"
	"time"
)

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
	resp, err := time.Parse(time.Kitchen, str)
	if err != nil {
		return err
	}

	now := time.Now()
	(*t).Time = time.Date(now.Year(), now.Month(), now.Day(), resp.Hour(), resp.Minute(), 0, 0, now.Location())
	return nil
}

func (t *Time) EncodeValues(key string, v *url.Values) error {
	v.Add(key, t.String())
	return nil
}

// Equal reports whether t and u are equal based on time.Equal
func (t Time) Equal(u Time) bool {
	return t.Time.Equal(u.Time)
}
