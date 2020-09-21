package service

import (
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"gopkg.in/djherbis/times.v1"
	"log"
	"os"
	"strconv"
	"strings"
)

func ReorganizeFiles(reorganizeData *model.ReorganizeData) []error {
	recoveryFilePath := reorganizeData.Dest + string(os.PathSeparator) + "arfRecover.txt"

	if file, err := os.Open(recoveryFilePath); err == nil {
		file.Close()
		log.Fatal("ERROR: ArfRecover file already exists at the destination path")
	}

	errs := make([]error, 0)

	Walk(reorganizeData.Src, reorganizeData.Recursive, func(path string, info os.FileInfo, err error) error {
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

		if err != nil {
			errs = append(errs, err)
		} else {
			writeRecoveryData(recoveryFilePath, path, newFilePath)
		}

		return nil
	})

	return errs
}

func generateFolderName(path string, reorganizeData *model.ReorganizeData) string {
	if reorganizeData.FileType {
		return generateFolderNameByType(path)
	} else if reorganizeData.FileSize != 0 {
		return generateFolderNameBySize(path, reorganizeData.FileSize)
	} else {
		return generateFolderNameByDate(path, reorganizeData.CreatedDate)
	}
}

func generateFolderNameByType(path string) string {
	fileType, _ := getFileContentType(path)
	return fileType
}

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
