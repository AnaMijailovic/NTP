package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/AnaMijailovic/NTP/arf/model"
	"gopkg.in/djherbis/times.v1"
	"os"
	"path/filepath"
	"strings"
)

// Renames files by given criteria.
// Returns a slice containing eventual errors.
// Fails and returns an error if recovery file
// already exists at a destination path.
func Rename(renameData *model.RenameData) []error {
	errs := make([]error, 0)

	// Check if recovery file already exists at a destination path
	recoveryFilePath := renameData.Path + string(os.PathSeparator) + "arfRecover.txt"

	if file, err := os.Open(recoveryFilePath); err == nil {
		file.Close()
		errs = append(errs, errors.New("ERROR: ArfRecover file already exists at the destination path"))
		return errs
	}

	// Walk func which will be called for each file
	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		newFileName, err := generateNewFileName(path, renameData)
		if err != nil {
			return err
		}

		ext := filepath.Ext(newFileName)
		withoutExt := strings.Replace(newFileName, ext, "", 1)
		if withoutExt == "" {
			return model.UnableToRenameFileError{Err: "Unable to rename " + path + " file."}
		} else {

			newFilePath := filepath.Dir(path) + string(os.PathSeparator) + newFileName

			// Rename a file
			err = renameIfNotExists(path, newFilePath)

			if err == nil {
				err = writeRecoveryData(recoveryFilePath, path, newFilePath)
			}

			if err != nil {
				return err
			}

		}
		return nil
	}

	// Call walk function
	walk := Walk{ root: renameData.Path, recursive: renameData.Recursive, walkFn: walkFn, poolSize: 5}
	walkErrs := walk.startWalking()

	return walkErrs

}

// Generates a new file name according to
// selected criteria.
// Returns an error if criteria is not valid.
func generateNewFileName( oldFilePath string, renameData *model.RenameData) (string, error) {
	oldFileName := filepath.Base(oldFilePath)
	extension := filepath.Ext(oldFilePath)
	oldWithoutExt := strings.Replace(oldFileName, extension, "",1)

	if renameData.Random {
		randomStr, _ := GenerateRandomString(12)
		return  randomStr + extension, nil
	} else if renameData.Remove != "" {
		return strings.ReplaceAll(oldWithoutExt, renameData.Remove, renameData.ReplaceWith) + extension, nil
	} else if renameData.Pattern != "" {
		newName, err := parsePatternString(renameData.Pattern, oldFilePath )
		if err != nil {
			return "", err
		}
		return newName, nil
	}

	return "", errors.New("ERROR: Invalid rename criteria")
}

// Generates a new file name based on a pattern.
// Returns an error if pattern is not valid.
func parsePatternString(patternStr string, oldFilePath string) (string, error){

	// Remove extension from file name
	oldFileName := filepath.Base(oldFilePath)
	ext := filepath.Ext(oldFileName)
	oldFileName = strings.Replace(oldFileName, ext, "", 1)

	firstIndex := strings.Index(patternStr, "{")
	if firstIndex == -1 {
		return "", errors.New("ERROR: The pattern is invalid")
	}

	for strings.Index(patternStr, "{") != -1 {
		startIndex := strings.Index(patternStr, "{")
		endIndex := strings.Index(patternStr, "}")

		if endIndex == -1 || startIndex > endIndex {
			return "", errors.New("ERROR: The pattern is invalid")
		}

		tag := patternStr[startIndex+1 : endIndex]

		if tag == "name" {
			patternStr = strings.Replace(patternStr, patternStr[startIndex:endIndex+1], oldFileName, 1)

		} else if tag == "random" {
			random, _ := GenerateRandomString(12)
			patternStr = strings.Replace(patternStr, patternStr[startIndex:endIndex+1], random, 1 )
		} else if tag == "cDate" {
			timeStat, _ := times.Stat(oldFilePath)
			cDateStr := timeStat.BirthTime().Format("02-01-2006")
			patternStr = strings.Replace(patternStr, patternStr[startIndex:endIndex+1], cDateStr, 1 )
		} else {
			return "", errors.New("ERROR: The pattern is invalid")
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

// Serial version
// Renames files by given criteria.
// Returns a slice containing eventual errors.
// Fails and returns an error if recovery file
// already exists at a destination path.
func RenameS(renameData *model.RenameData) []error {
	errs := make([]error, 0)

	recoveryFilePath := renameData.Path + string(os.PathSeparator) + "arfRecover.txt"

	if file, err := os.Open(recoveryFilePath); err == nil {
		file.Close()
		errs = append(errs, errors.New("ERROR: ArfRecover file already exists at the destination path"))
		return errs
	}

	err := WalkS(renameData.Path, renameData.Recursive, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		newFileName, err := generateNewFileName(path, renameData)
		if err != nil {
			return err
		}

		ext := filepath.Ext(newFileName)
		withoutExt := strings.Replace(newFileName, ext, "", 1)
		if withoutExt == "" {
			errs = append(errs, model.UnableToRenameFileError{Err: "Unable to rename " + path + " file."})
		} else {

			newFilePath := filepath.Dir(path) + string(os.PathSeparator) + newFileName

			// Rename a file
			err = renameIfNotExists(path, newFilePath)

			if err == nil {
				err = writeRecoveryData(recoveryFilePath, path, newFilePath)
			}

			if err != nil {
				errs = append(errs, err)
			}

		}
		return nil
	} )

	if err != nil {
		errs = append(errs, err)
	}
	return errs
}