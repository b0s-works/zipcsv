package summarizer

import (
	"strconv"
)

type sum struct {
	result int64
}

func (s sum) AddValue(stringValue string) error {
	d, err := strconv.ParseInt(stringValue, 0, 64)
	if err != nil {
		return err
	}
	s.result += d
	return nil
}

func (s sum) Result() string {
	return strconv.FormatInt(s.result, 10)
}
