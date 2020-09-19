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
	"github.com/AnaMijailovic/NTP/arf/service"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("you must provide path argument")
		}

		// Check if path is valid
		if _, err := os.Open(args[0]); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Get flag values
		recursiveFlag, _ := cmd.Flags().GetBool("recursive")
		emptyFlag, _ := cmd.Flags().GetBool("empty")
		createdBeforeFlag, _ := cmd.Flags().GetString("createdBefore")
		notAccessedAfterFlag, _ := cmd.Flags().GetString("notAccessedAfter")

		// Convert strings to dates
		cbTime := convertStringToDate(createdBeforeFlag, "createdBefore")
		naTime := convertStringToDate(notAccessedAfterFlag, "notAccessedAfter")

		// Check if criteria is provided
		if !emptyFlag && cbTime.IsZero() && naTime.IsZero() {
			log.Fatal("ERROR: You must provide a deletion criteria")
		}


		service.DeleteFiles(args[0], recursiveFlag, emptyFlag, cbTime, naTime )

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolP("recursive", "r", false, "Recursive or not")
	deleteCmd.Flags().BoolP("empty", "e", false, "Delete empty files")
	deleteCmd.Flags().StringP("createdBefore", "b", "",
		"Delete files created before the given date")
	deleteCmd.Flags().StringP("notAccessedAfter", "a", "",
		"Delete files not accessed after the given date")
}

func convertStringToDate(dateStr string, dateName string) time.Time {

	var err error
	var date time.Time

	if dateStr != "" {
		date, err = time.Parse("02-01-2006", dateStr)
	}else {
		date, err = time.Parse("02-01-2006", "01-01-0001")
	}

	if err != nil {
		if _, ok := err.(*time.ParseError); ok {
			log.Fatal("ERROR: " + dateName + " date format is not valid: ", err)
		} else {
			log.Fatal(err)
		}
	}

	return date

}
