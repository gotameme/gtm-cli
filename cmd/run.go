/*
Copyright (c) 2024 Sebastian Kroczek <me@xbug.de>

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
	"github.com/gotameme/core"
	"github.com/gotameme/core/ant"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [path to .go file]",
	Short: "Run a Go! Tame! Me! simulation plugin.",
	Long: `Run compiles and executes a Go! Tame! Me! simulation plugin from the specified .go file. For example:

This command compiles the provided Go file into a plugin and runs the simulation with optional parameters. It allows
for dynamic experimentation with different ant behaviors and simulation settings.

Usage:
  gtm-cli run /path/to/your/plugin.go

The 'run' command supports several flags to customize the simulation, including immediate start, headless mode, and 
setting the desired amount of sugar cones.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide the path to the .go file.")
			return
		}

		// resolve the absolute path of the .go file
		absPath, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Println("Error on resolve the absolute path", err)
			return
		}

		// create a temporary file to store the compiled plugin
		tempFile, err := os.CreateTemp("", "plugin-*.so")
		if err != nil {
			fmt.Println("Error on creating temporary file:", err)
			return
		}
		tempFilePath := tempFile.Name()
		tempFile.Close() // close the file to avoid file handle leaks
		// delete the temporary file when the program exits
		defer os.Remove(tempFilePath)

		// define the build command
		buildCmd := exec.Command("go", "build", "-o", tempFilePath, "-buildmode=plugin", absPath)
		// change the working directory to the directory of the .go file
		buildCmd.Dir = filepath.Dir(absPath)
		// redirect the standard output and error to the console
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		// run the build command
		if err := buildCmd.Run(); err != nil {
			fmt.Println("Error compiling ant:", err)
			return
		}

		fmt.Println("Ant successfully compiled. Starting the simulation...")
		// load the compiled plugin
		antPlugin, err := plugin.Open(tempFilePath)
		if err != nil {
			panic(err)
		}
		symNewAnt, err := antPlugin.Lookup("NewAnt")
		if err != nil {
			panic(err)
		}
		NewAnt, ok := symNewAnt.(func(ant.AntOs) ant.Ant)
		if !ok {
			panic("unexpected type from module symbol")
		}

		options := []core.Option{
			core.WithAntConstructor(NewAnt),
		}
		if startImmediately, _ := cmd.Flags().GetBool("startImmediately"); startImmediately {
			options = append(options, core.StartImmediately())
		}
		if headless, _ := cmd.Flags().GetBool("headless"); headless {
			options = append(options, core.Headless())
		}
		if sugar, _ := cmd.Flags().GetInt("sugar"); sugar > 0 {
			options = append(options, core.WithDesiredSugar(sugar))
		}
		core.Run(options...)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("startImmediately", "i", false, "Start the game immediately after the plugin has been compiled.")
	// run headless
	runCmd.Flags().Bool("headless", false, "Run the game in headless mode.")
	runCmd.Flags().Int("sugar", 1, "Set the desired amount of sugar cones.")
}
