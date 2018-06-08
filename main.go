package main

import (
	"bufio"
	"fmt"
	"log"
	"io/ioutil"
	"path/filepath"
	"strings"
	"os"
	"zipcsv/zipcsv"
	"zipcsv/config"
	"zipcsv/summarizer"
	"github.com/sirupsen/logrus"
)

/*
дата/время, //2017-07-01T09:28:22
id вестибюля,
номер билета,
UID носителя,
Тип билета,
Тип прохода ( 0 - проходы, 1 - внешние пересадки, -1 - внутренние пересадки),
ѕор порядковый номер поездки по билету (если будет проставлен в системе),
количество оставшихся поездок по билету (если будет проставлено в системе)
 */

func getStdinSize() (int64) {

	file := os.Stdin
	fi, err := file.Stat()
	checkError(err, "Error on file.Stat()")

	return fi.Size()
}

func tryStdInAsDataInput(cfg config.Config) summarizer.Storage {

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	columnsTitles := strings.Split(scanner.Text(), ";")
	checkError(scanner.Err(), "Error on reading standard input")

	summarizerStorage := summarizer.New(cfg, columnsTitles)

	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
		summarizerStorage = processRow(scanner.Text(), summarizerStorage)
	}
	checkError(scanner.Err(), "Error on reading standard input")

	return summarizerStorage
}

func tryFilesAsDataInput(cfg config.Config) summarizer.Storage {
	dir := getCurrentDir()
	files := listFilesOfDir(dir)

	fmt.Printf("Hello world. Current dir is: %q\n", dir)
	fmt.Printf("ZIP files in the directory: %v\n", files)
	if len(files) < 1 {
		fmt.Println("Have no found ZIP files in the directory")
		os.Exit(0)
	}

	rows, errs := zipcsv.ProcessFiles(files)
	header, ok := <-rows
	if !ok {
		fmt.Println("File is empty")
	}

	fmt.Println(header);
	headerColumns := strings.Split(header, ";")

	fmt.Println("Если в первой строке колонки, то они следующие:\n", headerColumns)

	summarizerStorage := summarizer.New(cfg, headerColumns)
	i := 0
loop:
	for {
		select {
		case row, ok := <-rows:
			if i == 0 {
				fmt.Println(row)
			}
			if !ok {
				break loop
			}
			summarizerStorage = processRow(row, summarizerStorage)
		case err, ok := <-errs:
			if !ok {
				break loop
			}
			fmt.Println(err)
		}
		i++
	}

	log.Printf("%+v", summarizerStorage)

	return summarizerStorage
}

func checkError(err error, errMsg string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, errMsg+":%+v", err)
		os.Exit(3)
	}
}

func main() {
	cfg := config.New()
	if size := getStdinSize(); size != 0 {
		logrus.Debugf("%v bytes available in Stdin\n", size)
		if summaryString := tryStdInAsDataInput(cfg).GetSummary(); len(summaryString) > 0 {
			logrus.Debugf(summaryString)
		}
	} else {
		fmt.Println("Stdin is empty. File reader will be tried...")

		if summaryString := tryFilesAsDataInput(cfg).GetSummary(); len(summaryString) > 0 {
			logrus.Debugf(summaryString)
		}
	}
}

func processRow(row string, summarizerStorage summarizer.Storage) summarizer.Storage {
	columnToValue := strings.Split(row, ";")

	if len(columnToValue) < 1 {
		log.Printf("Row «%+v» contain no columns! So row will be skipped.", row)
		return summarizerStorage
	}
	if len(columnToValue[0]) < 13 {
		log.Printf("First column value «%+v» lower then 13. So row will be skipped.", columnToValue)
		return summarizerStorage
	}

	for columnIndex, columnValue := range columnToValue {
		if summarizerStorage[columnIndex] != nil {
			summarizerStorage[columnIndex].AddValue(columnValue)
		}
	}

	return summarizerStorage
}

func getCurrentDir() string {
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		log.Fatal(err)
		return ""
	} else {
		return dir
	}
	return ""
}

func listFilesOfDir(dir string) []string {
	var result []string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() || !zipcsv.IsZIP(file.Name()) {
			continue
		}
		result = append(result, dir+string(os.PathSeparator)+file.Name())
	}

	return result
}
