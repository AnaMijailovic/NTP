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
	Use:   "rename",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		renameData := model.RenameData{path, recursiveFlag, randomFlag, removeFlag,
			replaceWithFlag, patternFlag }

		errs := service.Rename(&renameData)
		PrintErrors(errs)

	},
}

func init() {
	rootCmd.AddCommand(renameCmd)

	renameCmd.Flags().BoolP("recursive", "r", false, "Recursive or not")
	renameCmd.Flags().BoolP("random", "n", false, "Generate random new names")
	renameCmd.Flags().StringP("remove", "m", "",
		"Remove a given part of the file name")
	renameCmd.Flags().StringP("replaceWith", "w", "",
		"Replace with a given part of the file name")
	renameCmd.Flags().StringP("pattern", "p", "",
		"Replace with a given part of the file name")
}
