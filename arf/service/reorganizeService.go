package service

import (
	"errors"
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"gopkg.in/djherbis/times.v1"
	"os"
	"strconv"
	"strings"
)

// Reorganizes files by given criteria.
// Returns a slice containing eventual errors.
// Fails and returns an error if recovery file
// already exists at a destination path.
func ReorganizeFiles(reorganizeData *model.ReorganizeData) []error {
	recoveryFilePath := reorganizeData.Dest + string(os.PathSeparator) + "arfRecover.txt"
	errs := make([]error, 0)

	if file, err := os.Open(recoveryFilePath); err == nil {
		file.Close()
		errs = append(errs, errors.New("ERROR: ArfRecover file already exists at the destination path"))
		return errs
	}

	err := Walk(reorganizeData.Src, reorganizeData.Recursive, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		newFolderName := generateFolderName(path, reorganizeData)
		newFolderName = strings.Replace(newFolderName, "/", "_", -1)

		newFolderPath := reorganizeData.Dest + string(os.PathSeparator) + newFolderName

		// check if directory already exists
		if _, err := os.Open(newFolderPath); err != nil {
			err := os.Mkdir(newFolderPath, 0755)
			if err != nil {
				fmt.Println(err)
			}
		}

		// Move a file
		newFilePath :=  newFolderPath + string(os.PathSeparator) + info.Name()

		err = renameIfNotExists(path, newFilePath)

		if err == nil {
			err = writeRecoveryData(recoveryFilePath, path, newFilePath)
		}

		if err != nil {
			errs = append(errs, err)
		}

		return nil
	})

	if err != nil {
		errs = append(errs, err)
	}
	return errs
}

// Generates a new folder name according to
// selected criteria.
func generateFolderName(path string, reorganizeData *model.ReorganizeData) string {
	if reorganizeData.FileType {
		return generateFolderNameByType(path)
	} else if reorganizeData.FileSize != 0 {
		return generateFolderNameBySize(path, reorganizeData.FileSize)
	} else {
		return generateFolderNameByDate(path, reorganizeData.CreatedDate)
	}
}

// The folder name is the file type.
// Path is the absolute path to the file
// whose type is retrieved.
func generateFolderNameByType(path string) string {
	fileType, _ := getFileContentType(path)
	return fileType
}

// The folder name is the file size range (in MB).
// Path is the absolute path to the file
// whose size is retrieved.
func generateFolderNameBySize(path string, step int64) string {
	var label string
	stat, _ := os.Stat(path)
	size := stat.Size() / 1000000 // MB

	groupNum := size / step

	if groupNum == 0 {
		label = strconv.FormatInt(step*groupNum, 10) + "-" + strconv.FormatInt(step*groupNum + step, 10)
	} else {
		label = strconv.FormatInt(step*groupNum + 1, 10) + "-" + strconv.FormatInt(step*groupNum + step, 10)
	}

	return label
}

// The folder name is the creation date.
// It can be full date, month and year or just year,
// which is determined by the dateType.
// Valid values for dateType are 'd', 'm' and 'y'.
// Path is the absolute path to the file
// whose creation date is retrieved.
func generateFolderNameByDate(path string, dateType string) string {

	timeStat, _ := times.Stat(path)
	var folderName string

	if dateType == "d" {
		folderName = timeStat.BirthTime().Format("02-01-2006")
	} else if dateType == "m" {
		folderName = timeStat.BirthTime().Format("01-2006")
	} else { // year
		folderName = timeStat.BirthTime().Format("2006")
	}

	return folderName
}
