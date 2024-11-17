package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type emailStartMap = map[int]bool

const NEW_EMAIL_REGEX = `^From\s(\d+@\w+)\s(\w{3})\s(\w{3})\s(\d{2})\s(\d{2}:\d{2}:\d{2})\s([\+\-]\d{4})\s(\d{4})$`

const CREATED_FILE_PERMISSIONS = 0644
const CREATED_FILE_PREFIX = "converted_"
const CREATED_FILE_EXTENSION = ".eml"

func main() {
	fileToConvertName := "toConvert.mbox"

	emailStartMap := getEmailStartMap(fileToConvertName)

	converEmails(fileToConvertName, emailStartMap)
}

func getEmailStartMap(fileToConvertName string) emailStartMap {
	emailMap := emailStartMap{}

	file := openFile(fileToConvertName)
	fileScanner := createFileScanner(file)
	defer file.Close()

	for i := 0; fileScanner.Scan(); i++ {
		line := fileScanner.Text()
		isNewEmail := checkNewEmailStart(line)

		if isNewEmail {
			emailMap[i] = true
		}
	}

	return emailMap
}

func checkNewEmailStart(line string) bool {
	regex := regexp.MustCompile(NEW_EMAIL_REGEX)
	isNewEmail := regex.MatchString(line)

	return isNewEmail
}

func converEmails(fileToConvertName string, emailMap emailStartMap) {
	var fileWriter *bufio.Writer
	processed := 0

	file := openFile(fileToConvertName)
	fileScanner := createFileScanner(file)
	defer file.Close()

	for i := 0; fileScanner.Scan(); i++ {
		line := fileScanner.Text()
		isNewMail := emailMap[i]

		if isNewMail {
			processed++
			fileWriter = setWriterToNewEmail(fileWriter, processed)
		}

		if fileWriter != nil {
			fileWriter.WriteString(line + "\n")
		}
	}
}

func createFileScanner(file *os.File) *bufio.Scanner {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	return fileScanner
}
func openFile(fileName string) *os.File {
	file, err := os.Open(fileName)
	check(err)

	return file
}

func setWriterToNewEmail(fileWriter *bufio.Writer, number int) *bufio.Writer {
	if fileWriter != nil {
		fileWriter.Flush()
	}
	fileName := buildSavePathName(strconv.Itoa(number))
	fileWriter = createNewEmailFile(fileName)

	return fileWriter
}

func createNewEmailFile(fileName string) *bufio.Writer {
	fileToWright, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, CREATED_FILE_PERMISSIONS)
	check(err)

	fileWriter := bufio.NewWriter(fileToWright)
	fmt.Println("Start converting: " + fileName)

	return fileWriter
}

func buildSavePathName(number string) string {
	return CREATED_FILE_PREFIX + number + CREATED_FILE_EXTENSION
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
