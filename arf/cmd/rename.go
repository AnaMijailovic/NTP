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
	"log"
	"os"

	"github.com/spf13/cobra"
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

		// Check if path is valid (if provided)
		if len(args) > 0{
			if file, err := os.Open(args[0]); err != nil {
				return err
			}else {
				file.Close()
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rename called")
		var path string
		var err error

		if len(args) > 0 {
			path = args[0]
		}else {
			path, err = os.Getwd()
			if err != nil {
				log.Println(err)
			}
		}

		fmt.Println(path)

		// Get flag values
		recursiveFlag, _ := cmd.Flags().GetBool("recursive")
		randomFlag, _ := cmd.Flags().GetBool("random")
		removeFlag, _ := cmd.Flags().GetString("remove")
		replaceWithFlag, _ := cmd.Flags().GetString("replaceWith")
		patternFlag, _ := cmd.Flags().GetString("pattern")

		renameData := model.RenameData{path, recursiveFlag, randomFlag, removeFlag,
			replaceWithFlag, patternFlag }

		fmt.Println(renameData)
		service.Rename(&renameData)
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
