package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	fileToConvert := "toConvert.mbox"
	writeFilePrefix := "converted_"
	writeFileExt := ".eml"
	var fileToWright *os.File
	var fileWriter *bufio.Writer
	starts := map[int]bool{}
	limit := 30
	processed := 0

	file, err := os.Open(fileToConvert)
	check(err)

	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)
	pattern := `^From\s(\d+@\w+)\s(\w{3})\s(\w{3})\s(\d{2})\s(\d{2}:\d{2}:\d{2})\s([\+\-]\d{4})\s(\d{4})$`

	// Compile the regex
	regex := regexp.MustCompile(pattern)

	for i := 0; fileScanner.Scan(); i++ {
		line := fileScanner.Text()
		isNewEmail := regex.MatchString(line)

		if isNewEmail {
			starts[i] = true
		}
	}

	file.Close()

	file, err = os.Open(fileToConvert)
	check(err)
	defer file.Close()

	fileScanner = bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)

	for i := 0; fileScanner.Scan(); i++ {
		line := fileScanner.Text()
		isNewMail := starts[i]

		if isNewMail {
			processed++
			if fileWriter != nil {
				fileWriter.Flush()
			}
			if processed >= limit {
				break
			}
			fileName := writeFilePrefix + strconv.Itoa(processed) + writeFileExt
			fileToWright, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			check(err)
			fileWriter = bufio.NewWriter(fileToWright)
			fmt.Println("Start converting: " + fileName)
		}

		fileWriter.WriteString(line + "\n")

	}

}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
