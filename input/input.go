package input

import (
	"bufio"
	"fmt"
	"log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/sirupsen/logrus"
	"zipcsv/zipcsv"
	"zipcsv/config"
	"zipcsv/summarizer"
)

func UseStdInAsDataInput(cfg config.Config) summarizer.Result {
	var summarizerResult summarizer.Result

	file := os.Stdin
	fi, err := file.Stat()
	checkError(err, "Error on file.Stat()")

	if size := fi.Size(); size != 0 {

		logrus.Debugf("%v bytes available in Stdin\n", size)
		scanner := bufio.NewScanner(os.Stdin)

		scanner.Scan()
		columnsTitles := strings.Split(scanner.Text(), ";")
		checkError(scanner.Err(), "Error on reading standard input")

		summarizerResult = summarizer.New(cfg, columnsTitles)

		for scanner.Scan() {
			fmt.Println(scanner.Text()) // Println will add back the final '\n'
			summarizerResult.ProcessRow(scanner.Text())
		}
		checkError(scanner.Err(), "Error on reading standard input")

		return summarizerResult
	} else {
		logrus.Fatal("Stdin is empty.")
		os.Exit(3)
	}
	return summarizerResult
}

//TODO Remove zip support. Only Raw CSV input
func UseFilesAsDataInput(cfg config.Config) summarizer.Result {
	fmt.Println("File reader will be tried...")

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

	summarizerResult := summarizer.New(cfg, headerColumns)
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
			summarizerResult.ProcessRow(row)
		case err, ok := <-errs:
			if !ok {
				break loop
			}
			fmt.Println(err)
		}
		i++
	}

	log.Printf("%+v", summarizerResult)

	return summarizerResult
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

func checkError(err error, errMsg string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, errMsg+":%+v", err)
		os.Exit(3)
	}
}
