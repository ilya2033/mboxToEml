package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type emailStartMap = map[int]bool

const NEW_EMAIL_REGEX = `^From\s(\d+@\w+)\s(\w{3})\s(\w{3})\s(\d{2})\s(\d{2}:\d{2}:\d{2})\s([\+\-]\d{4})\s(\d{4})$`

const DEFAULT_CREATED_FILE_PERMISSIONS = 0644
const DEFAULT_CREATED_FILE_PREFIX = "converted_"
const DEFAULT_CREATED_FILE_EXTENSION = ".eml"

type Config struct {
	ConvertFrom              string
	SaveFolder               string
	ConvertedFilePrefix      string
	ConvertedFilePermissions uint64
	NewEmailPattern          *regexp.Regexp
}

func main() {
	config := createConfig()
	emailStartMap := getEmailStartMap(config)
	convertEmails(emailStartMap, config)
}

func createConfig() *Config {
	convertFrom := flag.String("from", "", "File to convert from")
	destinationFolder := flag.String("toFolder", "", "Folder save to")
	convertedFilePrefix := flag.String("prefix", DEFAULT_CREATED_FILE_PREFIX, "converted file prefix")
	convertedFilePermissions := flag.Uint64("permissions", DEFAULT_CREATED_FILE_PERMISSIONS, "converted file permissions")

	flag.Parse()

	if *convertFrom == "" {
		log.Fatal("Missing file to convert")
	}

	if *destinationFolder == "" {
		log.Fatal("Missing destinationFolder")
	}

	newEmailRegex := regexp.MustCompile(NEW_EMAIL_REGEX)

	config := &Config{
		ConvertFrom:              *convertFrom,
		SaveFolder:               *destinationFolder,
		ConvertedFilePrefix:      *convertedFilePrefix,
		ConvertedFilePermissions: *convertedFilePermissions,
		NewEmailPattern:          newEmailRegex,
	}

	return config
}

func getEmailStartMap(config *Config) emailStartMap {
	emailMap := emailStartMap{}

	file := openFile(config.ConvertFrom)
	fileScanner := createFileScanner(file)
	defer file.Close()

	for i := 0; fileScanner.Scan(); i++ {
		line := fileScanner.Text()
		isNewEmail := checkNewEmailStart(line, config)

		if isNewEmail {
			emailMap[i] = true
		}
	}

	return emailMap
}

func checkNewEmailStart(line string, config *Config) bool {
	isNewEmail := config.NewEmailPattern.MatchString(line)

	return isNewEmail
}

func convertEmails(emailMap emailStartMap, config *Config) {
	var fileWriter *bufio.Writer
	processed := 0

	file := openFile(config.ConvertFrom)
	fileScanner := createFileScanner(file)
	defer file.Close()

	for i := 0; fileScanner.Scan(); i++ {
		line := fileScanner.Text()
		isNewMail := emailMap[i]

		if isNewMail {
			processed++
			fileWriter = setWriterToNewEmail(fileWriter, processed, config)
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

func setWriterToNewEmail(fileWriter *bufio.Writer, number int, config *Config) *bufio.Writer {
	if fileWriter != nil {
		fileWriter.Flush()
	}
	fileName := buildSavePathName(strconv.Itoa(number), config)
	createDirAllPath(fileName)
	fileWriter = createNewEmailFile(fileName, config)

	return fileWriter
}

func createNewEmailFile(fileName string, config *Config) *bufio.Writer {
	fileToWright, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, fs.FileMode(config.ConvertedFilePermissions))
	check(err)

	fileWriter := bufio.NewWriter(fileToWright)
	fmt.Println("Start converting: " + fileName)

	return fileWriter
}

func buildSavePathName(number string, config *Config) string {
	return config.SaveFolder + "/" + config.ConvertedFilePrefix + number + DEFAULT_CREATED_FILE_EXTENSION
}

func createDirAllPath(path string) {
	err := os.MkdirAll(filepath.Dir(path), fs.ModePerm)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
