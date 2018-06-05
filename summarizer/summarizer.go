package summarizer

import (
	"fmt"
	"log"
	"strconv"
	"zipcsv/config"
)

type Storage map[int]ValueSaver

type ValueSaver interface {
	AddValue(string) error
	Result() string
}

func New(cfg config.Config, columnsTitles []string) Storage {
	if len(columnsTitles) != 0 {
		log.Fatal("ColumnsTitles variable is empty. Aggregation types can't be linked when no titles used.")
		return nil
	}

	var summarizer = make(Storage)
	for columnIndex, columnTitle := range columnsTitles {
		aggregationType := cfg.Aggregation[columnTitle]

		switch aggregationType {
		case "sum":
			summarizer[columnIndex] = sum{}
		case "avg":
			summarizer[columnIndex] = avg{}
			//TODO ADD GROUPING
		case "countUnique":
			summarizer[columnIndex] = countUnique{}
		default:
			log.Printf("Aggregation type «%v» is unknown! So column called «%v» will be skipped!", aggregationType, columnTitle)
			/*os.Exit(3)*/
		}
	}

	return summarizer
}

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
	return fmt.Sprintf("%d", s.result)
}

type avg struct {
	count  int64
	result int64
}

func (a avg) AddValue(stringValue string) error {
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

//TODO ADD GROUPING
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
