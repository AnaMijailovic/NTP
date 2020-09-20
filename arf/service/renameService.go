package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Rename(renameData *model.RenameData) {
	recoveryFilePath := renameData.Path + string(os.PathSeparator) + "arfRecover.txt"
	// TODO Something better here?
	os.Remove(recoveryFilePath)

	Walk(renameData.Path, renameData.Recursive, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		newFileName := generateNewFileName(path, renameData)

		newFilePath := filepath.Dir(path) + string(os.PathSeparator) + newFileName

		fmt.Println("New file path: ", newFilePath)

		// Rename a file
		err = os.Rename(path, newFilePath)


		if err != nil {
			log.Fatal(err)
		}

		writeRecoveryData(recoveryFilePath, path, newFilePath)

		return nil
	})
}

func generateNewFileName(oldFilePath string, renameData *model.RenameData) string {
	oldFileName := filepath.Base(oldFilePath)
	if renameData.Random {
		randomStr, _ := GenerateRandomString(12)
		extension := filepath.Ext(oldFilePath)
		randomStr = randomStr + extension
		return  randomStr
	} else if renameData.Remove != "" {
		return strings.ReplaceAll(oldFileName, renameData.Remove, renameData.ReplaceWith)
	} else if renameData.Pattern != "" {
		newName, _ := parsePatternString(renameData.Pattern, oldFileName )
		fmt.Println("New name: ", newName)
		return newName
	}

	//TODO other cases
	return "TODO..."
}

func parsePatternString(patternStr string, oldFileName string) (string, error){
	// var newName string
	// Remove extension from file name
	ext := filepath.Ext(oldFileName)
	oldFileName = strings.Replace(oldFileName, ext, "", 1)

	for strings.Index(patternStr, "{") != -1 {
		startIndex := strings.Index(patternStr, "{")
		endIndex := strings.Index(patternStr, "}")
		tag := patternStr[startIndex+1 : endIndex]

		fmt.Println("Between: ", tag)

		if tag == "name" {
			patternStr = strings.Replace(patternStr, patternStr[startIndex:endIndex+1], oldFileName, 1)

		}else if tag == "random" {
			random, _ := GenerateRandomString(12)
			patternStr = strings.Replace(patternStr, patternStr[startIndex:endIndex+1], random, 1 )
		}

	}
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
