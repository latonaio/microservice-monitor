package main

import "time"

type jsonTime struct {
	time.Time
}

func (j jsonTime) format() string {
	return j.Time.Format(time.RFC3339)
}

func (j jsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + j.format() + `"`), nil
}

func (j jsonTime) Unix() int64 {
	return j.Time.Unix()
}
