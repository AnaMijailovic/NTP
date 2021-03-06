/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"github.com/AnaMijailovic/NTP/arf/model"
	"github.com/AnaMijailovic/NTP/arf/service"
	"github.com/spf13/cobra"
	"log"
)

// reorganizeCmd represents the reorganize command
var reorganizeCmd = &cobra.Command {

	Use:   "reorganize [sourcePath] [destinationPath]",
	Short: "Reorganizes files by given criteria",
	Long: `Reorganizes all files (recursively or not) in the given path 
by the given criteria.
If no paths are specified, the current directory path will be used
both as a source and as a destination.
Changes are saved in a recovery file that can be later used to restore 
the original organization. (look at 'recover' command)`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 && len(args) != 2 {
			return errors.New("you must provide source and destination paths")
		}
		return CheckPaths(args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		src := GetPath(args, 0)
		dest := GetPath(args, 1)

		// Get flag values
		recursiveFlag, _ := cmd.Flags().GetBool("recursive")
		fileTypeFlag, _ := cmd.Flags().GetBool("fileType")
		fileSizeFlag, _ := cmd.Flags().GetInt64("fileSize")
		createdDateFlag, _ := cmd.Flags().GetString("createdDate")

		// Check if reorganize criteria is provided
		if !fileTypeFlag && fileSizeFlag == 0 && createdDateFlag == "" {
			log.Fatal("ERROR: You must provide reorganize criteria")
		}

		// Validate that there is only one criteria
		if (fileTypeFlag && fileSizeFlag != 0) ||
			(fileTypeFlag && createdDateFlag != "") ||
			(fileSizeFlag != 0 && createdDateFlag != ""){
			log.Fatal("ERROR: You must provide exactly one reorganize criteria: fileType, fileSize or createdDate")
		}

		// Validate createdDate value
		if createdDateFlag != "" && createdDateFlag != "d" && createdDateFlag != "m" && createdDateFlag != "y" {
			log.Fatal("ERROR: Invalid createdDate value. \n\t\t " +
				          "Valid values are: 'd' (day), 'm' (month) and 'y' (year)")
		}

		// Validate fileSize value
		if fileSizeFlag < 0 {
			log.Fatal("ERROR: Invalid fileSize value. It must be a positive number")
		}

		reorganizeData := model.ReorganizeData{src, dest, recursiveFlag, fileTypeFlag,
			fileSizeFlag, createdDateFlag}

		errs := service.ReorganizeFiles(&reorganizeData)
		PrintErrors(errs)
	},
}

func init() {
	rootCmd.AddCommand(reorganizeCmd)

	reorganizeCmd.Flags().BoolP("recursive", "r", false, "Recursive or not (default: false)")
	reorganizeCmd.Flags().BoolP("fileType", "t", false, "Reorganize by file types")
	reorganizeCmd.Flags().Int64P("fileSize", "s", 0,
		"Reorganize by file size; uses given size (in MB) as a step")
	reorganizeCmd.Flags().StringP("createdDate", "c", "",
		"Reorganize by file creation time - y (year); m (month); d(day)")
}