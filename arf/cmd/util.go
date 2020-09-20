package cmd

import (
	"log"
	"os"
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