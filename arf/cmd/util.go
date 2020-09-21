package cmd

import (
	"fmt"
	"log"
	"os"
	"time"
)

func  CheckPaths(paths []string) error {

	// Check if paths are valid
	for _, path := range paths {
		if file, err := os.Open(path); err != nil {
			return err
		} else {
			file.Close()
		}
	}

	return nil
}

func GetPath(args []string, index int) string {
	var err error
	var path string

	if len(args) > index  {
		path = args[index]
	} else {
		path, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	return path
}

func PrintErrors(errs []error, message string) {

	if len(errs) > 0 {
		fmt.Println(message)
		for _, err := range errs {
			fmt.Println("Error: ", err)
		}
	}
}

func ConvertStringToDate(dateStr string, dateName string) time.Time {

	var err error
	var date time.Time

	if dateStr != "" {
		date, err = time.Parse("02-01-2006", dateStr)
	}else {
		date, err = time.Parse("02-01-2006", "01-01-0001")
	}

	if err != nil {
		if _, ok := err.(*time.ParseError); ok {
			log.Fatal("ERROR: " + dateName + " date format is not valid: ", err)
		} else {
			log.Fatal(err)
		}
	}

	return date

}