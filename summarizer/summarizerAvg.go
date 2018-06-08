package summarizer

import (
	"log"
	"strconv"
	"fmt"
)

type avg struct {
	count  int64
	result int64
}

func (a avg) AddValue(stringValue string) error {
	log.Print("Hello!")
	d, err := strconv.ParseInt(stringValue, 0, 64)
	if err != nil {
		return err
	}
	a.count ++
	a.result += d
	return nil
}
func (a avg) Result() string {
	return fmt.Sprintf("%d", a.result/a.count)
}
