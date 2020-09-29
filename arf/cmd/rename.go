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
	"github.com/AnaMijailovic/NTP/arf/model"
	"github.com/AnaMijailovic/NTP/arf/service"
	"github.com/spf13/cobra"
	"log"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename [path]",
	Short: "Renames files by given criteria",
	Long: `Renames all files (recursively or not) in the given path 
by the given criteria. 
If no path is specified, the current directory path will be used.
Changes are saved in a recovery file that can be later used to restore 
the original names. (look at 'recover' command)`,
	Args: func(cmd *cobra.Command, args []string) error {

		return CheckPaths(args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		path := GetPath(args, 0)

		// Get flag values
		recursiveFlag, _ := cmd.Flags().GetBool("recursive")
		randomFlag, _ := cmd.Flags().GetBool("random")
		removeFlag, _ := cmd.Flags().GetString("remove")
		replaceWithFlag, _ := cmd.Flags().GetString("replaceWith")
		patternFlag, _ := cmd.Flags().GetString("pattern")

		// Check if rename criteria is provided
		if !randomFlag && removeFlag == "" && patternFlag == "" {
			log.Fatal("ERROR: You must provide rename criteria")
		}

		// Check if there is only one criteria provided
		if (randomFlag && removeFlag != "") ||
			(randomFlag && patternFlag != "") ||
			(removeFlag != "" && patternFlag != ""){
			log.Fatal("ERROR: You must provide exactly one rename criteria")
		}

		// Rename
		renameData := model.RenameData{path, recursiveFlag, randomFlag, removeFlag,
			replaceWithFlag, patternFlag }

		errs := service.Rename(&renameData)
		PrintErrors(errs)

	},
}

func init() {
	rootCmd.AddCommand(renameCmd)

	// Set local flags
	renameCmd.Flags().BoolP("recursive", "r", false, "Recursive or not (default: false)")
	renameCmd.Flags().BoolP("random", "n", false, "Generate random new names")
	renameCmd.Flags().StringP("remove", "m", "",
		"Remove a given part of the file name")
	renameCmd.Flags().StringP("replaceWith", "w", "",
		"Replace with a given part of the file name")
	renameCmd.Flags().StringP("pattern", "p", "",
		"Replace with a given part of the file name")
}
