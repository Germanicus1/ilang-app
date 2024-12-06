package models

import (
	"time"
)

type CustomTime time.Time

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Trim quotes around the timestamp string
	s := string(b)
	s = s[1 : len(s)-1]

	// Attempt to parse the timestamp
	t, err := time.Parse("2006-01-02T15:04:05.999999", s)
	if err != nil {
		return err
	}

	*ct = CustomTime(t)
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	return []byte(t.Format(`"` + "2006-01-02T15:04:05.999999" + `"`)), nil
}

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	CreatedAt CustomTime `json:"created_at"`
	Role      string     `json:"role"`
}

