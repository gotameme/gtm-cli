/*
Copyright © 2024 Sebastian Kroczek <me@xbug.de>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"fmt"
	"github.com/gotameme/gtm-cli/internal/ifacegen"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get name from flag
		name, _ := cmd.Flags().GetString("name")
		// check if directory exists
		if ifacegen.DirectoryExists(name) {
			fmt.Printf("Directory %s already exists", name)
			return
		}
		// create directory with name
		if err := ifacegen.CreateDirectory(name); err != nil {
			fmt.Printf("Error creating Directory: %v", err)
			return
		}
		// create main.go first so mod tidy doesn't  remove the dependency
		if err := ifacegen.CreateMainFile(name); err != nil {
			fmt.Printf("Error creating main.go: %v", err)
			return
		}
		// create go.mod
		if err := ifacegen.CreateModFile(name); err != nil {
			fmt.Printf("Error creating go.mod: %v", err)
			return
		}

		code, err := ifacegen.GenerateInterfaceCode(name)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Printf("%s", code)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// require name flag
	initCmd.PersistentFlags().String("name", "", "Name of the project")
	initCmd.MarkPersistentFlagRequired("name")
}
