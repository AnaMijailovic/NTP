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
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/model"
	"github.com/AnaMijailovic/NTP/arf/service"
	"github.com/spf13/cobra"
	"log"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [path]",
	Short: "Deletes files and folders by given criteria",
	Long: `Deletes all files and folders in the given path (recursively or not)
that match the selected criteria. 
Multiple criteria are connected with 'or' operator.`,
	Args: func(cmd *cobra.Command, args []string) error {
		return CheckPaths(args)

	},
	Run: func(cmd *cobra.Command, args []string) {

		// Get path
		path := GetPath(args, 0)

		// Get flag values
		recursiveFlag, _ := cmd.Flags().GetBool("recursive")
		emptyFlag, _ := cmd.Flags().GetBool("empty")
		createdBeforeFlag, _ := cmd.Flags().GetString("createdBefore")
		notAccessedAfterFlag, _ := cmd.Flags().GetString("notAccessedAfter")

		// Convert strings to dates
		cbTime := ConvertStringToDate(createdBeforeFlag, "createdBefore")
		naTime := ConvertStringToDate(notAccessedAfterFlag, "notAccessedAfter")

		// Check if criteria is provided
		if !emptyFlag && cbTime.IsZero() && naTime.IsZero() {
			log.Fatal("ERROR: You must provide a deletion criteria")
		}

		// Delete
		deleteData := model.DeleteData{path, recursiveFlag, emptyFlag, cbTime,
			naTime}
		filesDeleted := service.DeleteFiles( &deleteData )
		fmt.Println("Deleted files: ", *filesDeleted)

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Set local flags
	deleteCmd.Flags().BoolP("recursive", "r", false, "Recursive or not")
	deleteCmd.Flags().BoolP("empty", "e", false, "Delete empty files")
	deleteCmd.Flags().StringP("createdBefore", "b", "",
		"Delete files created before the given date (dd-mm-yyyy)")
	deleteCmd.Flags().StringP("notAccessedAfter", "a", "",
		"Delete files not accessed after the given date (dd-mm-yyyy)")
}