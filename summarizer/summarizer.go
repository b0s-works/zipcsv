package summarizer

import (
	"zipcsv/config"
	"os"
	"fmt"
)

type Storage map[int]ValueSaver

func (s Storage) buildResults() string {
	var summaryString string
	for _, columnValueInterface := range s {
		summaryString += columnValueInterface.Result() + ";"
	}

	return summaryString
}
func (s Storage) GetSummary() string {
	return s.buildResults()
}

type ValueSaver interface {
	AddValue(string) error
	Result() string
}

func New(cfg config.Config, columnsTitles []string) Storage {
	if len(columnsTitles) == 0 {
		fmt.Println("ColumnsTitles variable is empty. Aggregation types can't be linked when no titles used.")
		os.Exit(0)
	}

	var summarizer = make(map[int]ValueSaver)
		fmt.Printf( "ColumnTitles is «%v»\n", columnsTitles )

	var successIndex int
	for columnIndex, columnTitle := range columnsTitles {
		aggregationType := cfg.Aggregation[columnTitle]

		if aggregationType == "" {
			//fmt.Printf( "Column «%v» have not aggregation method specified in configuration file\n", columnTitle )
			continue
		}

		switch aggregationType {
		case "sum":
			fmt.Printf("Column «%+v» by «%+v» will be used for summarizing «%+v»...\n", columnTitle, columnIndex, aggregationType)
			summarizer[columnIndex] = sum{}
			successIndex++
		case "avg":
			fmt.Printf("Column «%+v» by «%+v» will be used for summarizing «%+v»...\n", columnTitle, columnIndex, aggregationType)
			summarizer[columnIndex] = avg{}
			successIndex++
			//TODO ADD GROUPING
		case "countUnique":
			fmt.Printf("Column «%+v» by «%+v» will be used for summarizing «%+v»...\n", columnTitle, columnIndex, aggregationType)
			summarizer[columnIndex] = countUnique{}
			successIndex++
		default:
			fmt.Printf("Aggregation type «%v» is unknown! So column called «%v» will be skipped!\n", aggregationType, columnTitle)
			/*os.Exit(3)*/
		}
	}

	if successIndex == 0 {
		fmt.Printf( "No aggregation methods can be applied. So processing not needed. Program will be terminated.\n" )
		os.Exit(0)
	}

	return summarizer
}

//TODO ADD GROUPING
