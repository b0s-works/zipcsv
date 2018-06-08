package summarizer

import (
	"fmt"
	"log"
	"os"
	"strings"
	"zipcsv/config"
)

type Result map[int]DataProcessor

func (s Result) ProcessRow(row string) {
	columnToValue := strings.Split(row, ";")

	if len(columnToValue) < 1 {
		log.Printf("Row «%+v» contain no columns! So row will be skipped.", row)
		return
	}
	if len(columnToValue[0]) < 13 {
		log.Printf("First column value «%+v» lower then 13. So row will be skipped.", columnToValue)
		return
	}

	for columnIndex, columnValue := range columnToValue {
		if s[columnIndex] != nil {
			s[columnIndex].AddValue(columnValue)
		}
	}
}

func (s Result) buildResults() string {
	var summaryString string
	for _, columnValueInterface := range s {
		summaryString += columnValueInterface.Result() + ";"
	}

	return summaryString
}

func (s Result) GetSummary() string {
	return s.buildResults()
}

func New(cfg config.Config, columnsTitles []string) Result {
	if len(columnsTitles) == 0 {
		fmt.Println("ColumnsTitles variable is empty. Aggregation types can't be linked when no titles used.")
		os.Exit(0)
	}

	var summarizer = make(map[int]DataProcessor)
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
