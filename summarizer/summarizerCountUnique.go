package summarizer

import "fmt"

type countUnique struct {
	unique map[string]bool
}

func (cU countUnique) AddValue(stringValue string) error {
	cU.unique[stringValue] = true

	return nil
}
func (cU countUnique) Result() string {
	return fmt.Sprintf("%d", len(cU.unique))
}
