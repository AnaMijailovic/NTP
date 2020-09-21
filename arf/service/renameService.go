package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/AnaMijailovic/NTP/arf/model"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Rename(renameData *model.RenameData) []error {
	errs := make([]error, 0)

	recoveryFilePath := renameData.Path + string(os.PathSeparator) + "arfRecover.txt"

	if file, err := os.Open(recoveryFilePath); err == nil {
		file.Close()
		log.Fatal("ERROR: ArfRecover file already exists at the destination path")
	}

	Walk(renameData.Path, renameData.Recursive, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		newFileName := generateNewFileName(path, renameData)
		ext := filepath.Ext(newFileName)
		withoutExt := strings.Replace(newFileName, ext, "", 1)
		if withoutExt == "" {
			errs = append(errs, errors.New("Unable to rename " + path + " file."))
		} else {

			newFilePath := filepath.Dir(path) + string(os.PathSeparator) + newFileName

			// Rename a file
			err = renameIfNotExists(path, newFilePath)

			if err != nil {
				errs = append(errs, err)
			} else {
				writeRecoveryData(recoveryFilePath, path, newFilePath)
			}

		}
		return nil
	} )

	return errs
}

func generateNewFileName(oldFilePath string, renameData *model.RenameData) string {
	oldFileName := filepath.Base(oldFilePath)
	extension := filepath.Ext(oldFilePath)
	oldWithoutExt := strings.Replace(oldFileName, extension, "",1)

	if renameData.Random {
		randomStr, _ := GenerateRandomString(12)
		return  randomStr + extension
	} else if renameData.Remove != "" {
		return strings.ReplaceAll(oldWithoutExt, renameData.Remove, renameData.ReplaceWith) + extension
	} else if renameData.Pattern != "" {
		newName, _ := parsePatternString(renameData.Pattern, oldFileName )
		return newName
	} else {
		log.Fatal("ERROR: Invalid rename criteria")
	}

	return ""
}

func parsePatternString(patternStr string, oldFileName string) (string, error){
	// var newName string
	// Remove extension from file name
	ext := filepath.Ext(oldFileName)
	oldFileName = strings.Replace(oldFileName, ext, "", 1)

	firstIndex := strings.Index(patternStr, "{")
	if firstIndex == -1 {
		log.Fatal("ERROR: The pattern is invalid")
	}

	for strings.Index(patternStr, "{") != -1 {
		startIndex := strings.Index(patternStr, "{")
		endIndex := strings.Index(patternStr, "}")

		if endIndex == -1 || startIndex > endIndex {
			log.Fatal("ERROR: The pattern is invalid")
		}

		tag := patternStr[startIndex+1 : endIndex]

		if tag == "name" {
			patternStr = strings.Replace(patternStr, patternStr[startIndex:endIndex+1], oldFileName, 1)

		}else if tag == "random" {
			random, _ := GenerateRandomString(12)
			patternStr = strings.Replace(patternStr, patternStr[startIndex:endIndex+1], random, 1 )
		}else {
			log.Fatal("ERROR: The pattern is invalid")
		}

	}

	// Check if pattern is not changed -> no {} in the pattern
	return patternStr + ext, nil

}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
