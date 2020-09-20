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
	"fmt"
	"github.com/AnaMijailovic/NTP/arf/service"
	"github.com/spf13/cobra"
)

// recoverCmd represents the recover command
var recoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("you must provide a path argument")
		}

		return CheckPaths(args)
	},
	Run: func(cmd *cobra.Command, args []string) {

		errs := service.Recover(args[0])

		if len(errs) > 0 {
			fmt.Println("ARF was unable to recover the following files: ")
			for _, err := range errs {
				fmt.Println("Error: ", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(recoverCmd)

}