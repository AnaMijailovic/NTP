package service

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Recover(recoveryFilePath string) []error {

	file, err := os.Open(recoveryFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	errs := recoverFiles(file)

	// Delete recoveryFile
	file.Close()
	err = os.Remove(recoveryFilePath)
	if err != nil {
		log.Fatal(err)
	}

	return errs
}

func recoverFiles(file io.Reader) []error {
	errs := make([]error, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		// Split next line
		paths := strings.Split(scanner.Text(), ",")
		if len(paths) != 2 {
			log.Fatal("ERROR: Recovery file is invalid")
		}
		src, dest := paths[0], paths[1]

		err := moveFile(dest, src)
		if err != nil {
			errs = append(errs, err)
		}

		// Delete destination directory if it is empty
		file, err := os.Stat(filepath.Dir(dest))
		if err != nil {
			errs = append(errs, err)
		} else {
			if file.Size() == 0 {
				os.Remove(filepath.Dir(dest))
			}
		}


	}

	if err := scanner.Err(); err != nil {

		log.Fatal("ERROR: Recovery file is invalid")
	}

	return errs
}