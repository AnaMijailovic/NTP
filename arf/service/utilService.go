package service

import (
	"errors"
	"github.com/AnaMijailovic/NTP/arf/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func getFileContentType(filePath string) (string, error) {
	out, _ := os.Open(filePath)
	defer out.Close()
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	contentType = strings.Split(contentType, ";")[0]
	if strings.Contains(contentType, "application") {
		contentType = strings.Split(contentType, "/")[1]
	}

	return contentType, nil
}

func moveFile(src string, dest string) error {

	dirPath := filepath.Dir(dest)

	if _, err := os.Open(dirPath); err != nil {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	// Move file
	err := os.Rename(src, dest)

	if err != nil {
		return err
	}

	return nil

}

func writeRecoveryData(recoveryFilePath string, src string, dest string) error {

	// If the file doesn't exist, create it, or append to the file
	recoveryFile, err := os.OpenFile(recoveryFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return errors.New("ERROR: Invalid recovery file path")
	}

	defer recoveryFile.Close()

	recoveryFile.WriteString(src + "," + dest + "\n")
	return nil
}

func renameIfNotExists(old string, new string) error {

	if old == new {
		return nil
	}

	if file, err := os.Open(new); err == nil {
		file.Close()
		return model.UnableToRenameFileError{Err: "The file with the same name already exists: " + new}
	}

	err := os.Rename(old, new)

	if err != nil {
		return model.UnableToRenameFileError{Err: "Unable to remove: " + new}
	}

	return nil
}